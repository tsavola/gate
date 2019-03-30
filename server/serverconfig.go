// Copyright (c) 2017 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"fmt"
	"io"

	"github.com/gogo/protobuf/proto"
	"github.com/tsavola/gate/image"
	"github.com/tsavola/gate/packet"
	"github.com/tsavola/gate/runtime"
	"github.com/tsavola/gate/server/detail"
	"github.com/tsavola/gate/snapshot"
)

type InstanceConnector interface {
	// Connect allocates a new I/O stream.  The returned function is used to
	// drive I/O between network connection and instance.  If it's non-nil, a
	// connection was made.
	Connect(context.Context) func(context.Context, io.Reader, io.Writer) error

	// Close disconnects remaining connections.  Currently blocked and future
	// Connect calls will return nil.
	Close() error
}

type InstanceServices interface {
	InstanceConnector
	runtime.ServiceRegistry
}

type ServiceStarter func(context.Context, runtime.ServiceConfig, []snapshot.Service, chan<- packet.Buf, <-chan packet.Buf) (runtime.ServiceDiscoverer, []runtime.ServiceState, error)

type instanceServices struct {
	InstanceConnector
	startServing ServiceStarter
}

func NewInstanceServices(c InstanceConnector, r runtime.ServiceRegistry) InstanceServices {
	return &instanceServices{c, r.StartServing}
}

func WrapInstanceServices(c InstanceConnector, f ServiceStarter) InstanceServices {
	return &instanceServices{c, f}
}

func (is *instanceServices) StartServing(ctx context.Context, config runtime.ServiceConfig, state []snapshot.Service, send chan<- packet.Buf, recv <-chan packet.Buf,
) (runtime.ServiceDiscoverer, []runtime.ServiceState, error) {
	return is.startServing(ctx, config, state, send, recv)
}

type Event interface {
	EventName() string
	EventType() int32
	proto.Message
}

type Config struct {
	ImageStorage   image.Storage
	ProcessFactory runtime.ProcessFactory
	AccessPolicy   Authorizer
	Monitor        func(Event, error)
}

func (c *Config) Configured() bool {
	return c.ProcessFactory != nil && c.AccessPolicy != nil
}

func AllocateIface(name string) detail.Iface {
	value, found := detail.Iface_value[name]
	if !found {
		value = int32(len(detail.Iface_name))
		detail.Iface_name[value] = name
		detail.Iface_value[name] = value
	}
	return detail.Iface(value)
}

func RegisterIface(value int32, name string) {
	if n, found := detail.Iface_name[value]; found && n != name {
		panic(fmt.Errorf("iface %d (%s) already exists with different name: %s", value, name, n))
	}
	if v, found := detail.Iface_value[name]; found && v != value {
		panic(fmt.Errorf("iface %s (%d) already exists with different value: %d", name, value, v))
	}
	detail.Iface_name[value] = name
	detail.Iface_value[name] = value
}
