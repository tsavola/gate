// Copyright (c) 2017 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"encoding/binary"
)

// Code represents the source or destination of a packet.
type Code int16

const (
	CodeServices Code = -1
)

type Domain uint8

const (
	DomainCall Domain = iota
	DomainInfo
	DomainFlow
	DomainData
)

func (d Domain) String() string {
	switch d {
	case DomainCall:
		return "call"

	case DomainInfo:
		return "info"

	case DomainFlow:
		return "flow"

	case DomainData:
		return "data"

	default:
		return "<invalid domain>"
	}
}

const (
	Alignment = 8

	// Packet header
	OffsetSize     = 0
	OffsetCode     = 4
	OffsetDomain   = 6
	offsetReserved = 7
	HeaderSize     = 8

	// Services packet header
	OffsetServicesCount = HeaderSize + 0
	ServicesHeaderSize  = HeaderSize + 2

	// Flow packet header
	FlowHeaderSize = HeaderSize

	// Data packet header
	OffsetDataID   = HeaderSize + 0
	OffsetDataNote = HeaderSize + 4
	DataHeaderSize = HeaderSize + 8
)

const (
	flowOffsetID        = 0
	flowOffsetIncrement = 4
	flowSize            = 8
)

// Align packet length up to a multiple of packet alignment.
func Align(length int) int {
	return (length + (Alignment - 1)) &^ (Alignment - 1)
}

// Buf holds a packet.
type Buf []byte

func Make(code Code, domain Domain, packetSize int) Buf {
	b := Buf(make([]byte, packetSize))
	binary.LittleEndian.PutUint16(b[OffsetCode:], uint16(code))
	b[OffsetDomain] = byte(domain)
	return b
}

func MakeCall(code Code, contentSize int) Buf {
	return Make(code, DomainCall, HeaderSize+contentSize)
}

func MakeInfo(code Code, contentSize int) Buf {
	return Make(code, DomainInfo, HeaderSize+contentSize)
}

func MakeFlow(code Code, id int32, increment int32) Buf {
	b := MakeFlows(code, 1)
	b.Set(0, id, increment)
	return Buf(b)
}

// Code is the program instance-specific service identifier.
func (b Buf) Code() Code {
	return Code(binary.LittleEndian.Uint16(b[OffsetCode:]))
}

func (b Buf) Domain() Domain {
	return Domain(b[OffsetDomain])
}

// Content of a received packet, or buffer for initializing sent packet.
func (b Buf) Content() []byte {
	return b[HeaderSize:]
}

// Split a packet into two parts.  The headerSize parameter determins how many
// bytes are initialized in the second part: the header is copied from the
// first part.  The length of the first part is given as the prefixLen
// parameter.  If the buffer is too short for the second part, the length of
// the second buffer will be zero.
func (b Buf) Split(headerSize, prefixLen int) (prefix, unused Buf) {
	prefixCap := Align(prefixLen)
	if prefixCap > len(b) {
		prefixCap = len(b)
	}

	prefix = b[:prefixLen:prefixCap]
	unused = b[prefixCap:]

	if len(unused) < headerSize {
		unused = unused[0:]
		return
	}

	copy(unused, prefix[:headerSize])
	return
}

// FlowBuf holds a flow packet.
type FlowBuf Buf

func MakeFlows(code Code, count int) FlowBuf {
	b := Make(code, DomainFlow, FlowHeaderSize+count*flowSize)
	return FlowBuf(b)
}

func (b FlowBuf) Num() int {
	return (len(b) - FlowHeaderSize) / flowSize
}

func (b FlowBuf) Get(i int) (id int32, increment int32) {
	flow := b[FlowHeaderSize+i*flowSize:]
	id = int32(binary.LittleEndian.Uint32(flow[flowOffsetID:]))
	increment = int32(binary.LittleEndian.Uint32(flow[flowOffsetIncrement:]))
	return
}

func (b FlowBuf) Set(i int, id int32, increment int32) {
	flow := b[FlowHeaderSize+i*flowSize:]
	binary.LittleEndian.PutUint32(flow[flowOffsetID:], uint32(id))
	binary.LittleEndian.PutUint32(flow[flowOffsetIncrement:], uint32(increment))
}

// DataBuf holds a data packet.
type DataBuf Buf

func MakeData(code Code, id int32, dataSize int) DataBuf {
	b := Make(code, DomainData, DataHeaderSize+dataSize)
	binary.LittleEndian.PutUint32(b[OffsetDataID:], uint32(id))
	return DataBuf(b)
}

func (b DataBuf) ID() int32 {
	return int32(binary.LittleEndian.Uint32(b[OffsetDataID:]))
}

// Note is a value associated with a data packet.  Each service interface
// specifies its semantics separately.
func (b DataBuf) Note() int32 {
	return int32(binary.LittleEndian.Uint32(b[OffsetDataNote:]))
}

// SetNote value.  It defaults to zero.
func (b DataBuf) SetNote(value int32) {
	binary.LittleEndian.PutUint32(b[OffsetDataNote:], uint32(value))
}

func (b DataBuf) Data() []byte {
	return b[DataHeaderSize:]
}

func (b DataBuf) DataLen() int {
	return len(b) - DataHeaderSize
}

func (b DataBuf) Split(dataLen int) (prefix Buf, unused DataBuf) {
	prefix, unusedBuf := Buf(b).Split(DataHeaderSize, DataHeaderSize+dataLen)
	unused = DataBuf(unusedBuf)
	return
}
