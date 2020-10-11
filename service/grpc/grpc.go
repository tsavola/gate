// Copyright (c) 2020 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"context"
	"errors"
	"io"

	"gate.computer/gate/packet"
	"gate.computer/gate/principal"
	"gate.computer/gate/runtime"
	"gate.computer/gate/service"
	"gate.computer/gate/service/grpc/api"
	"github.com/google/uuid"
	"github.com/tsavola/mu"
	"google.golang.org/grpc"
)

type procKey struct {
	b []byte
	n int
}

var (
	procMu   mu.Mutex
	procKeys = make(map[runtime.ProcessKey]procKey)
)

func getProcKey(ctx context.Context) (key []byte) {
	opaque := runtime.MustContextProcessKey(ctx)

	procMu.Guard(func() {
		if x, found := procKeys[opaque]; found {
			x.n++
			procKeys[opaque] = x
			key = x.b
		}
	})
	if key != nil {
		return
	}

	array := uuid.New()
	key = array[:]
	procMu.Guard(func() {
		if x, found := procKeys[opaque]; found {
			x.n++
			procKeys[opaque] = x
			key = x.b
		} else {
			procKeys[opaque] = procKey{key, 1}
		}
	})
	return
}

func putProcKey(ctx context.Context) {
	opaque := runtime.MustContextProcessKey(ctx)

	procMu.Lock()
	defer procMu.Unlock()

	x := procKeys[opaque]
	x.n--
	if x.n == 0 {
		delete(procKeys, opaque)
	} else {
		procKeys[opaque] = x
	}
}

// Conn is a connection to a gRPC server.
type Conn struct {
	Services []*Service

	conn *grpc.ClientConn
}

// New takes ownership of conn.
func New(ctx context.Context, conn *grpc.ClientConn) (*Conn, error) {
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	r, err := api.NewRootClient(conn).Init(ctx, &api.InitRequest{})
	if err != nil {
		return nil, err
	}

	c := &Conn{
		conn: conn,
	}

	servClient := api.NewServiceClient(conn)
	instClient := api.NewInstanceClient(conn)

	for _, info := range r.Services {
		c.Services = append(c.Services, newService(servClient, instClient, info))
	}

	conn = nil
	return c, nil
}

// DialContext connects to a gRPC server.
func DialContext(ctx context.Context, target string, opts ...grpc.DialOption) (*Conn, error) {
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	c, err := New(ctx, conn)
	if err != nil {
		return nil, err
	}

	conn = nil
	return c, nil
}

// Register the services which are accessible through the connection.
func (c *Conn) Register(r *service.Registry) error {
	for _, s := range c.Services {
		if err := r.Register(s); err != nil {
			return err
		}
	}
	return nil
}

// Close the gRPC client connection.
func (c *Conn) Close() error {
	err := c.conn.Close()
	c.conn = nil
	return err
}

type Service struct {
	c          api.ServiceClient
	instClient api.InstanceClient
	info       *api.ServiceInfo
}

func newService(c api.ServiceClient, instClient api.InstanceClient, info *api.ServiceInfo) *Service {
	return &Service{
		c:          c,
		instClient: instClient,
		info:       info,
	}
}

func (s *Service) Service() service.Service {
	return service.Service{
		Name:     s.info.Name,
		Revision: s.info.Revision,
	}
}

func (s *Service) Discoverable(ctx context.Context) bool {
	if s.info.RequirePrincipal && principal.ContextID(ctx) == nil {
		return false
	}

	return true
}

func (s *Service) CreateInstance(ctx context.Context, config service.InstanceConfig, snapshot []byte,
) (service.Instance, error) {
	key := getProcKey(ctx)
	defer func() {
		if key != nil {
			putProcKey(ctx)
		}
	}()

	r, err := s.c.CreateInstance(ctx, &api.CreateInstanceRequest{
		Name:     s.info.Name,
		Config:   newInstanceConfig(ctx, config, key),
		Snapshot: snapshot,
	})
	if err != nil {
		return nil, err
	}

	if r.Error != "" {
		return nil, errors.New(r.Error) // TODO: a ModuleError
	}

	key = nil
	return newInstance(ctx, s.instClient, r.Id, config), nil
}

func newInstanceConfig(ctx context.Context, config service.InstanceConfig, key []byte) *api.InstanceConfig {
	r := &api.InstanceConfig{
		MaxSendSize: int32(config.MaxSendSize),
		ProcessKey:  key,
	}

	if pri := principal.ContextID(ctx); pri != nil {
		r.PrincipalId = pri.String()
	}

	if id, ok := principal.ContextInstanceUUID(ctx); ok {
		r.InstanceUuid = id[:]
	}

	return r
}

type instance struct {
	service.InstanceBase

	c    api.InstanceClient
	id   []byte
	code packet.Code

	stream   api.Instance_ReceiveClient
	leftout  <-chan []byte
	incoming []byte
}

func newInstance(ctx context.Context, c api.InstanceClient, id []byte, config service.InstanceConfig) *instance {
	return &instance{
		c:    c,
		id:   id,
		code: config.Code,
	}
}

func (inst *instance) Ready(ctx context.Context) error {
	stream, err := inst.c.Receive(ctx, &api.ReceiveRequest{
		Id: inst.id,
	})
	if err != nil {
		return err
	}

	inst.stream = stream
	return nil
}

func (inst *instance) Start(ctx context.Context, out chan<- packet.Buf) error {
	c := make(chan []byte, 1)
	go receiveForward(ctx, inst.code, out, inst.stream, c)
	inst.leftout = c
	inst.stream = nil
	return nil
}

func (inst *instance) Handle(ctx context.Context, out chan<- packet.Buf, p packet.Buf) error {
	if len(inst.incoming) == 0 {
		_, err := inst.c.Handle(ctx, &api.HandleRequest{
			Id:   inst.id,
			Data: p,
		})
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():

		default:
			return err
		}

	}

	inst.incoming = append(inst.incoming, p...)
	return nil
}

func (inst *instance) Shutdown(ctx context.Context) error {
	putProcKey(ctx)

	_, err := inst.c.Shutdown(ctx, &api.ShutdownRequest{
		Id: inst.id,
	})
	return err
}

func (inst *instance) Suspend(ctx context.Context) ([]byte, error) {
	putProcKey(ctx)

	_, err := inst.c.Suspend(ctx, &api.SuspendRequest{
		Id: inst.id,
	})
	if err != nil {
		return nil, err
	}

	var outgoing []byte
	if inst.leftout != nil {
		outgoing = <-inst.leftout
	}

	r, err := inst.c.Snapshot(ctx, &api.SnapshotRequest{
		Id:       inst.id,
		Outgoing: outgoing,
		Incoming: inst.incoming,
	})
	if err != nil {
		return nil, err
	}

	return r.Value, nil
}

func receiveForward(ctx context.Context, code packet.Code, out chan<- packet.Buf, stream api.Instance_ReceiveClient, leftout chan<- []byte) {
	defer close(leftout)

	for {
		r, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}

		p := mustBePacket(r.Value)
		p.SetCode(code)

		select {
		case out <- p:

		case <-ctx.Done():
			leftout <- receiveBuffer(p, stream)
			return
		}
	}
}

func receiveBuffer(initial packet.Buf, stream api.Instance_ReceiveClient) (buf []byte) {
	initial.SetSize()
	buf = initial

	for {
		r, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}

		p := mustBePacket(r.Value)
		p.SetSize()

		buf = append(buf, p...)
	}
}

func mustBePacket(b []byte) packet.Buf {
	if len(b) < packet.HeaderSize {
		panic("invalid packet received from gRPC service")
	}
	return b
}
