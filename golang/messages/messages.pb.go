// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: proto/messages.proto

// 包定义

package messages

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

// 消息类型
type ProcessMessage_Type int32

const (
	ProcessMessage_START                ProcessMessage_Type = 0
	ProcessMessage_STOP                 ProcessMessage_Type = 1
	ProcessMessage_ROCKETMQ_MESSAGE     ProcessMessage_Type = 2
	ProcessMessage_ROCKETMQ_MESSAGE_ACK ProcessMessage_Type = 3
	ProcessMessage_RESULT               ProcessMessage_Type = 4
)

// Enum value maps for ProcessMessage_Type.
var (
	ProcessMessage_Type_name = map[int32]string{
		0: "START",
		1: "STOP",
		2: "ROCKETMQ_MESSAGE",
		3: "ROCKETMQ_MESSAGE_ACK",
		4: "RESULT",
	}
	ProcessMessage_Type_value = map[string]int32{
		"START":                0,
		"STOP":                 1,
		"ROCKETMQ_MESSAGE":     2,
		"ROCKETMQ_MESSAGE_ACK": 3,
		"RESULT":               4,
	}
)

func (x ProcessMessage_Type) Enum() *ProcessMessage_Type {
	p := new(ProcessMessage_Type)
	*p = x
	return p
}

func (x ProcessMessage_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProcessMessage_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_messages_proto_enumTypes[0].Descriptor()
}

func (ProcessMessage_Type) Type() protoreflect.EnumType {
	return &file_proto_messages_proto_enumTypes[0]
}

func (x ProcessMessage_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProcessMessage_Type.Descriptor instead.
func (ProcessMessage_Type) EnumDescriptor() ([]byte, []int) {
	return file_proto_messages_proto_rawDescGZIP(), []int{0, 0}
}

// 消息定义
type ProcessMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 请求 ID
	RequestId string `protobuf:"bytes,1,opt,name=requestId,proto3" json:"requestId,omitempty"`
	// 消息类型
	Type ProcessMessage_Type `protobuf:"varint,2,opt,name=type,proto3,enum=messages.ProcessMessage_Type" json:"type,omitempty"`
	// 信息
	Info *ProcessMessage_Info `protobuf:"bytes,3,opt,name=info,proto3,oneof" json:"info,omitempty"`
}

func (x *ProcessMessage) Reset() {
	*x = ProcessMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessMessage) ProtoMessage() {}

func (x *ProcessMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessMessage.ProtoReflect.Descriptor instead.
func (*ProcessMessage) Descriptor() ([]byte, []int) {
	return file_proto_messages_proto_rawDescGZIP(), []int{0}
}

func (x *ProcessMessage) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *ProcessMessage) GetType() ProcessMessage_Type {
	if x != nil {
		return x.Type
	}
	return ProcessMessage_START
}

func (x *ProcessMessage) GetInfo() *ProcessMessage_Info {
	if x != nil {
		return x.Info
	}
	return nil
}

type ProcessMessage_Info struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 状态码
	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// 消息 ID
	MessageId string `protobuf:"bytes,2,opt,name=messageId,proto3" json:"messageId,omitempty"`
	// 消息内容
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ProcessMessage_Info) Reset() {
	*x = ProcessMessage_Info{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessMessage_Info) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessMessage_Info) ProtoMessage() {}

func (x *ProcessMessage_Info) ProtoReflect() protoreflect.Message {
	mi := &file_proto_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessMessage_Info.ProtoReflect.Descriptor instead.
func (*ProcessMessage_Info) Descriptor() ([]byte, []int) {
	return file_proto_messages_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ProcessMessage_Info) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ProcessMessage_Info) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *ProcessMessage_Info) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_messages_proto protoreflect.FileDescriptor

var file_proto_messages_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x22, 0xcf, 0x02, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49,
	0x64, 0x12, 0x31, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1d, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x36, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x50, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x49, 0x6e, 0x66,
	0x6f, 0x48, 0x00, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x88, 0x01, 0x01, 0x1a, 0x52, 0x0a, 0x04,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0x57, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x54, 0x41, 0x52,
	0x54, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x54, 0x4f, 0x50, 0x10, 0x01, 0x12, 0x14, 0x0a,
	0x10, 0x52, 0x4f, 0x43, 0x4b, 0x45, 0x54, 0x4d, 0x51, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47,
	0x45, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x52, 0x4f, 0x43, 0x4b, 0x45, 0x54, 0x4d, 0x51, 0x5f,
	0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x41, 0x43, 0x4b, 0x10, 0x03, 0x12, 0x0a, 0x0a,
	0x06, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x10, 0x04, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x69, 0x6e,
	0x66, 0x6f, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_messages_proto_rawDescOnce sync.Once
	file_proto_messages_proto_rawDescData = file_proto_messages_proto_rawDesc
)

func file_proto_messages_proto_rawDescGZIP() []byte {
	file_proto_messages_proto_rawDescOnce.Do(func() {
		file_proto_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_messages_proto_rawDescData)
	})
	return file_proto_messages_proto_rawDescData
}

var file_proto_messages_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_messages_proto_goTypes = []any{
	(ProcessMessage_Type)(0),    // 0: messages.ProcessMessage.Type
	(*ProcessMessage)(nil),      // 1: messages.ProcessMessage
	(*ProcessMessage_Info)(nil), // 2: messages.ProcessMessage.Info
}
var file_proto_messages_proto_depIdxs = []int32{
	0, // 0: messages.ProcessMessage.type:type_name -> messages.ProcessMessage.Type
	2, // 1: messages.ProcessMessage.info:type_name -> messages.ProcessMessage.Info
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_messages_proto_init() }
func file_proto_messages_proto_init() {
	if File_proto_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_messages_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ProcessMessage); i {
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
		file_proto_messages_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ProcessMessage_Info); i {
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
	file_proto_messages_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_messages_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_messages_proto_goTypes,
		DependencyIndexes: file_proto_messages_proto_depIdxs,
		EnumInfos:         file_proto_messages_proto_enumTypes,
		MessageInfos:      file_proto_messages_proto_msgTypes,
	}.Build()
	File_proto_messages_proto = out.File
	file_proto_messages_proto_rawDesc = nil
	file_proto_messages_proto_goTypes = nil
	file_proto_messages_proto_depIdxs = nil
}
