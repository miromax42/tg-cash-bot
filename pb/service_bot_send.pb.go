// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: service_bot_send.proto

package pb

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_service_bot_send_proto protoreflect.FileDescriptor

var file_service_bot_send_proto_rawDesc = []byte{
	0x0a, 0x16, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x62, 0x6f, 0x74, 0x5f, 0x73, 0x65,
	0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32,
	0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x73, 0x65, 0x6e, 0x64,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x9a,
	0x01, 0x0a, 0x07, 0x42, 0x6f, 0x74, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x8e, 0x01, 0x0a, 0x0a, 0x53,
	0x65, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x2e, 0x70, 0x62, 0x2e, 0x53,
	0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x17, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x4f, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x14, 0x22, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x72, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x3a, 0x01, 0x2a, 0x92, 0x41, 0x32, 0x12, 0x17, 0x53, 0x75, 0x6d, 0x6d, 0x61,
	0x72, 0x79, 0x3a, 0x20, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x20, 0x72,
	0x70, 0x63, 0x1a, 0x17, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x20,
	0x53, 0x65, 0x6e, 0x64, 0x20, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x42, 0x99, 0x01, 0x5a, 0x29,
	0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x6f, 0x7a, 0x6f, 0x6e, 0x2e, 0x64, 0x65, 0x76, 0x2f,
	0x6d, 0x69, 0x72, 0x6f, 0x6d, 0x61, 0x78, 0x78, 0x73, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x67, 0x72,
	0x61, 0x6d, 0x2d, 0x62, 0x6f, 0x74, 0x2f, 0x70, 0x62, 0x92, 0x41, 0x6b, 0x12, 0x69, 0x0a, 0x18,
	0x54, 0x65, 0x6c, 0x65, 0x67, 0x72, 0x61, 0x6d, 0x20, 0x65, 0x78, 0x70, 0x65, 0x6e, 0x73, 0x65,
	0x2d, 0x62, 0x6f, 0x74, 0x20, 0x41, 0x50, 0x49, 0x22, 0x48, 0x0a, 0x14, 0x4d, 0x61, 0x78, 0x69,
	0x6d, 0x20, 0x4d, 0x69, 0x72, 0x6f, 0x73, 0x68, 0x69, 0x63, 0x68, 0x65, 0x6e, 0x63, 0x6b, 0x6f,
	0x12, 0x1b, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x6d, 0x69, 0x72, 0x6f, 0x6d, 0x61,
	0x78, 0x34, 0x32, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6f, 0x1a, 0x13, 0x6d,
	0x69, 0x72, 0x6f, 0x6d, 0x61, 0x78, 0x78, 0x73, 0x40, 0x67, 0x6d, 0x61, 0x69, 0x6c, 0x2e, 0x63,
	0x6f, 0x6d, 0x32, 0x03, 0x30, 0x2e, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_service_bot_send_proto_goTypes = []interface{}{
	(*SendMessageRequest)(nil),  // 0: pb.SendMessageRequest
	(*SendMessageResponse)(nil), // 1: pb.SendMessageResponse
}
var file_service_bot_send_proto_depIdxs = []int32{
	0, // 0: pb.BotSend.SendReport:input_type -> pb.SendMessageRequest
	1, // 1: pb.BotSend.SendReport:output_type -> pb.SendMessageResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_bot_send_proto_init() }
func file_service_bot_send_proto_init() {
	if File_service_bot_send_proto != nil {
		return
	}
	file_send_message_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_bot_send_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_bot_send_proto_goTypes,
		DependencyIndexes: file_service_bot_send_proto_depIdxs,
	}.Build()
	File_service_bot_send_proto = out.File
	file_service_bot_send_proto_rawDesc = nil
	file_service_bot_send_proto_goTypes = nil
	file_service_bot_send_proto_depIdxs = nil
}
