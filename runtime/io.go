// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/tsavola/gate/internal/error/badprogram"
	"github.com/tsavola/gate/packet"
)

const (
	packetReservedOffset = 7
)

// IOState of a program.  Contents are empty unless the program is suspended.
type IOState struct {
	Input  []byte // Buffered data which the program hasn't received yet.
	Output []byte // Buffered data which the program has already sent.
}

type read struct {
	buf packet.Buf
	err error
}

// ioLoop mutates Process and IOState (if any).
func ioLoop(ctx context.Context, services ServiceRegistry, subject *Process, ioState *IOState,
) (err error) {
	var (
		messageInput  = make(chan packet.Buf)
		messageOutput = make(chan packet.Buf)
	)
	discoverer := services.StartServing(ctx, ServiceConfig{maxPacketSize}, messageInput, messageOutput)
	defer close(messageOutput)

	var (
		initialRead []byte
		pendingMsg  packet.Buf
		pendingEvs  []packet.Buf
	)
	if ioState != nil {
		pendingMsg, initialRead, err = splitBufferedPackets(ioState.Output, discoverer)
		if err != nil {
			return
		}
		ioState.Output = nil

		if len(ioState.Input) > 0 {
			pendingEvs = []packet.Buf{ioState.Input} // No need to split packets.
		}
		ioState.Input = nil
	} else {
		subject.writerOut.Close()
		subject.writerOut = nil
	}

	subjectInput := subjectReadLoop(subject.reader, initialRead)
	defer func() {
		if subjectInput != nil {
			for range subjectInput {
			}
		}
	}()

	subjectOutput := subjectWriteLoop(subject.writer)
	subject.writer = nil
	defer func() {
		if subjectOutput != nil {
			close(subjectOutput)
		}
	}()

	var (
		suspended = subject.suspended
		done      = ctx.Done()
	)
	doSuspend := func() {
		suspended = nil
		done = nil

		subject.killSuspend()

		close(subjectOutput)
		subjectOutput = nil
	}

	for subjectInput != nil || pendingMsg != nil {
		var (
			doEv            packet.Buf
			doMessageInput  <-chan packet.Buf
			doMessageOutput chan<- packet.Buf
			doSubjectInput  <-chan read
			doSubjectOutput chan<- packet.Buf
		)

		if len(pendingEvs) > 0 {
			doEv = pendingEvs[0]
		}

		if pendingMsg != nil {
			doMessageOutput = messageOutput
		}

		if doEv == nil {
			doMessageInput = messageInput
			if pendingMsg == nil {
				doSubjectInput = subjectInput
			}
		} else {
			doSubjectOutput = subjectOutput
		}

		select {
		case p := <-doMessageInput:
			pendingEvs = append(pendingEvs, initMessagePacket(p))

		case read, ok := <-doSubjectInput:
			if !ok {
				err = errors.New("read loop terminated unexpectedly")
				return
			}

			switch {
			case read.err == nil:
				msg, ev, opErr := handlePacket(read.buf, discoverer)
				if opErr != nil {
					err = opErr
					return
				}

				if msg != nil {
					pendingMsg = msg
				}
				if ev != nil {
					pendingEvs = append(pendingEvs, ev)
				}

			case read.err == io.EOF:
				subjectInput = nil

				if subjectOutput != nil {
					close(subjectOutput)
					subjectOutput = nil
				}

				if ioState != nil {
					if suspended == nil { // Suspended.
						ioState.Output = append(pendingMsg, read.buf...)
						pendingMsg = nil

						ioState.Input, err = ioutil.ReadAll(subject.writerOut)
						if err != nil {
							return
						}

						var pendingLen int
						for _, ev := range pendingEvs {
							pendingLen += len(ev)
						}

						if n := len(ioState.Input) + pendingLen; cap(ioState.Input) < n {
							ioState.Input = append(make([]byte, 0, n), ioState.Input...)
						}

						for _, ev := range pendingEvs {
							ioState.Input = append(ioState.Input, ev...)
						}
					}

					subject.writerOut.Close()
					subject.writerOut = nil
				}

			case read.err != nil:
				err = read.err
				return
			}

		case doMessageOutput <- pendingMsg:
			pendingMsg = nil

		case doSubjectOutput <- doEv:
			if len(pendingEvs) > 0 {
				pendingEvs = pendingEvs[1:]
			}

		case <-suspended:
			doSuspend()

		case <-done:
			doSuspend()
		}
	}

	return
}

func subjectReadLoop(r *os.File, partial []byte) <-chan read {
	reads := make(chan read)

	go func() {
		defer close(reads)

		header := make([]byte, packet.HeaderSize)

		for {
			offset := copy(header, partial)
			partial = partial[offset:]

			if n, err := io.ReadFull(r, header[offset:]); err != nil {
				if err != io.EOF {
					err = fmt.Errorf("subject read: %v", err)
				}
				reads <- read{buf: header[:offset+n], err: err}
				return
			}

			size := binary.LittleEndian.Uint32(header)
			if size < packet.HeaderSize || size > maxPacketSize {
				reads <- read{err: badprogram.Errorf("runtime: invalid op packet size: %d", size)}
				return
			}

			buf := make([]byte, size)
			offset = copy(buf, header)
			offset += copy(buf[offset:], partial)
			partial = nil

			if n, err := io.ReadFull(r, buf[offset:]); err != nil {
				if err != io.EOF {
					err = fmt.Errorf("subject read: %v", err)
				}
				reads <- read{buf: buf[:offset+n], err: err}
				return
			}

			reads <- read{buf: buf}
		}
	}()

	return reads
}

func subjectWriteLoop(w *os.File) chan<- packet.Buf {
	writes := make(chan packet.Buf)

	go func() {
		defer w.Close()

		for buf := range writes {
			if _, err := w.Write(buf); err != nil {
				return
			}
		}
	}()

	return writes
}

func initMessagePacket(p packet.Buf) packet.Buf {
	if len(p) < packet.HeaderSize || len(p) > maxPacketSize {
		panic(errors.New("invalid outgoing message packet buffer length"))
	}

	if p[packetReservedOffset] != 0 {
		panic(errors.New("reserved byte is nonzero in outgoing message packet header"))
	}

	if p.Code() < 0 {
		panic(errors.New("negative service code in outgoing message packet header"))
	}

	if p.Domain() > packet.DomainData {
		panic(errors.New("invalid domain in outgoing message packet header"))
	}

	// Service implementations only need to initialize code and domain fields.
	binary.LittleEndian.PutUint32(p[packet.OffsetSize:], uint32(len(p)))

	return p
}

func handlePacket(p packet.Buf, discoverer ServiceDiscoverer) (msg, reply packet.Buf, err error) {
	switch code := p.Code(); {
	case code >= 0:
		msg, err = checkServicePacket(p, discoverer)
		if err != nil {
			return
		}

	case code == packet.CodeServices:
		reply, err = handleServicesPacket(p, discoverer)
		if err != nil {
			return
		}

	default:
		err = badprogram.Errorf("invalid code in incoming packet: %d", code)
		return
	}

	return
}

func splitBufferedPackets(buf []byte, discoverer ServiceDiscoverer,
) (msg packet.Buf, tail []byte, err error) {
	if len(buf) < packet.HeaderSize {
		tail = buf
		return
	}

	size := binary.LittleEndian.Uint32(buf[packet.OffsetSize:])
	if size < packet.HeaderSize || size > maxPacketSize {
		err = badprogram.Errorf("buffered packet has invalid size: %d", size)
		return
	}

	if uint32(len(buf)) < size {
		tail = buf
		return
	}

	p := packet.Buf(buf[:size])

	switch code := p.Code(); {
	case code >= 0:
		msg, err = checkServicePacket(p, discoverer)
		if err != nil {
			return
		}

	default:
		err = badprogram.Errorf("invalid code in buffered packet: %d", code)
		return
	}

	tail = buf[size:]
	return
}

func checkServicePacket(p packet.Buf, discoverer ServiceDiscoverer) (msg packet.Buf, err error) {
	if x := p[packetReservedOffset]; x != 0 {
		err = badprogram.Errorf("reserved byte has value 0x%02x in buffered packet header", x)
		return
	}

	if int(p.Code()) >= discoverer.NumServices() {
		err = badprogram.Errorf("invalid service code in packet: %d", p.Code())
		return
	}

	switch p.Domain() {
	case packet.DomainCall, packet.DomainFlow:

	case packet.DomainData:
		if n := len(p); n < packet.DataHeaderSize {
			err = badprogram.Errorf("data packet is too short: %d bytes", n)
			return
		}

	default:
		err = badprogram.Errorf("invalid domain in packet: %d", p.Domain())
		return
	}

	msg = p
	return
}
