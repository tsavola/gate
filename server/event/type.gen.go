// Code generated by internal/cmd/event-types.  DO NOT EDIT.

package event

func (x *FailInternal) EventName() string         { return "FAIL_INTERNAL" }
func (x *FailNetwork) EventName() string          { return "FAIL_NETWORK" }
func (x *FailProtocol) EventName() string         { return "FAIL_PROTOCOL" }
func (x *FailRequest) EventName() string          { return "FAIL_REQUEST" }
func (x *IfaceAccess) EventName() string          { return "IFACE_ACCESS" }
func (x *InstanceConnect) EventName() string      { return "INSTANCE_CONNECT" }
func (x *InstanceCreateLocal) EventName() string  { return "INSTANCE_CREATE_LOCAL" }
func (x *InstanceCreateStream) EventName() string { return "INSTANCE_CREATE_STREAM" }
func (x *InstanceDelete) EventName() string       { return "INSTANCE_DELETE" }
func (x *InstanceDisconnect) EventName() string   { return "INSTANCE_DISCONNECT" }
func (x *InstanceKill) EventName() string         { return "INSTANCE_KILL" }
func (x *InstanceList) EventName() string         { return "INSTANCE_LIST" }
func (x *InstanceResume) EventName() string       { return "INSTANCE_RESUME" }
func (x *InstanceSnapshot) EventName() string     { return "INSTANCE_SNAPSHOT" }
func (x *InstanceStatus) EventName() string       { return "INSTANCE_STATUS" }
func (x *InstanceSuspend) EventName() string      { return "INSTANCE_SUSPEND" }
func (x *InstanceWait) EventName() string         { return "INSTANCE_WAIT" }
func (x *ModuleDownload) EventName() string       { return "MODULE_DOWNLOAD" }
func (x *ModuleList) EventName() string           { return "MODULE_LIST" }
func (x *ModuleSourceExist) EventName() string    { return "MODULE_SOURCE_EXIST" }
func (x *ModuleSourceNew) EventName() string      { return "MODULE_SOURCE_NEW" }
func (x *ModuleUnref) EventName() string          { return "MODULE_UNREF" }
func (x *ModuleUploadExist) EventName() string    { return "MODULE_UPLOAD_EXIST" }
func (x *ModuleUploadNew) EventName() string      { return "MODULE_UPLOAD_NEW" }

func (*FailInternal) EventType() int32         { return int32(Type_FAIL_INTERNAL) }
func (*FailNetwork) EventType() int32          { return int32(Type_FAIL_NETWORK) }
func (*FailProtocol) EventType() int32         { return int32(Type_FAIL_PROTOCOL) }
func (*FailRequest) EventType() int32          { return int32(Type_FAIL_REQUEST) }
func (*IfaceAccess) EventType() int32          { return int32(Type_IFACE_ACCESS) }
func (*InstanceConnect) EventType() int32      { return int32(Type_INSTANCE_CONNECT) }
func (*InstanceCreateLocal) EventType() int32  { return int32(Type_INSTANCE_CREATE_LOCAL) }
func (*InstanceCreateStream) EventType() int32 { return int32(Type_INSTANCE_CREATE_STREAM) }
func (*InstanceDelete) EventType() int32       { return int32(Type_INSTANCE_DELETE) }
func (*InstanceDisconnect) EventType() int32   { return int32(Type_INSTANCE_DISCONNECT) }
func (*InstanceKill) EventType() int32         { return int32(Type_INSTANCE_KILL) }
func (*InstanceList) EventType() int32         { return int32(Type_INSTANCE_LIST) }
func (*InstanceResume) EventType() int32       { return int32(Type_INSTANCE_RESUME) }
func (*InstanceSnapshot) EventType() int32     { return int32(Type_INSTANCE_SNAPSHOT) }
func (*InstanceStatus) EventType() int32       { return int32(Type_INSTANCE_STATUS) }
func (*InstanceSuspend) EventType() int32      { return int32(Type_INSTANCE_SUSPEND) }
func (*InstanceWait) EventType() int32         { return int32(Type_INSTANCE_WAIT) }
func (*ModuleDownload) EventType() int32       { return int32(Type_MODULE_DOWNLOAD) }
func (*ModuleList) EventType() int32           { return int32(Type_MODULE_LIST) }
func (*ModuleSourceExist) EventType() int32    { return int32(Type_MODULE_SOURCE_EXIST) }
func (*ModuleSourceNew) EventType() int32      { return int32(Type_MODULE_SOURCE_NEW) }
func (*ModuleUnref) EventType() int32          { return int32(Type_MODULE_UNREF) }
func (*ModuleUploadExist) EventType() int32    { return int32(Type_MODULE_UPLOAD_EXIST) }
func (*ModuleUploadNew) EventType() int32      { return int32(Type_MODULE_UPLOAD_NEW) }
