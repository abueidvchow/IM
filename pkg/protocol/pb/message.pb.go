// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: message.proto

package pb

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

// 会话类型
type SessionType int32

const (
	SessionType_ST_UNKNOWN SessionType = 0 // 未知
	SessionType_ST_SINGLE  SessionType = 1 // 私聊
	SessionType_ST_GROUP   SessionType = 2 // 群聊
)

// Enum value maps for SessionType.
var (
	SessionType_name = map[int32]string{
		0: "ST_UNKNOWN",
		1: "ST_SINGLE",
		2: "ST_GROUP",
	}
	SessionType_value = map[string]int32{
		"ST_UNKNOWN": 0,
		"ST_SINGLE":  1,
		"ST_GROUP":   2,
	}
)

func (x SessionType) Enum() *SessionType {
	p := new(SessionType)
	*p = x
	return p
}

func (x SessionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SessionType) Descriptor() protoreflect.EnumDescriptor {
	return file_message_proto_enumTypes[0].Descriptor()
}

func (SessionType) Type() protoreflect.EnumType {
	return &file_message_proto_enumTypes[0]
}

func (x SessionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SessionType.Descriptor instead.
func (SessionType) EnumDescriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

// 消息类型
type MessageType int32

const (
	MessageType_MT_UNKNOWN MessageType = 0 // 未知
	MessageType_MT_TEXT    MessageType = 1 // 文字类型
	MessageType_MT_PICTURE MessageType = 2 // 图片类型
	MessageType_MT_VOICE   MessageType = 3 // 语音类型
)

// Enum value maps for MessageType.
var (
	MessageType_name = map[int32]string{
		0: "MT_UNKNOWN",
		1: "MT_TEXT",
		2: "MT_PICTURE",
		3: "MT_VOICE",
	}
	MessageType_value = map[string]int32{
		"MT_UNKNOWN": 0,
		"MT_TEXT":    1,
		"MT_PICTURE": 2,
		"MT_VOICE":   3,
	}
)

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}

func (x MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_message_proto_enumTypes[1].Descriptor()
}

func (MessageType) Type() protoreflect.EnumType {
	return &file_message_proto_enumTypes[1]
}

func (x MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageType.Descriptor instead.
func (MessageType) EnumDescriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{1}
}

// websocket消息类型
type CmdType int32

const (
	CmdType_CT_UNKNOWN CmdType = 0 // 未知
	CmdType_CT_MESSAGE CmdType = 1 // 消息投递
	CmdType_CT_SYNC    CmdType = 2 // 离线消息同步
)

// Enum value maps for CmdType.
var (
	CmdType_name = map[int32]string{
		0: "CT_UNKNOWN",
		1: "CT_MESSAGE",
		2: "CT_SYNC",
	}
	CmdType_value = map[string]int32{
		"CT_UNKNOWN": 0,
		"CT_MESSAGE": 1,
		"CT_SYNC":    2,
	}
)

func (x CmdType) Enum() *CmdType {
	p := new(CmdType)
	*p = x
	return p
}

func (x CmdType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CmdType) Descriptor() protoreflect.EnumDescriptor {
	return file_message_proto_enumTypes[2].Descriptor()
}

func (CmdType) Type() protoreflect.EnumType {
	return &file_message_proto_enumTypes[2]
}

func (x CmdType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CmdType.Descriptor instead.
func (CmdType) EnumDescriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{2}
}

type Input struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type CmdType `protobuf:"varint,1,opt,name=type,proto3,enum=pb.CmdType" json:"type,omitempty"`
	Data []byte  `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Input) Reset() {
	*x = Input{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Input) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Input) ProtoMessage() {}

func (x *Input) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Input.ProtoReflect.Descriptor instead.
func (*Input) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

func (x *Input) GetType() CmdType {
	if x != nil {
		return x.Type
	}
	return CmdType_CT_UNKNOWN
}

func (x *Input) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionType SessionType `protobuf:"varint,1,opt,name=session_type,json=sessionType,proto3,enum=pb.SessionType" json:"session_type,omitempty"` // 会话类型 单聊、群聊
	ReceiverId  int64       `protobuf:"varint,2,opt,name=receiver_id,json=receiverId,proto3" json:"receiver_id,omitempty"`                        // 接收者id 用户id/群组id
	SenderId    int64       `protobuf:"varint,3,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`                              // 发送者id
	MessageType MessageType `protobuf:"varint,4,opt,name=message_type,json=messageType,proto3,enum=pb.MessageType" json:"message_type,omitempty"` // 消息类型 文本、图片、语音
	Content     []byte      `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`                                                 // 实际用户所发数据
	Seq         int64       `protobuf:"varint,6,opt,name=seq,proto3" json:"seq,omitempty"`                                                        // 客户端的最大消息同步序号
	SendTime    int64       `protobuf:"varint,7,opt,name=send_time,json=sendTime,proto3" json:"send_time,omitempty"`                              // 消息发送时间戳，ms
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[1]
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
	return file_message_proto_rawDescGZIP(), []int{1}
}

func (x *Message) GetSessionType() SessionType {
	if x != nil {
		return x.SessionType
	}
	return SessionType_ST_UNKNOWN
}

func (x *Message) GetReceiverId() int64 {
	if x != nil {
		return x.ReceiverId
	}
	return 0
}

func (x *Message) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *Message) GetMessageType() MessageType {
	if x != nil {
		return x.MessageType
	}
	return MessageType_MT_UNKNOWN
}

func (x *Message) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *Message) GetSeq() int64 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *Message) GetSendTime() int64 {
	if x != nil {
		return x.SendTime
	}
	return 0
}

var File_message_proto protoreflect.FileDescriptor

var file_message_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x22, 0x3c, 0x0a, 0x05, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x1f, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x70, 0x62, 0x2e,
	0x43, 0x6d, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x22, 0xf8, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x32, 0x0a,
	0x0c, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x32, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x73, 0x65, 0x71, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x73, 0x65, 0x71, 0x12,
	0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x2a, 0x3a, 0x0a, 0x0b,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x53,
	0x54, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x53,
	0x54, 0x5f, 0x53, 0x49, 0x4e, 0x47, 0x4c, 0x45, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x54,
	0x5f, 0x47, 0x52, 0x4f, 0x55, 0x50, 0x10, 0x02, 0x2a, 0x48, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x4d, 0x54, 0x5f, 0x55, 0x4e,
	0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x4d, 0x54, 0x5f, 0x54, 0x45,
	0x58, 0x54, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x4d, 0x54, 0x5f, 0x50, 0x49, 0x43, 0x54, 0x55,
	0x52, 0x45, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x4d, 0x54, 0x5f, 0x56, 0x4f, 0x49, 0x43, 0x45,
	0x10, 0x03, 0x2a, 0x36, 0x0a, 0x07, 0x43, 0x6d, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a,
	0x0a, 0x43, 0x54, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0e, 0x0a,
	0x0a, 0x43, 0x54, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a,
	0x07, 0x43, 0x54, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x10, 0x02, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2e,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_proto_rawDescOnce sync.Once
	file_message_proto_rawDescData = file_message_proto_rawDesc
)

func file_message_proto_rawDescGZIP() []byte {
	file_message_proto_rawDescOnce.Do(func() {
		file_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_proto_rawDescData)
	})
	return file_message_proto_rawDescData
}

var file_message_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_message_proto_goTypes = []any{
	(SessionType)(0), // 0: pb.SessionType
	(MessageType)(0), // 1: pb.MessageType
	(CmdType)(0),     // 2: pb.CmdType
	(*Input)(nil),    // 3: pb.Input
	(*Message)(nil),  // 4: pb.Message
}
var file_message_proto_depIdxs = []int32{
	2, // 0: pb.Input.type:type_name -> pb.CmdType
	0, // 1: pb.Message.session_type:type_name -> pb.SessionType
	1, // 2: pb.Message.message_type:type_name -> pb.MessageType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_message_proto_init() }
func file_message_proto_init() {
	if File_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Input); i {
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
		file_message_proto_msgTypes[1].Exporter = func(v any, i int) any {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_proto_goTypes,
		DependencyIndexes: file_message_proto_depIdxs,
		EnumInfos:         file_message_proto_enumTypes,
		MessageInfos:      file_message_proto_msgTypes,
	}.Build()
	File_message_proto = out.File
	file_message_proto_rawDesc = nil
	file_message_proto_goTypes = nil
	file_message_proto_depIdxs = nil
}