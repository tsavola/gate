// Copyright (c) 2017 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"net"
	"os/exec"
	"sync"
	"syscall"

	"github.com/tsavola/gate/internal/defaultlog"
	"github.com/tsavola/gate/internal/file"
)

var errExecutorDead = errors.New("executor died unexpectedly")

const (
	execOpCreate uint8 = iota
	execOpKill
	execOpSuspend
)

// Executor manages Process resources in an isolated environment.
type Executor struct {
	conn          *net.UnixConn
	idAlloc       <-chan int16
	idFree        chan<- int16
	execRequests  chan execRequest
	killRequests  chan int16
	doneSending   chan struct{}
	doneReceiving chan struct{}

	lock    sync.Mutex
	waiters map[int16]chan<- syscall.WaitStatus
}

func NewExecutor(ctx context.Context, config *Config) (e *Executor, err error) {
	maxProcs := config.maxProcesses()
	if maxProcs > idAllocRangeLen {
		err = errors.New("executor process limit is too high")
		return
	}

	errorLog := config.ErrorLog
	if errorLog == nil {
		errorLog = defaultlog.StandardLogger{}
	}

	var (
		conn *net.UnixConn
		cmd  *exec.Cmd
	)

	if config.DaemonSocket != "" {
		conn, err = dialContainerDaemon(ctx, config)
	} else {
		cmd, conn, err = startContainer(ctx, config)
	}
	if err != nil {
		return
	}

	e = &Executor{
		conn:          conn,
		execRequests:  make(chan execRequest), // No buffering.  Request must be released.
		killRequests:  make(chan int16, 16),   // TODO: how much buffering?
		doneSending:   make(chan struct{}),
		doneReceiving: make(chan struct{}),
		waiters:       make(map[int16]chan<- syscall.WaitStatus),
	}

	e.idAlloc, e.idFree = makeIdAllocator(maxProcs)

	go e.sender(errorLog)
	go e.receiver(errorLog)

	if cmd != nil {
		go containerWaiter(cmd, e.doneSending, errorLog)
	}

	return
}

func (e *Executor) NewProcess(ctx context.Context) (*Process, error) {
	return newProcess(ctx, e)
}

func (e *Executor) execute(ctx context.Context, proc *execProcess, input *file.Ref, output *file.File,
) (err error) {
	select {
	case id, ok := <-e.idAlloc:
		if !ok {
			err = context.Canceled // TODO: ?
			return
		}
		proc.init(e, id)

	case <-ctx.Done():
		err = ctx.Err()
		return
	}

	input.Ref()
	defer func() {
		if err != nil {
			input.Close()
		}
	}()

	select {
	case e.execRequests <- execRequest{proc, input, output}:
		return

	case <-e.doneSending:
		err = errExecutorDead
		return

	case <-e.doneReceiving:
		err = errExecutorDead
		return

	case <-ctx.Done():
		err = ctx.Err() // TODO: include subsystem in error
		return
	}
}

// Close kills all processes.
func (e *Executor) Close() error {
	select {
	case e.killRequests <- math.MaxInt16:
		<-e.doneSending

	case <-e.doneSending:
		// It died on its own.
	}

	<-e.doneReceiving

	// TODO: terminate id allocator

	return e.conn.Close()
}

func (e *Executor) sender(errorLog Logger) {
	var closed bool
	defer func() {
		if !closed {
			close(e.doneSending)
		}
	}()

	buf := make([]byte, 4) // sizeof(struct exec_request)

	// TODO: send multiple entries at once
	for {
		var (
			req  execRequest
			cmsg []byte
		)

		select {
		case req = <-e.execRequests:
			e.lock.Lock()
			e.waiters[req.proc.id] = req.proc.waiter
			e.lock.Unlock()

			// This is like exec_request in runtime/executor/executor.h
			binary.LittleEndian.PutUint16(buf[0:], uint16(req.proc.id))
			buf[2] = execOpCreate

			cmsg = unixRights(req.fds()...)

		case id := <-e.killRequests:
			if id == math.MaxInt16 {
				close(e.doneSending)
				closed = true

				if err := e.conn.CloseWrite(); err != nil {
					errorLog.Printf("executor socket: %v", err)
				}
				return
			}

			op := execOpKill
			if id < 0 {
				id = ^id
				op = execOpSuspend
			}

			// This is like exec_request in runtime/executor/executor.h
			binary.LittleEndian.PutUint16(buf[0:], uint16(id))
			buf[2] = op
		}

		_, _, err := e.conn.WriteMsgUnix(buf, cmsg, nil)
		req.release()
		if err != nil {
			errorLog.Printf("executor socket: %v", err)
			return
		}
	}
}

func (e *Executor) receiver(errorLog Logger) {
	defer close(e.doneReceiving)

	buf := make([]byte, 512*8) // N * sizeof(struct exec_status)
	buffered := 0

	for {
		n, err := e.conn.Read(buf[buffered:])
		if err != nil {
			if err != io.EOF {
				errorLog.Printf("executor socket: %v", err)
			}
			return
		}

		buffered += n
		b := buf[:buffered]

		e.lock.Lock()

		for ; len(b) >= 8; b = b[8:] {
			// This is like exec_status in runtime/executor/executor.h
			var (
				id     = int16(binary.LittleEndian.Uint16(b[0:]))
				status = int32(binary.LittleEndian.Uint32(b[4:]))
			)

			e.waiters[id] <- syscall.WaitStatus(status)
			delete(e.waiters, id)
		}

		e.lock.Unlock()

		buffered = copy(buf, b)
	}
}

// Dead channel will be closed when the executor process dies.  If that wasn't
// requested by calling Close, it indicates an internal error.
func (e *Executor) Dead() <-chan struct{} {
	return e.doneReceiving
}

// Low-level process, tightly coupled with Executor.  See process.go for the
// high-level Process type.
type execProcess struct {
	executor *Executor
	waiter   chan syscall.WaitStatus
	id       int16
}

func (p *execProcess) init(e *Executor, id int16) {
	p.executor = e
	p.waiter = make(chan syscall.WaitStatus, 1)
	p.id = id
}

func (p *execProcess) kill(suspend bool) {
	if p.id < 0 {
		return
	}

	value := p.id
	if suspend {
		value = ^value
	}

	select {
	case p.executor.killRequests <- value:
		p.executor.idFree <- p.id
		p.id = -1

	case <-p.executor.doneSending:

	case <-p.executor.doneReceiving:
	}
}

func (p *execProcess) killWait() (status syscall.WaitStatus, err error) {
	var killRequests chan<- int16
	if p.id >= 0 {
		killRequests = p.executor.killRequests
	}

	for {
		select {
		case killRequests <- p.id:
			killRequests = nil
			p.executor.idFree <- p.id
			p.id = -1

		case <-p.executor.doneSending:
			// No way to kill it anymore.
			killRequests = nil

		case status = <-p.waiter:
			return

		case <-p.executor.doneReceiving:
			err = errExecutorDead
			return
		}
	}
}

type execRequest struct {
	proc   *execProcess
	input  *file.Ref
	output *file.File
}

func (req *execRequest) fds() []int {
	return []int{
		int(req.input.Fd()),
		int(req.output.Fd()),
	}
}

func (req *execRequest) release() {
	if req.proc == nil {
		return
	}

	req.input.Close()
	req.output.Close()
}
