// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.6.1
// source: server/api/server.proto

package api

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type State int32

const (
	State_NONEXISTENT State = 0
	State_RUNNING     State = 1
	State_SUSPENDED   State = 2
	State_HALTED      State = 3
	State_TERMINATED  State = 4
	State_KILLED      State = 5
)

// Enum value maps for State.
var (
	State_name = map[int32]string{
		0: "NONEXISTENT",
		1: "RUNNING",
		2: "SUSPENDED",
		3: "HALTED",
		4: "TERMINATED",
		5: "KILLED",
	}
	State_value = map[string]int32{
		"NONEXISTENT": 0,
		"RUNNING":     1,
		"SUSPENDED":   2,
		"HALTED":      3,
		"TERMINATED":  4,
		"KILLED":      5,
	}
)

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}

func (x State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (State) Descriptor() protoreflect.EnumDescriptor {
	return file_server_api_server_proto_enumTypes[0].Descriptor()
}

func (State) Type() protoreflect.EnumType {
	return &file_server_api_server_proto_enumTypes[0]
}

func (x State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use State.Descriptor instead.
func (State) EnumDescriptor() ([]byte, []int) {
	return file_server_api_server_proto_rawDescGZIP(), []int{0}
}

type Cause int32

const (
	Cause_NORMAL                            Cause = 0
	Cause_UNREACHABLE                       Cause = 3
	Cause_CALL_STACK_EXHAUSTED              Cause = 4
	Cause_MEMORY_ACCESS_OUT_OF_BOUNDS       Cause = 5
	Cause_INDIRECT_CALL_INDEX_OUT_OF_BOUNDS Cause = 6
	Cause_INDIRECT_CALL_SIGNATURE_MISMATCH  Cause = 7
	Cause_INTEGER_DIVIDE_BY_ZERO            Cause = 8
	Cause_INTEGER_OVERFLOW                  Cause = 9
	Cause_BREAKPOINT                        Cause = 10
	Cause_ABI_DEFICIENCY                    Cause = 27
	Cause_ABI_VIOLATION                     Cause = 28
	Cause_INTERNAL                          Cause = 29
)

// Enum value maps for Cause.
var (
	Cause_name = map[int32]string{
		0:  "NORMAL",
		3:  "UNREACHABLE",
		4:  "CALL_STACK_EXHAUSTED",
		5:  "MEMORY_ACCESS_OUT_OF_BOUNDS",
		6:  "INDIRECT_CALL_INDEX_OUT_OF_BOUNDS",
		7:  "INDIRECT_CALL_SIGNATURE_MISMATCH",
		8:  "INTEGER_DIVIDE_BY_ZERO",
		9:  "INTEGER_OVERFLOW",
		10: "BREAKPOINT",
		27: "ABI_DEFICIENCY",
		28: "ABI_VIOLATION",
		29: "INTERNAL",
	}
	Cause_value = map[string]int32{
		"NORMAL":                            0,
		"UNREACHABLE":                       3,
		"CALL_STACK_EXHAUSTED":              4,
		"MEMORY_ACCESS_OUT_OF_BOUNDS":       5,
		"INDIRECT_CALL_INDEX_OUT_OF_BOUNDS": 6,
		"INDIRECT_CALL_SIGNATURE_MISMATCH":  7,
		"INTEGER_DIVIDE_BY_ZERO":            8,
		"INTEGER_OVERFLOW":                  9,
		"BREAKPOINT":                        10,
		"ABI_DEFICIENCY":                    27,
		"ABI_VIOLATION":                     28,
		"INTERNAL":                          29,
	}
)

func (x Cause) Enum() *Cause {
	p := new(Cause)
	*p = x
	return p
}

func (x Cause) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Cause) Descriptor() protoreflect.EnumDescriptor {
	return file_server_api_server_proto_enumTypes[1].Descriptor()
}

func (Cause) Type() protoreflect.EnumType {
	return &file_server_api_server_proto_enumTypes[1]
}

func (x Cause) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Cause.Descriptor instead.
func (Cause) EnumDescriptor() ([]byte, []int) {
	return file_server_api_server_proto_rawDescGZIP(), []int{1}
}

type DebugOp int32

const (
	DebugOp_CONFIG_GET        DebugOp = 0
	DebugOp_CONFIG_SET        DebugOp = 1
	DebugOp_CONFIG_UNION      DebugOp = 2
	DebugOp_CONFIG_COMPLEMENT DebugOp = 3
	DebugOp_READ_GLOBALS      DebugOp = 4
	DebugOp_READ_MEMORY       DebugOp = 5
	DebugOp_READ_STACK        DebugOp = 6
)

// Enum value maps for DebugOp.
var (
	DebugOp_name = map[int32]string{
		0: "CONFIG_GET",
		1: "CONFIG_SET",
		2: "CONFIG_UNION",
		3: "CONFIG_COMPLEMENT",
		4: "READ_GLOBALS",
		5: "READ_MEMORY",
		6: "READ_STACK",
	}
	DebugOp_value = map[string]int32{
		"CONFIG_GET":        0,
		"CONFIG_SET":        1,
		"CONFIG_UNION":      2,
		"CONFIG_COMPLEMENT": 3,
		"READ_GLOBALS":      4,
		"READ_MEMORY":       5,
		"READ_STACK":        6,
	}
)

func (x DebugOp) Enum() *DebugOp {
	p := new(DebugOp)
	*p = x
	return p
}

func (x DebugOp) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DebugOp) Descriptor() protoreflect.EnumDescriptor {
	return file_server_api_server_proto_enumTypes[2].Descriptor()
}

func (DebugOp) Type() protoreflect.EnumType {
	return &file_server_api_server_proto_enumTypes[2]
}

func (x DebugOp) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DebugOp.Descriptor instead.
func (DebugOp) EnumDescriptor() ([]byte, []int) {
	return file_server_api_server_proto_rawDescGZIP(), []int{2}
}

type ModuleRef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ModuleRef) Reset() {
	*x = ModuleRef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_api_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModuleRef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModuleRef) ProtoMessage() {}

func (x *ModuleRef) ProtoReflect() protoreflect.Message {
	mi := &file_server_api_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ModuleRef.ProtoReflect.Descriptor instead.
func (*ModuleRef) Descriptor() ([]byte, []int) {
	return file_server_api_server_proto_rawDescGZIP(), []int{0}
}

func (x *ModuleRef) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ModuleRefs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Modules []*ModuleRef `protobuf:"bytes,1,rep,name=modules,proto3" json:"modules,omitempty"`
}

func (x *ModuleRefs) Reset() {
	*x = ModuleRefs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_api_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModuleRefs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModuleRefs) ProtoMessage() {}

func (x *ModuleRefs) ProtoReflect() protoreflect.Message {
	mi := &file_server_api_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ModuleRefs.ProtoReflect.Descriptor instead.
func (*ModuleRefs) Descriptor() ([]byte, []int) {
	return file_server_api_server_proto_rawDescGZIP(), []int{1}
}

func (x *ModuleRefs) GetModules() []*ModuleRef {
	if x != nil {
		return x.Modules
	}
	return nil
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State  State  `protobuf:"varint,1,opt,name=state,proto3,enum=gate.server.api.State" json:"state,omitempty"`
	Cause  Cause  `protobuf:"varint,2,opt,name=cause,proto3,enum=gate.server.api.Cause" json:"cause,omitempty"`
	Result int32  `protobuf:"varint,3,opt,name=result,proto3" json:"result,omitempty"`
	Error  string `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_api_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_server_api_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_server_api_server_proto_rawDescGZIP(), []int{2}
}

func (x *Status) GetState() State {
	if x != nil {
		return x.State
	}
	return State_NONEXISTENT
}

func (x *Status) GetCause() Cause {
	if x != nil {
		return x.Cause
	}
	return Cause_NORMAL
}

func (x *Status) GetResult() int32 {
	if x != nil {
		return x.Result
	}
	return 0
}

func (x *Status) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type InstanceStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Instance  string  `protobuf:"bytes,1,opt,name=instance,proto3" json:"instance,omitempty"`
	Status    *Status `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Transient bool    `protobuf:"varint,3,opt,name=transient,proto3" json:"transient,omitempty"`
	Debugging bool    `protobuf:"varint,4,opt,name=debugging,proto3" json:"debugging,omitempty"`
}

func (x *InstanceStatus) Reset() {
	*x = InstanceStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_api_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceStatus) ProtoMessage() {}

func (x *InstanceStatus) ProtoReflect() protoreflect.Message {
	mi := &file_server_api_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceStatus.ProtoReflect.Descriptor instead.
func (*InstanceStatus) Descriptor() ([]byte, []int) {
	return file_server_api_server_proto_rawDescGZIP(), []int{3}
}

func (x *InstanceStatus) GetInstance() string {
	if x != nil {
		return x.Instance
	}
	return ""
}

func (x *InstanceStatus) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *InstanceStatus) GetTransient() bool {
	if x != nil {
		return x.Transient
	}
	return false
}

func (x *InstanceStatus) GetDebugging() bool {
	if x != nil {
		return x.Debugging
	}
	return false
}

type Instances struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Instances []*InstanceStatus `protobuf:"bytes,1,rep,name=instances,proto3" json:"instances,omitempty"`
}

func (x *Instances) Reset() {
	*x = Instances{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_api_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Instances) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Instances) ProtoMessage() {}

func (x *Instances) ProtoReflect() protoreflect.Message {
	mi := &file_server_api_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Instances.ProtoReflect.Descriptor instead.
func (*Instances) Descriptor() ([]byte, []int) {
	return file_server_api_server_proto_rawDescGZIP(), []int{4}
}

func (x *Instances) GetInstances() []*InstanceStatus {
	if x != nil {
		return x.Instances
	}
	return nil
}

type DebugRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Op     DebugOp      `protobuf:"varint,1,opt,name=op,proto3,enum=gate.server.api.DebugOp" json:"op,omitempty"`
	Config *DebugConfig `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	Addr   uint64       `protobuf:"varint,3,opt,name=addr,proto3" json:"addr,omitempty"`
	Size   uint64       `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
}

func (x *DebugRequest) Reset() {
	*x = DebugRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_api_server_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DebugRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DebugRequest) ProtoMessage() {}

func (x *DebugRequest) ProtoReflect() protoreflect.Message {
	mi := &file_server_api_server_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DebugRequest.ProtoReflect.Descriptor instead.
func (*DebugRequest) Descriptor() ([]byte, []int) {
	return file_server_api_server_proto_rawDescGZIP(), []int{5}
}

func (x *DebugRequest) GetOp() DebugOp {
	if x != nil {
		return x.Op
	}
	return DebugOp_CONFIG_GET
}

func (x *DebugRequest) GetConfig() *DebugConfig {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *DebugRequest) GetAddr() uint64 {
	if x != nil {
		return x.Addr
	}
	return 0
}

func (x *DebugRequest) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

type DebugResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Module string       `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty"`
	Status *Status      `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Config *DebugConfig `protobuf:"bytes,3,opt,name=config,proto3" json:"config,omitempty"`
	Data   []byte       `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *DebugResponse) Reset() {
	*x = DebugResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_api_server_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DebugResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DebugResponse) ProtoMessage() {}

func (x *DebugResponse) ProtoReflect() protoreflect.Message {
	mi := &file_server_api_server_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DebugResponse.ProtoReflect.Descriptor instead.
func (*DebugResponse) Descriptor() ([]byte, []int) {
	return file_server_api_server_proto_rawDescGZIP(), []int{6}
}

func (x *DebugResponse) GetModule() string {
	if x != nil {
		return x.Module
	}
	return ""
}

func (x *DebugResponse) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *DebugResponse) GetConfig() *DebugConfig {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *DebugResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type DebugConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DebugInfo   bool     `protobuf:"varint,1,opt,name=debug_info,json=debugInfo,proto3" json:"debug_info,omitempty"` // TODO: specify semantics
	Breakpoints []uint64 `protobuf:"varint,2,rep,packed,name=breakpoints,proto3" json:"breakpoints,omitempty"`
}

func (x *DebugConfig) Reset() {
	*x = DebugConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_api_server_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DebugConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DebugConfig) ProtoMessage() {}

func (x *DebugConfig) ProtoReflect() protoreflect.Message {
	mi := &file_server_api_server_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DebugConfig.ProtoReflect.Descriptor instead.
func (*DebugConfig) Descriptor() ([]byte, []int) {
	return file_server_api_server_proto_rawDescGZIP(), []int{7}
}

func (x *DebugConfig) GetDebugInfo() bool {
	if x != nil {
		return x.DebugInfo
	}
	return false
}

func (x *DebugConfig) GetBreakpoints() []uint64 {
	if x != nil {
		return x.Breakpoints
	}
	return nil
}

var File_server_api_server_proto protoreflect.FileDescriptor

var file_server_api_server_proto_rawDesc = []byte{
	0x0a, 0x17, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x67, 0x61, 0x74, 0x65, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x22, 0x1b, 0x0a, 0x09, 0x4d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x66, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x42, 0x0a, 0x0a, 0x4d, 0x6f, 0x64, 0x75, 0x6c,
	0x65, 0x52, 0x65, 0x66, 0x73, 0x12, 0x34, 0x0a, 0x07, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x52,
	0x65, 0x66, 0x52, 0x07, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x22, 0x92, 0x01, 0x0a, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2c, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x63, 0x61, 0x75, 0x73, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x61, 0x75, 0x73, 0x65, 0x52, 0x05, 0x63, 0x61, 0x75,
	0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x22, 0x99, 0x01, 0x0a, 0x0e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12,
	0x2f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x1c, 0x0a, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x64, 0x65, 0x62, 0x75, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x09, 0x64, 0x65, 0x62, 0x75, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x22, 0x4a, 0x0a, 0x09,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x3d, 0x0a, 0x09, 0x69, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x67,
	0x61, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x09, 0x69,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x22, 0x96, 0x01, 0x0a, 0x0c, 0x44, 0x65, 0x62,
	0x75, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x02, 0x6f, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x62, 0x75, 0x67, 0x4f, 0x70, 0x52,
	0x02, 0x6f, 0x70, 0x12, 0x34, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x62, 0x75, 0x67, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x22, 0xa2, 0x01, 0x0a, 0x0d, 0x44, 0x65, 0x62, 0x75, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x2f, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x61,
	0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x34, 0x0a, 0x06,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x61, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44,
	0x65, 0x62, 0x75, 0x67, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x4e, 0x0a, 0x0b, 0x44, 0x65, 0x62, 0x75, 0x67, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x65, 0x62, 0x75, 0x67, 0x5f, 0x69,
	0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x64, 0x65, 0x62, 0x75, 0x67,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x20, 0x0a, 0x0b, 0x62, 0x72, 0x65, 0x61, 0x6b, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x72, 0x65, 0x61, 0x6b,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x2a, 0x5c, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x0f, 0x0a, 0x0b, 0x4e, 0x4f, 0x4e, 0x45, 0x58, 0x49, 0x53, 0x54, 0x45, 0x4e, 0x54, 0x10, 0x00,
	0x12, 0x0b, 0x0a, 0x07, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0d, 0x0a,
	0x09, 0x53, 0x55, 0x53, 0x50, 0x45, 0x4e, 0x44, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06,
	0x48, 0x41, 0x4c, 0x54, 0x45, 0x44, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x45, 0x52, 0x4d,
	0x49, 0x4e, 0x41, 0x54, 0x45, 0x44, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x4b, 0x49, 0x4c, 0x4c,
	0x45, 0x44, 0x10, 0x05, 0x2a, 0xa3, 0x02, 0x0a, 0x05, 0x43, 0x61, 0x75, 0x73, 0x65, 0x12, 0x0a,
	0x0a, 0x06, 0x4e, 0x4f, 0x52, 0x4d, 0x41, 0x4c, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e,
	0x52, 0x45, 0x41, 0x43, 0x48, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x03, 0x12, 0x18, 0x0a, 0x14, 0x43,
	0x41, 0x4c, 0x4c, 0x5f, 0x53, 0x54, 0x41, 0x43, 0x4b, 0x5f, 0x45, 0x58, 0x48, 0x41, 0x55, 0x53,
	0x54, 0x45, 0x44, 0x10, 0x04, 0x12, 0x1f, 0x0a, 0x1b, 0x4d, 0x45, 0x4d, 0x4f, 0x52, 0x59, 0x5f,
	0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4f, 0x55, 0x54, 0x5f, 0x4f, 0x46, 0x5f, 0x42, 0x4f,
	0x55, 0x4e, 0x44, 0x53, 0x10, 0x05, 0x12, 0x25, 0x0a, 0x21, 0x49, 0x4e, 0x44, 0x49, 0x52, 0x45,
	0x43, 0x54, 0x5f, 0x43, 0x41, 0x4c, 0x4c, 0x5f, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x5f, 0x4f, 0x55,
	0x54, 0x5f, 0x4f, 0x46, 0x5f, 0x42, 0x4f, 0x55, 0x4e, 0x44, 0x53, 0x10, 0x06, 0x12, 0x24, 0x0a,
	0x20, 0x49, 0x4e, 0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x5f, 0x43, 0x41, 0x4c, 0x4c, 0x5f, 0x53,
	0x49, 0x47, 0x4e, 0x41, 0x54, 0x55, 0x52, 0x45, 0x5f, 0x4d, 0x49, 0x53, 0x4d, 0x41, 0x54, 0x43,
	0x48, 0x10, 0x07, 0x12, 0x1a, 0x0a, 0x16, 0x49, 0x4e, 0x54, 0x45, 0x47, 0x45, 0x52, 0x5f, 0x44,
	0x49, 0x56, 0x49, 0x44, 0x45, 0x5f, 0x42, 0x59, 0x5f, 0x5a, 0x45, 0x52, 0x4f, 0x10, 0x08, 0x12,
	0x14, 0x0a, 0x10, 0x49, 0x4e, 0x54, 0x45, 0x47, 0x45, 0x52, 0x5f, 0x4f, 0x56, 0x45, 0x52, 0x46,
	0x4c, 0x4f, 0x57, 0x10, 0x09, 0x12, 0x0e, 0x0a, 0x0a, 0x42, 0x52, 0x45, 0x41, 0x4b, 0x50, 0x4f,
	0x49, 0x4e, 0x54, 0x10, 0x0a, 0x12, 0x12, 0x0a, 0x0e, 0x41, 0x42, 0x49, 0x5f, 0x44, 0x45, 0x46,
	0x49, 0x43, 0x49, 0x45, 0x4e, 0x43, 0x59, 0x10, 0x1b, 0x12, 0x11, 0x0a, 0x0d, 0x41, 0x42, 0x49,
	0x5f, 0x56, 0x49, 0x4f, 0x4c, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x1c, 0x12, 0x0c, 0x0a, 0x08,
	0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x4c, 0x10, 0x1d, 0x2a, 0x85, 0x01, 0x0a, 0x07, 0x44,
	0x65, 0x62, 0x75, 0x67, 0x4f, 0x70, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47,
	0x5f, 0x47, 0x45, 0x54, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47,
	0x5f, 0x53, 0x45, 0x54, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47,
	0x5f, 0x55, 0x4e, 0x49, 0x4f, 0x4e, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x4f, 0x4e, 0x46,
	0x49, 0x47, 0x5f, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x03, 0x12,
	0x10, 0x0a, 0x0c, 0x52, 0x45, 0x41, 0x44, 0x5f, 0x47, 0x4c, 0x4f, 0x42, 0x41, 0x4c, 0x53, 0x10,
	0x04, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x45, 0x41, 0x44, 0x5f, 0x4d, 0x45, 0x4d, 0x4f, 0x52, 0x59,
	0x10, 0x05, 0x12, 0x0e, 0x0a, 0x0a, 0x52, 0x45, 0x41, 0x44, 0x5f, 0x53, 0x54, 0x41, 0x43, 0x4b,
	0x10, 0x06, 0x42, 0x1f, 0x5a, 0x1d, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x75,
	0x74, 0x65, 0x72, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f,
	0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_server_api_server_proto_rawDescOnce sync.Once
	file_server_api_server_proto_rawDescData = file_server_api_server_proto_rawDesc
)

func file_server_api_server_proto_rawDescGZIP() []byte {
	file_server_api_server_proto_rawDescOnce.Do(func() {
		file_server_api_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_api_server_proto_rawDescData)
	})
	return file_server_api_server_proto_rawDescData
}

var file_server_api_server_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_server_api_server_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_server_api_server_proto_goTypes = []interface{}{
	(State)(0),             // 0: gate.server.api.State
	(Cause)(0),             // 1: gate.server.api.Cause
	(DebugOp)(0),           // 2: gate.server.api.DebugOp
	(*ModuleRef)(nil),      // 3: gate.server.api.ModuleRef
	(*ModuleRefs)(nil),     // 4: gate.server.api.ModuleRefs
	(*Status)(nil),         // 5: gate.server.api.Status
	(*InstanceStatus)(nil), // 6: gate.server.api.InstanceStatus
	(*Instances)(nil),      // 7: gate.server.api.Instances
	(*DebugRequest)(nil),   // 8: gate.server.api.DebugRequest
	(*DebugResponse)(nil),  // 9: gate.server.api.DebugResponse
	(*DebugConfig)(nil),    // 10: gate.server.api.DebugConfig
}
var file_server_api_server_proto_depIdxs = []int32{
	3,  // 0: gate.server.api.ModuleRefs.modules:type_name -> gate.server.api.ModuleRef
	0,  // 1: gate.server.api.Status.state:type_name -> gate.server.api.State
	1,  // 2: gate.server.api.Status.cause:type_name -> gate.server.api.Cause
	5,  // 3: gate.server.api.InstanceStatus.status:type_name -> gate.server.api.Status
	6,  // 4: gate.server.api.Instances.instances:type_name -> gate.server.api.InstanceStatus
	2,  // 5: gate.server.api.DebugRequest.op:type_name -> gate.server.api.DebugOp
	10, // 6: gate.server.api.DebugRequest.config:type_name -> gate.server.api.DebugConfig
	5,  // 7: gate.server.api.DebugResponse.status:type_name -> gate.server.api.Status
	10, // 8: gate.server.api.DebugResponse.config:type_name -> gate.server.api.DebugConfig
	9,  // [9:9] is the sub-list for method output_type
	9,  // [9:9] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_server_api_server_proto_init() }
func file_server_api_server_proto_init() {
	if File_server_api_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_api_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ModuleRef); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_api_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ModuleRefs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_api_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_api_server_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_api_server_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Instances); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_api_server_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DebugRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_api_server_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DebugResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_server_api_server_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DebugConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_server_api_server_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_server_api_server_proto_goTypes,
		DependencyIndexes: file_server_api_server_proto_depIdxs,
		EnumInfos:         file_server_api_server_proto_enumTypes,
		MessageInfos:      file_server_api_server_proto_msgTypes,
	}.Build()
	File_server_api_server_proto = out.File
	file_server_api_server_proto_rawDesc = nil
	file_server_api_server_proto_goTypes = nil
	file_server_api_server_proto_depIdxs = nil
}
