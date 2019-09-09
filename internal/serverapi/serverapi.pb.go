// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: internal/serverapi/serverapi.proto

package serverapi // import "github.com/tsavola/gate/internal/serverapi"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type State int32

const (
	State_NONEXISTENT State = 0
	State_RUNNING     State = 1
	State_SUSPENDED   State = 2
	State_TERMINATED  State = 3
	State_KILLED      State = 4
)

var State_name = map[int32]string{
	0: "NONEXISTENT",
	1: "RUNNING",
	2: "SUSPENDED",
	3: "TERMINATED",
	4: "KILLED",
}
var State_value = map[string]int32{
	"NONEXISTENT": 0,
	"RUNNING":     1,
	"SUSPENDED":   2,
	"TERMINATED":  3,
	"KILLED":      4,
}

func (x State) String() string {
	return proto.EnumName(State_name, int32(x))
}
func (State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_serverapi_4494df2946e13c40, []int{0}
}

type Cause int32

const (
	Cause_NORMAL                            Cause = 0
	Cause_NO_FUNCTION                       Cause = 1
	Cause_UNREACHABLE                       Cause = 3
	Cause_CALL_STACK_EXHAUSTED              Cause = 4
	Cause_MEMORY_ACCESS_OUT_OF_BOUNDS       Cause = 5
	Cause_INDIRECT_CALL_INDEX_OUT_OF_BOUNDS Cause = 6
	Cause_INDIRECT_CALL_SIGNATURE_MISMATCH  Cause = 7
	Cause_INTEGER_DIVIDE_BY_ZERO            Cause = 8
	Cause_INTEGER_OVERFLOW                  Cause = 9
	Cause_ABI_VIOLATION                     Cause = 65
)

var Cause_name = map[int32]string{
	0:  "NORMAL",
	1:  "NO_FUNCTION",
	3:  "UNREACHABLE",
	4:  "CALL_STACK_EXHAUSTED",
	5:  "MEMORY_ACCESS_OUT_OF_BOUNDS",
	6:  "INDIRECT_CALL_INDEX_OUT_OF_BOUNDS",
	7:  "INDIRECT_CALL_SIGNATURE_MISMATCH",
	8:  "INTEGER_DIVIDE_BY_ZERO",
	9:  "INTEGER_OVERFLOW",
	65: "ABI_VIOLATION",
}
var Cause_value = map[string]int32{
	"NORMAL":                            0,
	"NO_FUNCTION":                       1,
	"UNREACHABLE":                       3,
	"CALL_STACK_EXHAUSTED":              4,
	"MEMORY_ACCESS_OUT_OF_BOUNDS":       5,
	"INDIRECT_CALL_INDEX_OUT_OF_BOUNDS": 6,
	"INDIRECT_CALL_SIGNATURE_MISMATCH":  7,
	"INTEGER_DIVIDE_BY_ZERO":            8,
	"INTEGER_OVERFLOW":                  9,
	"ABI_VIOLATION":                     65,
}

func (x Cause) String() string {
	return proto.EnumName(Cause_name, int32(x))
}
func (Cause) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_serverapi_4494df2946e13c40, []int{1}
}

type ModuleRef struct {
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (m *ModuleRef) Reset()         { *m = ModuleRef{} }
func (m *ModuleRef) String() string { return proto.CompactTextString(m) }
func (*ModuleRef) ProtoMessage()    {}
func (*ModuleRef) Descriptor() ([]byte, []int) {
	return fileDescriptor_serverapi_4494df2946e13c40, []int{0}
}
func (m *ModuleRef) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ModuleRef) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ModuleRef.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ModuleRef) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ModuleRef.Merge(dst, src)
}
func (m *ModuleRef) XXX_Size() int {
	return m.Size()
}
func (m *ModuleRef) XXX_DiscardUnknown() {
	xxx_messageInfo_ModuleRef.DiscardUnknown(m)
}

var xxx_messageInfo_ModuleRef proto.InternalMessageInfo

type ModuleRefs struct {
	Modules []ModuleRef `protobuf:"bytes,1,rep,name=modules" json:"modules"`
}

func (m *ModuleRefs) Reset()         { *m = ModuleRefs{} }
func (m *ModuleRefs) String() string { return proto.CompactTextString(m) }
func (*ModuleRefs) ProtoMessage()    {}
func (*ModuleRefs) Descriptor() ([]byte, []int) {
	return fileDescriptor_serverapi_4494df2946e13c40, []int{1}
}
func (m *ModuleRefs) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ModuleRefs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ModuleRefs.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ModuleRefs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ModuleRefs.Merge(dst, src)
}
func (m *ModuleRefs) XXX_Size() int {
	return m.Size()
}
func (m *ModuleRefs) XXX_DiscardUnknown() {
	xxx_messageInfo_ModuleRefs.DiscardUnknown(m)
}

var xxx_messageInfo_ModuleRefs proto.InternalMessageInfo

type Status struct {
	State  State  `protobuf:"varint,1,opt,name=state,proto3,enum=server.State" json:"state,omitempty"`
	Cause  Cause  `protobuf:"varint,2,opt,name=cause,proto3,enum=server.Cause" json:"cause,omitempty"`
	Result int32  `protobuf:"varint,3,opt,name=result,proto3" json:"result,omitempty"`
	Error  string `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	Debug  string `protobuf:"bytes,5,opt,name=debug,proto3" json:"debug,omitempty"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_serverapi_4494df2946e13c40, []int{2}
}
func (m *Status) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Status.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(dst, src)
}
func (m *Status) XXX_Size() int {
	return m.Size()
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

type InstanceStatus struct {
	Instance string `protobuf:"bytes,1,opt,name=instance,proto3" json:"instance,omitempty"`
	Status   Status `protobuf:"bytes,2,opt,name=status" json:"status"`
}

func (m *InstanceStatus) Reset()         { *m = InstanceStatus{} }
func (m *InstanceStatus) String() string { return proto.CompactTextString(m) }
func (*InstanceStatus) ProtoMessage()    {}
func (*InstanceStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_serverapi_4494df2946e13c40, []int{3}
}
func (m *InstanceStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InstanceStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InstanceStatus.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *InstanceStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstanceStatus.Merge(dst, src)
}
func (m *InstanceStatus) XXX_Size() int {
	return m.Size()
}
func (m *InstanceStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_InstanceStatus.DiscardUnknown(m)
}

var xxx_messageInfo_InstanceStatus proto.InternalMessageInfo

type Instances struct {
	Instances []InstanceStatus `protobuf:"bytes,1,rep,name=instances" json:"instances"`
}

func (m *Instances) Reset()         { *m = Instances{} }
func (m *Instances) String() string { return proto.CompactTextString(m) }
func (*Instances) ProtoMessage()    {}
func (*Instances) Descriptor() ([]byte, []int) {
	return fileDescriptor_serverapi_4494df2946e13c40, []int{4}
}
func (m *Instances) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Instances) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Instances.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Instances) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Instances.Merge(dst, src)
}
func (m *Instances) XXX_Size() int {
	return m.Size()
}
func (m *Instances) XXX_DiscardUnknown() {
	xxx_messageInfo_Instances.DiscardUnknown(m)
}

var xxx_messageInfo_Instances proto.InternalMessageInfo

type IOConnection struct {
	Connected bool   `protobuf:"varint,1,opt,name=connected,proto3" json:"connected,omitempty"`
	Status    Status `protobuf:"bytes,2,opt,name=status" json:"status"`
}

func (m *IOConnection) Reset()         { *m = IOConnection{} }
func (m *IOConnection) String() string { return proto.CompactTextString(m) }
func (*IOConnection) ProtoMessage()    {}
func (*IOConnection) Descriptor() ([]byte, []int) {
	return fileDescriptor_serverapi_4494df2946e13c40, []int{5}
}
func (m *IOConnection) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IOConnection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IOConnection.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *IOConnection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IOConnection.Merge(dst, src)
}
func (m *IOConnection) XXX_Size() int {
	return m.Size()
}
func (m *IOConnection) XXX_DiscardUnknown() {
	xxx_messageInfo_IOConnection.DiscardUnknown(m)
}

var xxx_messageInfo_IOConnection proto.InternalMessageInfo

type ConnectionStatus struct {
	Status Status `protobuf:"bytes,1,opt,name=status" json:"status"`
}

func (m *ConnectionStatus) Reset()         { *m = ConnectionStatus{} }
func (m *ConnectionStatus) String() string { return proto.CompactTextString(m) }
func (*ConnectionStatus) ProtoMessage()    {}
func (*ConnectionStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_serverapi_4494df2946e13c40, []int{6}
}
func (m *ConnectionStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ConnectionStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ConnectionStatus.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ConnectionStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectionStatus.Merge(dst, src)
}
func (m *ConnectionStatus) XXX_Size() int {
	return m.Size()
}
func (m *ConnectionStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectionStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectionStatus proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ModuleRef)(nil), "server.ModuleRef")
	proto.RegisterType((*ModuleRefs)(nil), "server.ModuleRefs")
	proto.RegisterType((*Status)(nil), "server.Status")
	proto.RegisterType((*InstanceStatus)(nil), "server.InstanceStatus")
	proto.RegisterType((*Instances)(nil), "server.Instances")
	proto.RegisterType((*IOConnection)(nil), "server.IOConnection")
	proto.RegisterType((*ConnectionStatus)(nil), "server.ConnectionStatus")
	proto.RegisterEnum("server.State", State_name, State_value)
	proto.RegisterEnum("server.Cause", Cause_name, Cause_value)
}
func (m *ModuleRef) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ModuleRef) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Key) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintServerapi(dAtA, i, uint64(len(m.Key)))
		i += copy(dAtA[i:], m.Key)
	}
	return i, nil
}

func (m *ModuleRefs) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ModuleRefs) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Modules) > 0 {
		for _, msg := range m.Modules {
			dAtA[i] = 0xa
			i++
			i = encodeVarintServerapi(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Status) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Status) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.State != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintServerapi(dAtA, i, uint64(m.State))
	}
	if m.Cause != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintServerapi(dAtA, i, uint64(m.Cause))
	}
	if m.Result != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintServerapi(dAtA, i, uint64(m.Result))
	}
	if len(m.Error) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintServerapi(dAtA, i, uint64(len(m.Error)))
		i += copy(dAtA[i:], m.Error)
	}
	if len(m.Debug) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintServerapi(dAtA, i, uint64(len(m.Debug)))
		i += copy(dAtA[i:], m.Debug)
	}
	return i, nil
}

func (m *InstanceStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InstanceStatus) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Instance) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintServerapi(dAtA, i, uint64(len(m.Instance)))
		i += copy(dAtA[i:], m.Instance)
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintServerapi(dAtA, i, uint64(m.Status.Size()))
	n1, err := m.Status.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	return i, nil
}

func (m *Instances) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Instances) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Instances) > 0 {
		for _, msg := range m.Instances {
			dAtA[i] = 0xa
			i++
			i = encodeVarintServerapi(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *IOConnection) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IOConnection) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Connected {
		dAtA[i] = 0x8
		i++
		if m.Connected {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintServerapi(dAtA, i, uint64(m.Status.Size()))
	n2, err := m.Status.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	return i, nil
}

func (m *ConnectionStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ConnectionStatus) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintServerapi(dAtA, i, uint64(m.Status.Size()))
	n3, err := m.Status.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	return i, nil
}

func encodeVarintServerapi(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ModuleRef) Size() (n int) {
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovServerapi(uint64(l))
	}
	return n
}

func (m *ModuleRefs) Size() (n int) {
	var l int
	_ = l
	if len(m.Modules) > 0 {
		for _, e := range m.Modules {
			l = e.Size()
			n += 1 + l + sovServerapi(uint64(l))
		}
	}
	return n
}

func (m *Status) Size() (n int) {
	var l int
	_ = l
	if m.State != 0 {
		n += 1 + sovServerapi(uint64(m.State))
	}
	if m.Cause != 0 {
		n += 1 + sovServerapi(uint64(m.Cause))
	}
	if m.Result != 0 {
		n += 1 + sovServerapi(uint64(m.Result))
	}
	l = len(m.Error)
	if l > 0 {
		n += 1 + l + sovServerapi(uint64(l))
	}
	l = len(m.Debug)
	if l > 0 {
		n += 1 + l + sovServerapi(uint64(l))
	}
	return n
}

func (m *InstanceStatus) Size() (n int) {
	var l int
	_ = l
	l = len(m.Instance)
	if l > 0 {
		n += 1 + l + sovServerapi(uint64(l))
	}
	l = m.Status.Size()
	n += 1 + l + sovServerapi(uint64(l))
	return n
}

func (m *Instances) Size() (n int) {
	var l int
	_ = l
	if len(m.Instances) > 0 {
		for _, e := range m.Instances {
			l = e.Size()
			n += 1 + l + sovServerapi(uint64(l))
		}
	}
	return n
}

func (m *IOConnection) Size() (n int) {
	var l int
	_ = l
	if m.Connected {
		n += 2
	}
	l = m.Status.Size()
	n += 1 + l + sovServerapi(uint64(l))
	return n
}

func (m *ConnectionStatus) Size() (n int) {
	var l int
	_ = l
	l = m.Status.Size()
	n += 1 + l + sovServerapi(uint64(l))
	return n
}

func sovServerapi(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozServerapi(x uint64) (n int) {
	return sovServerapi(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ModuleRef) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerapi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ModuleRef: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ModuleRef: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerapi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerapi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerapi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ModuleRefs) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerapi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ModuleRefs: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ModuleRefs: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Modules", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerapi
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Modules = append(m.Modules, ModuleRef{})
			if err := m.Modules[len(m.Modules)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerapi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerapi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Status) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerapi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Status: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Status: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= (State(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cause", wireType)
			}
			m.Cause = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Cause |= (Cause(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Result", wireType)
			}
			m.Result = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Result |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerapi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Error = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Debug", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerapi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Debug = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerapi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerapi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *InstanceStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerapi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: InstanceStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InstanceStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Instance", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServerapi
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Instance = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerapi
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Status.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerapi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerapi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Instances) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerapi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Instances: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Instances: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Instances", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerapi
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Instances = append(m.Instances, InstanceStatus{})
			if err := m.Instances[len(m.Instances)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerapi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerapi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *IOConnection) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerapi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IOConnection: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IOConnection: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Connected", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Connected = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerapi
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Status.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerapi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerapi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ConnectionStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServerapi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ConnectionStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ConnectionStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthServerapi
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Status.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServerapi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServerapi
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipServerapi(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowServerapi
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowServerapi
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthServerapi
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowServerapi
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipServerapi(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthServerapi = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowServerapi   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("internal/serverapi/serverapi.proto", fileDescriptor_serverapi_4494df2946e13c40)
}

var fileDescriptor_serverapi_4494df2946e13c40 = []byte{
	// 634 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0xcb, 0x6e, 0xd3, 0x4c,
	0x1c, 0xc5, 0xe3, 0xe6, 0xd2, 0xe6, 0x9f, 0xaf, 0xf9, 0xa6, 0xa3, 0xaa, 0x8a, 0x0a, 0xa4, 0xc1,
	0x80, 0x54, 0x55, 0x28, 0x11, 0x65, 0xc7, 0x06, 0x1c, 0x7b, 0x9a, 0x8c, 0x6a, 0x8f, 0xab, 0xb1,
	0x5d, 0xda, 0x6c, 0x2c, 0x37, 0x1d, 0x4a, 0x44, 0x6a, 0x57, 0xbe, 0x54, 0x62, 0xcb, 0x13, 0xb0,
	0xe2, 0x99, 0xba, 0xec, 0x92, 0x15, 0x82, 0xf6, 0x45, 0x90, 0x6f, 0x0d, 0x85, 0x05, 0x62, 0x37,
	0xe7, 0xfc, 0xcf, 0xfc, 0x74, 0x66, 0x46, 0x03, 0xf2, 0xcc, 0x8f, 0x45, 0xe8, 0x7b, 0xf3, 0x41,
	0x24, 0xc2, 0x4b, 0x11, 0x7a, 0x17, 0xb3, 0xc5, 0xaa, 0x7f, 0x11, 0x06, 0x71, 0x80, 0x1b, 0xb9,
	0x21, 0x3f, 0x82, 0xa6, 0x11, 0x9c, 0x26, 0x73, 0xc1, 0xc5, 0x3b, 0x8c, 0xa0, 0xfa, 0x41, 0x7c,
	0xec, 0x48, 0x3d, 0x69, 0xbb, 0xc9, 0xd3, 0xa5, 0xfc, 0x1a, 0xe0, 0x6e, 0x1c, 0xe1, 0x17, 0xb0,
	0x7c, 0x9e, 0xa9, 0xa8, 0x23, 0xf5, 0xaa, 0xdb, 0xad, 0xdd, 0xb5, 0x7e, 0x8e, 0xe9, 0xdf, 0x85,
	0x86, 0xb5, 0xab, 0x6f, 0x5b, 0x15, 0x5e, 0xe6, 0xe4, 0x2f, 0x12, 0x34, 0xac, 0xd8, 0x8b, 0x93,
	0x08, 0x3f, 0x81, 0x7a, 0x14, 0x7b, 0xb1, 0xc8, 0xf8, 0xed, 0xdd, 0xd5, 0x72, 0x6f, 0x3a, 0x16,
	0x3c, 0x9f, 0xa5, 0xa1, 0xa9, 0x97, 0x44, 0xa2, 0xb3, 0x74, 0x3f, 0xa4, 0xa6, 0x26, 0xcf, 0x67,
	0x78, 0x03, 0x1a, 0xa1, 0x88, 0x92, 0x79, 0xdc, 0xa9, 0xf6, 0xa4, 0xed, 0x3a, 0x2f, 0x14, 0x5e,
	0x87, 0xba, 0x08, 0xc3, 0x20, 0xec, 0xd4, 0xb2, 0x13, 0xe4, 0x22, 0x75, 0x4f, 0xc5, 0x49, 0x72,
	0xd6, 0xa9, 0xe7, 0x6e, 0x26, 0xe4, 0x09, 0xb4, 0xa9, 0x1f, 0xc5, 0x9e, 0x3f, 0x15, 0x45, 0xbf,
	0x4d, 0x58, 0x99, 0x15, 0x4e, 0x71, 0x05, 0x77, 0x1a, 0x3f, 0x87, 0x46, 0x94, 0xa5, 0xb2, 0x5e,
	0xad, 0xdd, 0xf6, 0xaf, 0xe5, 0x93, 0xa8, 0x38, 0x75, 0x91, 0x91, 0x47, 0xd0, 0x2c, 0xd9, 0x11,
	0x7e, 0x05, 0xcd, 0x12, 0x53, 0x5e, 0xdb, 0x46, 0xb9, 0xfb, 0x7e, 0x83, 0x82, 0xb2, 0x88, 0xcb,
	0x13, 0xf8, 0x8f, 0x9a, 0x6a, 0xe0, 0xfb, 0x62, 0x1a, 0xcf, 0x02, 0x1f, 0x3f, 0x84, 0xe6, 0x34,
	0x57, 0xe2, 0x34, 0xeb, 0xb8, 0xc2, 0x17, 0xc6, 0x3f, 0x96, 0x7c, 0x03, 0x68, 0x41, 0x2e, 0xae,
	0x60, 0x41, 0x90, 0xfe, 0x4e, 0xd8, 0x39, 0x80, 0x7a, 0xf6, 0x76, 0xf8, 0x7f, 0x68, 0x31, 0x93,
	0x91, 0x23, 0x6a, 0xd9, 0x84, 0xd9, 0xa8, 0x82, 0x5b, 0xb0, 0xcc, 0x1d, 0xc6, 0x28, 0x1b, 0x21,
	0x09, 0xaf, 0x42, 0xd3, 0x72, 0xac, 0x03, 0xc2, 0x34, 0xa2, 0xa1, 0x25, 0xdc, 0x06, 0xb0, 0x09,
	0x37, 0x28, 0x53, 0x6c, 0xa2, 0xa1, 0x2a, 0x06, 0x68, 0xec, 0x53, 0x5d, 0x27, 0x1a, 0xaa, 0xed,
	0x7c, 0x5a, 0x82, 0x7a, 0xf6, 0xd2, 0xa9, 0xcb, 0x4c, 0x6e, 0x28, 0x3a, 0xaa, 0xe4, 0x78, 0x77,
	0xcf, 0x61, 0xaa, 0x4d, 0x4d, 0x86, 0xa4, 0xd4, 0x70, 0x18, 0x27, 0x8a, 0x3a, 0x56, 0x86, 0x3a,
	0x41, 0x55, 0xdc, 0x81, 0x75, 0x55, 0xd1, 0x75, 0xd7, 0xb2, 0x15, 0x75, 0xdf, 0x25, 0x47, 0x63,
	0xc5, 0xb1, 0x52, 0x7a, 0x0d, 0x6f, 0xc1, 0x03, 0x83, 0x18, 0x26, 0x3f, 0x76, 0x15, 0x55, 0x25,
	0x96, 0xe5, 0x9a, 0x8e, 0xed, 0x9a, 0x7b, 0xee, 0xd0, 0x74, 0x98, 0x66, 0xa1, 0x3a, 0x7e, 0x06,
	0x8f, 0x29, 0xd3, 0x28, 0x27, 0xaa, 0xed, 0x66, 0x0c, 0xca, 0x34, 0x72, 0xf4, 0x5b, 0xac, 0x81,
	0x9f, 0x42, 0xef, 0x7e, 0xcc, 0xa2, 0x23, 0xa6, 0xd8, 0x0e, 0x27, 0xae, 0x41, 0x2d, 0x43, 0xb1,
	0xd5, 0x31, 0x5a, 0xc6, 0x9b, 0xb0, 0x41, 0x99, 0x4d, 0x46, 0x84, 0xbb, 0x1a, 0x3d, 0xa4, 0x1a,
	0x71, 0x87, 0xc7, 0xee, 0x84, 0x70, 0x13, 0xad, 0xe0, 0x75, 0x40, 0xe5, 0xcc, 0x3c, 0x24, 0x7c,
	0x4f, 0x37, 0xdf, 0xa2, 0x26, 0x5e, 0x83, 0x55, 0x65, 0x48, 0xdd, 0x43, 0x6a, 0xea, 0x4a, 0x76,
	0x3a, 0x65, 0x38, 0xbe, 0xfa, 0xd1, 0xad, 0x5c, 0xdd, 0x74, 0xa5, 0xeb, 0x9b, 0xae, 0xf4, 0xfd,
	0xa6, 0x2b, 0x7d, 0xbe, 0xed, 0x56, 0xae, 0x6f, 0xbb, 0x95, 0xaf, 0xb7, 0xdd, 0xca, 0x64, 0xe7,
	0x6c, 0x16, 0xbf, 0x4f, 0x4e, 0xfa, 0xd3, 0xe0, 0x7c, 0x10, 0x47, 0xde, 0x65, 0x30, 0xf7, 0x06,
	0x67, 0x5e, 0x2c, 0x06, 0x7f, 0x7e, 0xfa, 0x93, 0x46, 0xf6, 0xd7, 0x5f, 0xfe, 0x0c, 0x00, 0x00,
	0xff, 0xff, 0x24, 0xce, 0xb5, 0x0e, 0x11, 0x04, 0x00, 0x00,
}
