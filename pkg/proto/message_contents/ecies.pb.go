// ECIES is a wrapper for ECIES payloads

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: message_contents/ecies.proto

package message_contents

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

// EciesMessage is a wrapper for ECIES encrypted payloads
type EciesMessage struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Version:
	//
	//	*EciesMessage_V1
	Version       isEciesMessage_Version `protobuf_oneof:"version"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EciesMessage) Reset() {
	*x = EciesMessage{}
	mi := &file_message_contents_ecies_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EciesMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EciesMessage) ProtoMessage() {}

func (x *EciesMessage) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_ecies_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EciesMessage.ProtoReflect.Descriptor instead.
func (*EciesMessage) Descriptor() ([]byte, []int) {
	return file_message_contents_ecies_proto_rawDescGZIP(), []int{0}
}

func (x *EciesMessage) GetVersion() isEciesMessage_Version {
	if x != nil {
		return x.Version
	}
	return nil
}

func (x *EciesMessage) GetV1() []byte {
	if x != nil {
		if x, ok := x.Version.(*EciesMessage_V1); ok {
			return x.V1
		}
	}
	return nil
}

type isEciesMessage_Version interface {
	isEciesMessage_Version()
}

type EciesMessage_V1 struct {
	// Expected to be an ECIES encrypted SignedPayload
	V1 []byte `protobuf:"bytes,1,opt,name=v1,proto3,oneof"`
}

func (*EciesMessage_V1) isEciesMessage_Version() {}

var File_message_contents_ecies_proto protoreflect.FileDescriptor

var file_message_contents_ecies_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x73, 0x2f, 0x65, 0x63, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15,
	0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x2b, 0x0a, 0x0c, 0x45, 0x63, 0x69, 0x65, 0x73, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x02, 0x76, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x48, 0x00, 0x52, 0x02, 0x76, 0x31, 0x42, 0x09, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x42, 0xca, 0x01, 0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73,
	0x42, 0x0a, 0x45, 0x63, 0x69, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x30,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x6d, 0x74, 0x70, 0x2f,
	0x78, 0x6d, 0x74, 0x70, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73,
	0xa2, 0x02, 0x03, 0x58, 0x4d, 0x58, 0xaa, 0x02, 0x14, 0x58, 0x6d, 0x74, 0x70, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0xca, 0x02, 0x14,
	0x58, 0x6d, 0x74, 0x70, 0x5c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x73, 0xe2, 0x02, 0x20, 0x58, 0x6d, 0x74, 0x70, 0x5c, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x15, 0x58, 0x6d, 0x74, 0x70, 0x3a, 0x3a,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_contents_ecies_proto_rawDescOnce sync.Once
	file_message_contents_ecies_proto_rawDescData = file_message_contents_ecies_proto_rawDesc
)

func file_message_contents_ecies_proto_rawDescGZIP() []byte {
	file_message_contents_ecies_proto_rawDescOnce.Do(func() {
		file_message_contents_ecies_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_contents_ecies_proto_rawDescData)
	})
	return file_message_contents_ecies_proto_rawDescData
}

var file_message_contents_ecies_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_message_contents_ecies_proto_goTypes = []any{
	(*EciesMessage)(nil), // 0: xmtp.message_contents.EciesMessage
}
var file_message_contents_ecies_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_message_contents_ecies_proto_init() }
func file_message_contents_ecies_proto_init() {
	if File_message_contents_ecies_proto != nil {
		return
	}
	file_message_contents_ecies_proto_msgTypes[0].OneofWrappers = []any{
		(*EciesMessage_V1)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_contents_ecies_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_contents_ecies_proto_goTypes,
		DependencyIndexes: file_message_contents_ecies_proto_depIdxs,
		MessageInfos:      file_message_contents_ecies_proto_msgTypes,
	}.Build()
	File_message_contents_ecies_proto = out.File
	file_message_contents_ecies_proto_rawDesc = nil
	file_message_contents_ecies_proto_goTypes = nil
	file_message_contents_ecies_proto_depIdxs = nil
}
