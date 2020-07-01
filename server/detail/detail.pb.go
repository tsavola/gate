// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: server/detail/detail.proto

package detail

import (
	fmt "fmt"
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

type Iface int32

const (
	Iface_DEFAULT Iface = 0
)

var Iface_name = map[int32]string{
	0: "DEFAULT",
}

var Iface_value = map[string]int32{
	"DEFAULT": 0,
}

func (x Iface) String() string {
	return proto.EnumName(Iface_name, int32(x))
}

func (Iface) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ceb992a7f462a1c4, []int{0}
}

type Op int32

const (
	Op_UNKNOWN           Op = 0
	Op_MODULE_LIST       Op = 1
	Op_MODULE_DOWNLOAD   Op = 2
	Op_MODULE_UPLOAD     Op = 3
	Op_MODULE_SOURCE     Op = 4
	Op_MODULE_UNREF      Op = 5
	Op_CALL_EXTANT       Op = 6
	Op_CALL_UPLOAD       Op = 7
	Op_CALL_SOURCE       Op = 8
	Op_LAUNCH_EXTANT     Op = 9
	Op_LAUNCH_UPLOAD     Op = 10
	Op_LAUNCH_SOURCE     Op = 11
	Op_INSTANCE_LIST     Op = 12
	Op_INSTANCE_CONNECT  Op = 13
	Op_INSTANCE_STATUS   Op = 14
	Op_INSTANCE_WAIT     Op = 15
	Op_INSTANCE_KILL     Op = 16
	Op_INSTANCE_SUSPEND  Op = 17
	Op_INSTANCE_RESUME   Op = 18
	Op_INSTANCE_SNAPSHOT Op = 19
	Op_INSTANCE_DELETE   Op = 20
	Op_INSTANCE_DEBUG    Op = 21
)

var Op_name = map[int32]string{
	0:  "UNKNOWN",
	1:  "MODULE_LIST",
	2:  "MODULE_DOWNLOAD",
	3:  "MODULE_UPLOAD",
	4:  "MODULE_SOURCE",
	5:  "MODULE_UNREF",
	6:  "CALL_EXTANT",
	7:  "CALL_UPLOAD",
	8:  "CALL_SOURCE",
	9:  "LAUNCH_EXTANT",
	10: "LAUNCH_UPLOAD",
	11: "LAUNCH_SOURCE",
	12: "INSTANCE_LIST",
	13: "INSTANCE_CONNECT",
	14: "INSTANCE_STATUS",
	15: "INSTANCE_WAIT",
	16: "INSTANCE_KILL",
	17: "INSTANCE_SUSPEND",
	18: "INSTANCE_RESUME",
	19: "INSTANCE_SNAPSHOT",
	20: "INSTANCE_DELETE",
	21: "INSTANCE_DEBUG",
}

var Op_value = map[string]int32{
	"UNKNOWN":           0,
	"MODULE_LIST":       1,
	"MODULE_DOWNLOAD":   2,
	"MODULE_UPLOAD":     3,
	"MODULE_SOURCE":     4,
	"MODULE_UNREF":      5,
	"CALL_EXTANT":       6,
	"CALL_UPLOAD":       7,
	"CALL_SOURCE":       8,
	"LAUNCH_EXTANT":     9,
	"LAUNCH_UPLOAD":     10,
	"LAUNCH_SOURCE":     11,
	"INSTANCE_LIST":     12,
	"INSTANCE_CONNECT":  13,
	"INSTANCE_STATUS":   14,
	"INSTANCE_WAIT":     15,
	"INSTANCE_KILL":     16,
	"INSTANCE_SUSPEND":  17,
	"INSTANCE_RESUME":   18,
	"INSTANCE_SNAPSHOT": 19,
	"INSTANCE_DELETE":   20,
	"INSTANCE_DEBUG":    21,
}

func (x Op) String() string {
	return proto.EnumName(Op_name, int32(x))
}

func (Op) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ceb992a7f462a1c4, []int{1}
}

type Context struct {
	Iface                Iface    `protobuf:"varint,1,opt,name=iface,proto3,enum=server.detail.Iface" json:"iface,omitempty"`
	Req                  uint64   `protobuf:"varint,2,opt,name=req,proto3" json:"req,omitempty"`
	Addr                 string   `protobuf:"bytes,3,opt,name=addr,proto3" json:"addr,omitempty"`
	Op                   Op       `protobuf:"varint,4,opt,name=op,proto3,enum=server.detail.Op" json:"op,omitempty"`
	Principal            string   `protobuf:"bytes,5,opt,name=principal,proto3" json:"principal,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Context) Reset()         { *m = Context{} }
func (m *Context) String() string { return proto.CompactTextString(m) }
func (*Context) ProtoMessage()    {}
func (*Context) Descriptor() ([]byte, []int) {
	return fileDescriptor_ceb992a7f462a1c4, []int{0}
}
func (m *Context) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Context) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Context.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Context) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Context.Merge(m, src)
}
func (m *Context) XXX_Size() int {
	return m.Size()
}
func (m *Context) XXX_DiscardUnknown() {
	xxx_messageInfo_Context.DiscardUnknown(m)
}

var xxx_messageInfo_Context proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("server.detail.Iface", Iface_name, Iface_value)
	proto.RegisterEnum("server.detail.Op", Op_name, Op_value)
	proto.RegisterType((*Context)(nil), "server.detail.Context")
}

func init() { proto.RegisterFile("server/detail/detail.proto", fileDescriptor_ceb992a7f462a1c4) }

var fileDescriptor_ceb992a7f462a1c4 = []byte{
	// 446 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0x41, 0x6e, 0xd3, 0x40,
	0x14, 0x86, 0x63, 0x27, 0x69, 0xc8, 0x4b, 0x93, 0xbc, 0x4c, 0x53, 0xc9, 0x42, 0xc8, 0x0a, 0xac,
	0xa2, 0x2c, 0x1c, 0x09, 0x24, 0xf6, 0xae, 0x3d, 0xa5, 0x56, 0xa7, 0x33, 0x91, 0x67, 0x46, 0x41,
	0x6c, 0x2a, 0x93, 0x18, 0x14, 0xa9, 0xd4, 0xc6, 0x18, 0xc4, 0x31, 0x38, 0x00, 0x27, 0xe0, 0x24,
	0x5d, 0x72, 0x04, 0x08, 0x17, 0x41, 0x76, 0xdc, 0x34, 0xee, 0xca, 0xe3, 0x6f, 0xfe, 0xff, 0x1b,
	0x3d, 0xe9, 0xc1, 0xd3, 0x2f, 0x71, 0xf6, 0x2d, 0xce, 0xe6, 0xeb, 0x38, 0x8f, 0x36, 0x37, 0xd5,
	0xc7, 0x49, 0xb3, 0x24, 0x4f, 0x48, 0x7f, 0x77, 0xe7, 0xec, 0xe0, 0x8b, 0x9f, 0x06, 0x74, 0xbc,
	0xe4, 0x36, 0x8f, 0xbf, 0xe7, 0x64, 0x06, 0xed, 0xcd, 0x87, 0x68, 0x15, 0x5b, 0xc6, 0xc4, 0x98,
	0x0e, 0x5e, 0x8e, 0x9d, 0x5a, 0xd4, 0x09, 0x8a, 0xbb, 0x70, 0x17, 0x21, 0x08, 0xcd, 0x2c, 0xfe,
	0x6c, 0x99, 0x13, 0x63, 0xda, 0x0a, 0x8b, 0x23, 0x21, 0xd0, 0x8a, 0xd6, 0xeb, 0xcc, 0x6a, 0x4e,
	0x8c, 0x69, 0x37, 0x2c, 0xcf, 0xe4, 0x39, 0x98, 0x49, 0x6a, 0xb5, 0x4a, 0xdd, 0xe8, 0x91, 0x4e,
	0xa4, 0xa1, 0x99, 0xa4, 0xe4, 0x19, 0x74, 0xd3, 0x6c, 0x73, 0xbb, 0xda, 0xa4, 0xd1, 0x8d, 0xd5,
	0x2e, 0xbb, 0x0f, 0x60, 0x36, 0x86, 0x76, 0xf9, 0x2c, 0xe9, 0x41, 0xc7, 0xa7, 0xe7, 0xae, 0x66,
	0x0a, 0x1b, 0xb3, 0x5f, 0x4d, 0x30, 0x45, 0x5a, 0x30, 0xcd, 0x2f, 0xb9, 0x58, 0x72, 0x6c, 0x90,
	0x21, 0xf4, 0xae, 0x84, 0xaf, 0x19, 0xbd, 0x66, 0x81, 0x54, 0x68, 0x90, 0x13, 0x18, 0x56, 0xc0,
	0x17, 0x4b, 0xce, 0x84, 0xeb, 0xa3, 0x49, 0x46, 0xd0, 0xaf, 0xa0, 0x5e, 0x94, 0xa8, 0x79, 0x80,
	0xa4, 0xd0, 0xa1, 0x47, 0xb1, 0x45, 0x10, 0x8e, 0xef, 0x53, 0x3c, 0xa4, 0xe7, 0xd8, 0x2e, 0xec,
	0x9e, 0xcb, 0xd8, 0x35, 0x7d, 0xab, 0x5c, 0xae, 0xf0, 0x68, 0x0f, 0x2a, 0x4d, 0x67, 0x0f, 0x2a,
	0xc9, 0x93, 0xc2, 0xcb, 0x5c, 0xcd, 0xbd, 0x8b, 0xfb, 0x52, 0xf7, 0x00, 0x55, 0x35, 0x38, 0x40,
	0x55, 0xb1, 0x57, 0xa0, 0x80, 0x4b, 0xe5, 0x72, 0xaf, 0x9a, 0xe5, 0x98, 0x8c, 0x01, 0xf7, 0xc8,
	0x13, 0x9c, 0x53, 0x4f, 0x61, 0xbf, 0x98, 0x70, 0x4f, 0xa5, 0x72, 0x95, 0x96, 0x38, 0xa8, 0xb5,
	0x97, 0x6e, 0xa0, 0x70, 0x58, 0x43, 0x97, 0x01, 0x63, 0x88, 0x35, 0xa1, 0xd4, 0x72, 0x41, 0xb9,
	0x8f, 0xa3, 0x9a, 0x30, 0xa4, 0x52, 0x5f, 0x51, 0x24, 0xe4, 0x14, 0x46, 0x0f, 0x51, 0xee, 0x2e,
	0xe4, 0x85, 0x50, 0x78, 0x52, 0xcb, 0xfa, 0x94, 0x51, 0x45, 0x71, 0x4c, 0x08, 0x0c, 0x0e, 0xe0,
	0x99, 0x7e, 0x83, 0xa7, 0x67, 0xaf, 0xef, 0xfe, 0xda, 0x8d, 0xbb, 0xad, 0x6d, 0xfc, 0xde, 0xda,
	0xc6, 0x9f, 0xad, 0x6d, 0xfc, 0xf8, 0x67, 0x37, 0xde, 0x4d, 0x3e, 0x46, 0x79, 0xec, 0xac, 0x92,
	0x4f, 0xe9, 0xd7, 0x3c, 0xce, 0xe6, 0xc5, 0xdf, 0xbc, 0xb6, 0xb5, 0xef, 0x8f, 0xca, 0x7d, 0x7d,
	0xf5, 0x3f, 0x00, 0x00, 0xff, 0xff, 0xf4, 0x1a, 0x1f, 0x8a, 0xcd, 0x02, 0x00, 0x00,
}

func (m *Context) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Context) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Context) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Principal) > 0 {
		i -= len(m.Principal)
		copy(dAtA[i:], m.Principal)
		i = encodeVarintDetail(dAtA, i, uint64(len(m.Principal)))
		i--
		dAtA[i] = 0x2a
	}
	if m.Op != 0 {
		i = encodeVarintDetail(dAtA, i, uint64(m.Op))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Addr) > 0 {
		i -= len(m.Addr)
		copy(dAtA[i:], m.Addr)
		i = encodeVarintDetail(dAtA, i, uint64(len(m.Addr)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Req != 0 {
		i = encodeVarintDetail(dAtA, i, uint64(m.Req))
		i--
		dAtA[i] = 0x10
	}
	if m.Iface != 0 {
		i = encodeVarintDetail(dAtA, i, uint64(m.Iface))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintDetail(dAtA []byte, offset int, v uint64) int {
	offset -= sovDetail(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Context) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Iface != 0 {
		n += 1 + sovDetail(uint64(m.Iface))
	}
	if m.Req != 0 {
		n += 1 + sovDetail(uint64(m.Req))
	}
	l = len(m.Addr)
	if l > 0 {
		n += 1 + l + sovDetail(uint64(l))
	}
	if m.Op != 0 {
		n += 1 + sovDetail(uint64(m.Op))
	}
	l = len(m.Principal)
	if l > 0 {
		n += 1 + l + sovDetail(uint64(l))
	}
	return n
}

func sovDetail(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDetail(x uint64) (n int) {
	return sovDetail(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Context) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDetail
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
			return fmt.Errorf("proto: Context: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Context: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Iface", wireType)
			}
			m.Iface = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDetail
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Iface |= Iface(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Req", wireType)
			}
			m.Req = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDetail
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Req |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDetail
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDetail
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDetail
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Op", wireType)
			}
			m.Op = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDetail
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Op |= Op(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Principal", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDetail
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDetail
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDetail
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Principal = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDetail(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDetail
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthDetail
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
func skipDetail(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDetail
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
					return 0, ErrIntOverflowDetail
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
					return 0, ErrIntOverflowDetail
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
				return 0, ErrInvalidLengthDetail
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDetail
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDetail
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDetail        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDetail          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDetail = fmt.Errorf("proto: unexpected end of group")
)
