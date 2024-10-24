// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: rpc_conn.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DeliverMessageReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReceiverId int64  `protobuf:"varint,1,opt,name=receiver_id,json=receiverId,proto3" json:"receiver_id,omitempty"` //  消息接收者
	Data       []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`                                // 要投递的消息
}

func (x *DeliverMessageReq) Reset() {
	*x = DeliverMessageReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_conn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliverMessageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliverMessageReq) ProtoMessage() {}

func (x *DeliverMessageReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_conn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliverMessageReq.ProtoReflect.Descriptor instead.
func (*DeliverMessageReq) Descriptor() ([]byte, []int) {
	return file_rpc_conn_proto_rawDescGZIP(), []int{0}
}

func (x *DeliverMessageReq) GetReceiverId() int64 {
	if x != nil {
		return x.ReceiverId
	}
	return 0
}

func (x *DeliverMessageReq) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type DeliverMessageGroupReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReceiverId_2Data map[int64][]byte `protobuf:"bytes,1,rep,name=receiver_id_2_data,json=receiverId2Data,proto3" json:"receiver_id_2_data,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // 消息接受者到要投递的消息的映射
}

func (x *DeliverMessageGroupReq) Reset() {
	*x = DeliverMessageGroupReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_conn_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliverMessageGroupReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliverMessageGroupReq) ProtoMessage() {}

func (x *DeliverMessageGroupReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_conn_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliverMessageGroupReq.ProtoReflect.Descriptor instead.
func (*DeliverMessageGroupReq) Descriptor() ([]byte, []int) {
	return file_rpc_conn_proto_rawDescGZIP(), []int{1}
}

func (x *DeliverMessageGroupReq) GetReceiverId_2Data() map[int64][]byte {
	if x != nil {
		return x.ReceiverId_2Data
	}
	return nil
}

var File_rpc_conn_proto protoreflect.FileDescriptor

var file_rpc_conn_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x70, 0x62, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x48, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x72, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xba, 0x01, 0x0a, 0x16,
	0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x12, 0x5c, 0x0a, 0x12, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x5f, 0x32, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x2e, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x49, 0x64, 0x32, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x0f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x49, 0x64, 0x32,
	0x44, 0x61, 0x74, 0x61, 0x1a, 0x42, 0x0a, 0x14, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72,
	0x49, 0x64, 0x32, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x95, 0x01, 0x0a, 0x07, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x12, 0x3f, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x15, 0x2e, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x69,
	0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x49, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x2e, 0x70,
	0x62, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_rpc_conn_proto_rawDescOnce sync.Once
	file_rpc_conn_proto_rawDescData = file_rpc_conn_proto_rawDesc
)

func file_rpc_conn_proto_rawDescGZIP() []byte {
	file_rpc_conn_proto_rawDescOnce.Do(func() {
		file_rpc_conn_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_conn_proto_rawDescData)
	})
	return file_rpc_conn_proto_rawDescData
}

var file_rpc_conn_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_rpc_conn_proto_goTypes = []any{
	(*DeliverMessageReq)(nil),      // 0: pb.DeliverMessageReq
	(*DeliverMessageGroupReq)(nil), // 1: pb.DeliverMessageGroupReq
	nil,                            // 2: pb.DeliverMessageGroupReq.ReceiverId2DataEntry
	(*emptypb.Empty)(nil),          // 3: google.protobuf.Empty
}
var file_rpc_conn_proto_depIdxs = []int32{
	2, // 0: pb.DeliverMessageGroupReq.receiver_id_2_data:type_name -> pb.DeliverMessageGroupReq.ReceiverId2DataEntry
	0, // 1: pb.Connect.DeliverMessage:input_type -> pb.DeliverMessageReq
	1, // 2: pb.Connect.DeliverGroupMessage:input_type -> pb.DeliverMessageGroupReq
	3, // 3: pb.Connect.DeliverMessage:output_type -> google.protobuf.Empty
	3, // 4: pb.Connect.DeliverGroupMessage:output_type -> google.protobuf.Empty
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_conn_proto_init() }
func file_rpc_conn_proto_init() {
	if File_rpc_conn_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_conn_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*DeliverMessageReq); i {
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
		file_rpc_conn_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*DeliverMessageGroupReq); i {
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
			RawDescriptor: file_rpc_conn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_conn_proto_goTypes,
		DependencyIndexes: file_rpc_conn_proto_depIdxs,
		MessageInfos:      file_rpc_conn_proto_msgTypes,
	}.Build()
	File_rpc_conn_proto = out.File
	file_rpc_conn_proto_rawDesc = nil
	file_rpc_conn_proto_goTypes = nil
	file_rpc_conn_proto_depIdxs = nil
}
