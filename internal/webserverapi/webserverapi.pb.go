// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: internal/webserverapi/webserverapi.proto

package webserverapi

import (
	fmt "fmt"
	serverapi "gate.computer/gate/serverapi"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type IOConnection struct {
	Connected            bool     `protobuf:"varint,1,opt,name=connected,proto3" json:"connected,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IOConnection) Reset()         { *m = IOConnection{} }
func (m *IOConnection) String() string { return proto.CompactTextString(m) }
func (*IOConnection) ProtoMessage()    {}
func (*IOConnection) Descriptor() ([]byte, []int) {
	return fileDescriptor_9222a6c75676987e, []int{0}
}
func (m *IOConnection) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IOConnection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IOConnection.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IOConnection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IOConnection.Merge(m, src)
}
func (m *IOConnection) XXX_Size() int {
	return m.Size()
}
func (m *IOConnection) XXX_DiscardUnknown() {
	xxx_messageInfo_IOConnection.DiscardUnknown(m)
}

var xxx_messageInfo_IOConnection proto.InternalMessageInfo

type ConnectionStatus struct {
	Status               serverapi.Status `protobuf:"bytes,1,opt,name=status,proto3" json:"status"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ConnectionStatus) Reset()         { *m = ConnectionStatus{} }
func (m *ConnectionStatus) String() string { return proto.CompactTextString(m) }
func (*ConnectionStatus) ProtoMessage()    {}
func (*ConnectionStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_9222a6c75676987e, []int{1}
}
func (m *ConnectionStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ConnectionStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ConnectionStatus.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ConnectionStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectionStatus.Merge(m, src)
}
func (m *ConnectionStatus) XXX_Size() int {
	return m.Size()
}
func (m *ConnectionStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectionStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectionStatus proto.InternalMessageInfo

func init() {
	proto.RegisterType((*IOConnection)(nil), "webserverapi.IOConnection")
	proto.RegisterType((*ConnectionStatus)(nil), "webserverapi.ConnectionStatus")
}

func init() {
	proto.RegisterFile("internal/webserverapi/webserverapi.proto", fileDescriptor_9222a6c75676987e)
}

var fileDescriptor_9222a6c75676987e = []byte{
	// 187 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0xc8, 0xcc, 0x2b, 0x49,
	0x2d, 0xca, 0x4b, 0xcc, 0xd1, 0x2f, 0x4f, 0x4d, 0x2a, 0x4e, 0x2d, 0x2a, 0x4b, 0x2d, 0x4a, 0x2c,
	0xc8, 0x44, 0xe1, 0xe8, 0x15, 0x14, 0xe5, 0x97, 0xe4, 0x0b, 0xf1, 0x20, 0x8b, 0x49, 0x49, 0x22,
	0xd4, 0xa2, 0x29, 0x54, 0xd2, 0xe1, 0xe2, 0xf1, 0xf4, 0x77, 0xce, 0xcf, 0xcb, 0x4b, 0x4d, 0x2e,
	0xc9, 0xcc, 0xcf, 0x13, 0x92, 0xe1, 0xe2, 0x4c, 0x86, 0xf0, 0x52, 0x53, 0x24, 0x18, 0x15, 0x18,
	0x35, 0x38, 0x82, 0x10, 0x02, 0x4a, 0xce, 0x5c, 0x02, 0x08, 0xb5, 0xc1, 0x25, 0x89, 0x25, 0xa5,
	0xc5, 0x42, 0xfa, 0x5c, 0x6c, 0xc5, 0x60, 0x16, 0x58, 0x39, 0xb7, 0x91, 0xa0, 0x1e, 0xc2, 0x0e,
	0x88, 0x12, 0x27, 0x96, 0x13, 0xf7, 0xe4, 0x19, 0x82, 0xa0, 0xca, 0x9c, 0xec, 0x4e, 0x3c, 0x94,
	0x63, 0x38, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c,
	0x96, 0x63, 0x88, 0xd2, 0x48, 0x4f, 0x2c, 0x49, 0xd5, 0x4b, 0xce, 0xcf, 0x2d, 0x28, 0x2d, 0x49,
	0x2d, 0xd2, 0x07, 0xf1, 0xf4, 0xb1, 0x7a, 0x37, 0x89, 0x0d, 0xec, 0x72, 0x63, 0x40, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xe1, 0x0e, 0xaf, 0xb0, 0x0e, 0x01, 0x00, 0x00,
}

func (m *IOConnection) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IOConnection) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IOConnection) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Connected {
		i--
		if m.Connected {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ConnectionStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ConnectionStatus) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ConnectionStatus) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Status.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintWebserverapi(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintWebserverapi(dAtA []byte, offset int, v uint64) int {
	offset -= sovWebserverapi(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *IOConnection) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Connected {
		n += 2
	}
	return n
}

func (m *ConnectionStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Status.Size()
	n += 1 + l + sovWebserverapi(uint64(l))
	return n
}

func sovWebserverapi(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozWebserverapi(x uint64) (n int) {
	return sovWebserverapi(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *IOConnection) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowWebserverapi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
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
					return ErrIntOverflowWebserverapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Connected = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipWebserverapi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthWebserverapi
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthWebserverapi
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
				return ErrIntOverflowWebserverapi
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
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
					return ErrIntOverflowWebserverapi
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthWebserverapi
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthWebserverapi
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Status.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipWebserverapi(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthWebserverapi
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthWebserverapi
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
func skipWebserverapi(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowWebserverapi
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
					return 0, ErrIntOverflowWebserverapi
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowWebserverapi
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
			if length < 0 {
				return 0, ErrInvalidLengthWebserverapi
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupWebserverapi
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthWebserverapi
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthWebserverapi        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowWebserverapi          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupWebserverapi = fmt.Errorf("proto: unexpected end of group")
)
