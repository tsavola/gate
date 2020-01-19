// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: server/detail/detail.proto

package detail // import "github.com/tsavola/gate/server/detail"

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
	return fileDescriptor_detail_2ad7ef16060d8100, []int{0}
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
	return fileDescriptor_detail_2ad7ef16060d8100, []int{1}
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
	return fileDescriptor_detail_2ad7ef16060d8100, []int{0}
}
func (m *Context) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Context) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Context.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Context) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Context.Merge(dst, src)
}
func (m *Context) XXX_Size() int {
	return m.Size()
}
func (m *Context) XXX_DiscardUnknown() {
	xxx_messageInfo_Context.DiscardUnknown(m)
}

var xxx_messageInfo_Context proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Context)(nil), "server.detail.Context")
	proto.RegisterEnum("server.detail.Iface", Iface_name, Iface_value)
	proto.RegisterEnum("server.detail.Op", Op_name, Op_value)
}
func (m *Context) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Context) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Iface != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintDetail(dAtA, i, uint64(m.Iface))
	}
	if m.Req != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintDetail(dAtA, i, uint64(m.Req))
	}
	if len(m.Addr) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintDetail(dAtA, i, uint64(len(m.Addr)))
		i += copy(dAtA[i:], m.Addr)
	}
	if m.Op != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintDetail(dAtA, i, uint64(m.Op))
	}
	if len(m.Principal) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintDetail(dAtA, i, uint64(len(m.Principal)))
		i += copy(dAtA[i:], m.Principal)
	}
	return i, nil
}

func encodeVarintDetail(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Context) Size() (n int) {
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
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
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
			wire |= (uint64(b) & 0x7F) << shift
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
				m.Iface |= (Iface(b) & 0x7F) << shift
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
				m.Req |= (uint64(b) & 0x7F) << shift
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
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDetail
			}
			postIndex := iNdEx + intStringLen
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
				m.Op |= (Op(b) & 0x7F) << shift
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
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthDetail
			}
			postIndex := iNdEx + intStringLen
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
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthDetail
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowDetail
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
				next, err := skipDetail(dAtA[start:])
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
	ErrInvalidLengthDetail = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDetail   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("server/detail/detail.proto", fileDescriptor_detail_2ad7ef16060d8100) }

var fileDescriptor_detail_2ad7ef16060d8100 = []byte{
	// 455 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0x4f, 0x6e, 0xd3, 0x40,
	0x14, 0xc6, 0xe3, 0xfc, 0x69, 0xc8, 0x4b, 0x93, 0xbc, 0x4c, 0x53, 0xc9, 0x42, 0xc8, 0x0a, 0x48,
	0x48, 0x51, 0x16, 0x8e, 0x04, 0x4b, 0x56, 0xae, 0x3d, 0xa5, 0x56, 0xa7, 0xe3, 0xc8, 0x33, 0xa3,
	0x20, 0x36, 0x95, 0x9b, 0x98, 0x62, 0x29, 0xd4, 0xc6, 0x35, 0x15, 0xc7, 0xe0, 0x00, 0x9c, 0x80,
	0x93, 0x74, 0xc9, 0x11, 0x20, 0x5c, 0x04, 0xd9, 0x71, 0xd3, 0xb8, 0x2b, 0x8f, 0x7e, 0xef, 0xfb,
	0x7e, 0xd6, 0x93, 0x1e, 0x3c, 0xbf, 0x0d, 0xd3, 0xbb, 0x30, 0x9d, 0xad, 0xc2, 0x2c, 0x88, 0xd6,
	0xe5, 0xc7, 0x4c, 0xd2, 0x38, 0x8b, 0x49, 0x6f, 0x3b, 0x33, 0xb7, 0xf0, 0xd5, 0x4f, 0x0d, 0xda,
	0x76, 0x7c, 0x93, 0x85, 0xdf, 0x33, 0x32, 0x85, 0x56, 0xf4, 0x29, 0x58, 0x86, 0xba, 0x36, 0xd6,
	0x26, 0xfd, 0x37, 0x23, 0xb3, 0x12, 0x35, 0xdd, 0x7c, 0xe6, 0x6f, 0x23, 0x04, 0xa1, 0x91, 0x86,
	0x5f, 0xf5, 0xfa, 0x58, 0x9b, 0x34, 0xfd, 0xfc, 0x49, 0x08, 0x34, 0x83, 0xd5, 0x2a, 0xd5, 0x1b,
	0x63, 0x6d, 0xd2, 0xf1, 0x8b, 0x37, 0x79, 0x09, 0xf5, 0x38, 0xd1, 0x9b, 0x85, 0x6e, 0xf8, 0x44,
	0xe7, 0x25, 0x7e, 0x3d, 0x4e, 0xc8, 0x0b, 0xe8, 0x24, 0x69, 0x74, 0xb3, 0x8c, 0x92, 0x60, 0xad,
	0xb7, 0x8a, 0xee, 0x23, 0x98, 0x8e, 0xa0, 0x55, 0xfc, 0x96, 0x74, 0xa1, 0xed, 0xd0, 0x53, 0x4b,
	0x31, 0x89, 0xb5, 0xe9, 0xaf, 0x06, 0xd4, 0xbd, 0x24, 0x67, 0x8a, 0x9f, 0x73, 0x6f, 0xc1, 0xb1,
	0x46, 0x06, 0xd0, 0xbd, 0xf0, 0x1c, 0xc5, 0xe8, 0x25, 0x73, 0x85, 0x44, 0x8d, 0x1c, 0xc1, 0xa0,
	0x04, 0x8e, 0xb7, 0xe0, 0xcc, 0xb3, 0x1c, 0xac, 0x93, 0x21, 0xf4, 0x4a, 0xa8, 0xe6, 0x05, 0x6a,
	0xec, 0x21, 0xe1, 0x29, 0xdf, 0xa6, 0xd8, 0x24, 0x08, 0x87, 0x0f, 0x29, 0xee, 0xd3, 0x53, 0x6c,
	0xe5, 0x76, 0xdb, 0x62, 0xec, 0x92, 0x7e, 0x90, 0x16, 0x97, 0x78, 0xb0, 0x03, 0xa5, 0xa6, 0xbd,
	0x03, 0xa5, 0xe4, 0x59, 0xee, 0x65, 0x96, 0xe2, 0xf6, 0xd9, 0x43, 0xa9, 0xb3, 0x87, 0xca, 0x1a,
	0xec, 0xa1, 0xb2, 0xd8, 0xcd, 0x91, 0xcb, 0x85, 0xb4, 0xb8, 0x5d, 0xee, 0x72, 0x48, 0x46, 0x80,
	0x3b, 0x64, 0x7b, 0x9c, 0x53, 0x5b, 0x62, 0x2f, 0xdf, 0x70, 0x47, 0x85, 0xb4, 0xa4, 0x12, 0xd8,
	0xaf, 0xb4, 0x17, 0x96, 0x2b, 0x71, 0x50, 0x41, 0xe7, 0x2e, 0x63, 0x88, 0x15, 0xa1, 0x50, 0x62,
	0x4e, 0xb9, 0x83, 0xc3, 0x8a, 0xd0, 0xa7, 0x42, 0x5d, 0x50, 0x24, 0xe4, 0x18, 0x86, 0x8f, 0x51,
	0x6e, 0xcd, 0xc5, 0x99, 0x27, 0xf1, 0xa8, 0x92, 0x75, 0x28, 0xa3, 0x92, 0xe2, 0x88, 0x10, 0xe8,
	0xef, 0xc1, 0x13, 0xf5, 0x1e, 0x8f, 0x4f, 0xde, 0xdd, 0xff, 0x35, 0x6a, 0xf7, 0x1b, 0x43, 0xfb,
	0xbd, 0x31, 0xb4, 0x3f, 0x1b, 0x43, 0xfb, 0xf1, 0xcf, 0xa8, 0x7d, 0x7c, 0x7d, 0x1d, 0x65, 0x9f,
	0xbf, 0x5d, 0x99, 0xcb, 0xf8, 0xcb, 0x2c, 0xbb, 0x0d, 0xee, 0xe2, 0x75, 0x30, 0xbb, 0x0e, 0xb2,
	0x70, 0x56, 0x39, 0xdd, 0xab, 0x83, 0xe2, 0x68, 0xdf, 0xfe, 0x0f, 0x00, 0x00, 0xff, 0xff, 0x73,
	0x9c, 0xc4, 0xdb, 0xd2, 0x02, 0x00, 0x00,
}
