// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.0
// source: cmd.proto

package internal

import (
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

type CrawlState int32

const (
	CrawlState_Uninitialized CrawlState = 0
	CrawlState_NotLogged     CrawlState = 1
	CrawlState_Ready         CrawlState = 2
)

// Enum value maps for CrawlState.
var (
	CrawlState_name = map[int32]string{
		0: "Uninitialized",
		1: "NotLogged",
		2: "Ready",
	}
	CrawlState_value = map[string]int32{
		"Uninitialized": 0,
		"NotLogged":     1,
		"Ready":         2,
	}
)

func (x CrawlState) Enum() *CrawlState {
	p := new(CrawlState)
	*p = x
	return p
}

func (x CrawlState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CrawlState) Descriptor() protoreflect.EnumDescriptor {
	return file_cmd_proto_enumTypes[0].Descriptor()
}

func (CrawlState) Type() protoreflect.EnumType {
	return &file_cmd_proto_enumTypes[0]
}

func (x CrawlState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CrawlState.Descriptor instead.
func (CrawlState) EnumDescriptor() ([]byte, []int) {
	return file_cmd_proto_rawDescGZIP(), []int{0}
}

type CrawlCmd int32

const (
	CrawlCmd_Initial  CrawlCmd = 0
	CrawlCmd_Login    CrawlCmd = 1
	CrawlCmd_LoginAck CrawlCmd = 2
	CrawlCmd_StateAck CrawlCmd = 3
	CrawlCmd_Crawl    CrawlCmd = 4
)

// Enum value maps for CrawlCmd.
var (
	CrawlCmd_name = map[int32]string{
		0: "Initial",
		1: "Login",
		2: "LoginAck",
		3: "StateAck",
		4: "Account",
	}
	CrawlCmd_value = map[string]int32{
		"Initial":  0,
		"Login":    1,
		"LoginAck": 2,
		"StateAck": 3,
		"Account":  4,
	}
)

func (x CrawlCmd) Enum() *CrawlCmd {
	p := new(CrawlCmd)
	*p = x
	return p
}

func (x CrawlCmd) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CrawlCmd) Descriptor() protoreflect.EnumDescriptor {
	return file_cmd_proto_enumTypes[1].Descriptor()
}

func (CrawlCmd) Type() protoreflect.EnumType {
	return &file_cmd_proto_enumTypes[1]
}

func (x CrawlCmd) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CrawlCmd.Descriptor instead.
func (CrawlCmd) EnumDescriptor() ([]byte, []int) {
	return file_cmd_proto_rawDescGZIP(), []int{1}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   CrawlCmd    `protobuf:"varint,1,opt,name=cmd,proto3,enum=CrawlCmd" json:"cmd,omitempty"`
	State *CrawlState `protobuf:"varint,2,opt,name=state,proto3,enum=CrawlState,oneof" json:"state,omitempty"`
	Data  []byte      `protobuf:"bytes,3,opt,name=data,proto3,oneof" json:"data,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cmd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_cmd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_cmd_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetCmd() CrawlCmd {
	if x != nil {
		return x.Cmd
	}
	return CrawlCmd_Initial
}

func (x *Message) GetState() CrawlState {
	if x != nil && x.State != nil {
		return *x.State
	}
	return CrawlState_Uninitialized
}

func (x *Message) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_cmd_proto protoreflect.FileDescriptor

var file_cmd_proto_rawDesc = []byte{
	0x0a, 0x09, 0x63, 0x6d, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7a, 0x0a, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x09, 0x2e, 0x43, 0x72, 0x61, 0x77, 0x6c, 0x43, 0x6d, 0x64, 0x52, 0x03,
	0x63, 0x6d, 0x64, 0x12, 0x26, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x43, 0x72, 0x61, 0x77, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48,
	0x00, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x01, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x42, 0x07,
	0x0a, 0x05, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x39, 0x0a, 0x0a, 0x43, 0x72, 0x61, 0x77, 0x6c,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x6e, 0x69, 0x6e, 0x69, 0x74, 0x69,
	0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x4e, 0x6f, 0x74, 0x4c,
	0x6f, 0x67, 0x67, 0x65, 0x64, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x52, 0x65, 0x61, 0x64, 0x79,
	0x10, 0x02, 0x2a, 0x49, 0x0a, 0x08, 0x43, 0x72, 0x61, 0x77, 0x6c, 0x43, 0x6d, 0x64, 0x12, 0x0b,
	0x0a, 0x07, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x41,
	0x63, 0x6b, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x74, 0x61, 0x74, 0x65, 0x41, 0x63, 0x6b,
	0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x72, 0x61, 0x77, 0x6c, 0x10, 0x04, 0x42, 0x0c, 0x5a,
	0x0a, 0x2e, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_cmd_proto_rawDescOnce sync.Once
	file_cmd_proto_rawDescData = file_cmd_proto_rawDesc
)

func file_cmd_proto_rawDescGZIP() []byte {
	file_cmd_proto_rawDescOnce.Do(func() {
		file_cmd_proto_rawDescData = protoimpl.X.CompressGZIP(file_cmd_proto_rawDescData)
	})
	return file_cmd_proto_rawDescData
}

var file_cmd_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_cmd_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_cmd_proto_goTypes = []any{
	(CrawlState)(0), // 0: CrawlState
	(CrawlCmd)(0),   // 1: CrawlCmd
	(*Message)(nil), // 2: Message
}
var file_cmd_proto_depIdxs = []int32{
	1, // 0: Message.cmd:type_name -> CrawlCmd
	0, // 1: Message.state:type_name -> CrawlState
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_cmd_proto_init() }
func file_cmd_proto_init() {
	if File_cmd_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cmd_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Message); i {
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
	file_cmd_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cmd_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cmd_proto_goTypes,
		DependencyIndexes: file_cmd_proto_depIdxs,
		EnumInfos:         file_cmd_proto_enumTypes,
		MessageInfos:      file_cmd_proto_msgTypes,
	}.Build()
	File_cmd_proto = out.File
	file_cmd_proto_rawDesc = nil
	file_cmd_proto_goTypes = nil
	file_cmd_proto_depIdxs = nil
}
