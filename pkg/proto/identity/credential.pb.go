// Credentials

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: identity/credential.proto

package identity

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// A credential that can be used in MLS leaf nodes
type MlsCredential struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	InboxId       string                 `protobuf:"bytes,1,opt,name=inbox_id,json=inboxId,proto3" json:"inbox_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MlsCredential) Reset() {
	*x = MlsCredential{}
	mi := &file_identity_credential_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MlsCredential) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MlsCredential) ProtoMessage() {}

func (x *MlsCredential) ProtoReflect() protoreflect.Message {
	mi := &file_identity_credential_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MlsCredential.ProtoReflect.Descriptor instead.
func (*MlsCredential) Descriptor() ([]byte, []int) {
	return file_identity_credential_proto_rawDescGZIP(), []int{0}
}

func (x *MlsCredential) GetInboxId() string {
	if x != nil {
		return x.InboxId
	}
	return ""
}

var File_identity_credential_proto protoreflect.FileDescriptor

const file_identity_credential_proto_rawDesc = "" +
	"\n" +
	"\x19identity/credential.proto\x12\rxmtp.identity\"*\n" +
	"\rMlsCredential\x12\x19\n" +
	"\binbox_id\x18\x01 \x01(\tR\ainboxIdB\xa3\x01\n" +
	"\x11com.xmtp.identityB\x0fCredentialProtoP\x01Z(github.com/xmtp/xmtpd/pkg/proto/identity\xa2\x02\x03XIX\xaa\x02\rXmtp.Identity\xca\x02\rXmtp\\Identity\xe2\x02\x19Xmtp\\Identity\\GPBMetadata\xea\x02\x0eXmtp::Identityb\x06proto3"

var (
	file_identity_credential_proto_rawDescOnce sync.Once
	file_identity_credential_proto_rawDescData []byte
)

func file_identity_credential_proto_rawDescGZIP() []byte {
	file_identity_credential_proto_rawDescOnce.Do(func() {
		file_identity_credential_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_identity_credential_proto_rawDesc), len(file_identity_credential_proto_rawDesc)))
	})
	return file_identity_credential_proto_rawDescData
}

var file_identity_credential_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_identity_credential_proto_goTypes = []any{
	(*MlsCredential)(nil), // 0: xmtp.identity.MlsCredential
}
var file_identity_credential_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_identity_credential_proto_init() }
func file_identity_credential_proto_init() {
	if File_identity_credential_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_identity_credential_proto_rawDesc), len(file_identity_credential_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_identity_credential_proto_goTypes,
		DependencyIndexes: file_identity_credential_proto_depIdxs,
		MessageInfos:      file_identity_credential_proto_msgTypes,
	}.Build()
	File_identity_credential_proto = out.File
	file_identity_credential_proto_goTypes = nil
	file_identity_credential_proto_depIdxs = nil
}
