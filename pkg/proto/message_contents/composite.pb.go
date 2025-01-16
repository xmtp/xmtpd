// Composite ContentType

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        (unknown)
// source: message_contents/composite.proto

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

// Composite is used to implement xmtp.org/composite content type
type Composite struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Parts         []*Composite_Part      `protobuf:"bytes,1,rep,name=parts,proto3" json:"parts,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Composite) Reset() {
	*x = Composite{}
	mi := &file_message_contents_composite_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Composite) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Composite) ProtoMessage() {}

func (x *Composite) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_composite_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Composite.ProtoReflect.Descriptor instead.
func (*Composite) Descriptor() ([]byte, []int) {
	return file_message_contents_composite_proto_rawDescGZIP(), []int{0}
}

func (x *Composite) GetParts() []*Composite_Part {
	if x != nil {
		return x.Parts
	}
	return nil
}

// Part represents one section of a composite message
type Composite_Part struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Element:
	//
	//	*Composite_Part_Part
	//	*Composite_Part_Composite
	Element       isComposite_Part_Element `protobuf_oneof:"element"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Composite_Part) Reset() {
	*x = Composite_Part{}
	mi := &file_message_contents_composite_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Composite_Part) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Composite_Part) ProtoMessage() {}

func (x *Composite_Part) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_composite_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Composite_Part.ProtoReflect.Descriptor instead.
func (*Composite_Part) Descriptor() ([]byte, []int) {
	return file_message_contents_composite_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Composite_Part) GetElement() isComposite_Part_Element {
	if x != nil {
		return x.Element
	}
	return nil
}

func (x *Composite_Part) GetPart() *EncodedContent {
	if x != nil {
		if x, ok := x.Element.(*Composite_Part_Part); ok {
			return x.Part
		}
	}
	return nil
}

func (x *Composite_Part) GetComposite() *Composite {
	if x != nil {
		if x, ok := x.Element.(*Composite_Part_Composite); ok {
			return x.Composite
		}
	}
	return nil
}

type isComposite_Part_Element interface {
	isComposite_Part_Element()
}

type Composite_Part_Part struct {
	Part *EncodedContent `protobuf:"bytes,1,opt,name=part,proto3,oneof"`
}

type Composite_Part_Composite struct {
	Composite *Composite `protobuf:"bytes,2,opt,name=composite,proto3,oneof"`
}

func (*Composite_Part_Part) isComposite_Part_Element() {}

func (*Composite_Part_Composite) isComposite_Part_Element() {}

var File_message_contents_composite_proto protoreflect.FileDescriptor

var file_message_contents_composite_proto_rawDesc = []byte{
	0x0a, 0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x15, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x1e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdb, 0x01, 0x0a, 0x09, 0x43, 0x6f,
	0x6d, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x65, 0x12, 0x3b, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x43,
	0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x52, 0x05, 0x70,
	0x61, 0x72, 0x74, 0x73, 0x1a, 0x90, 0x01, 0x0a, 0x04, 0x50, 0x61, 0x72, 0x74, 0x12, 0x3b, 0x0a,
	0x04, 0x70, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x78, 0x6d,
	0x74, 0x70, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x48, 0x00, 0x52, 0x04, 0x70, 0x61, 0x72, 0x74, 0x12, 0x40, 0x0a, 0x09, 0x63, 0x6f,
	0x6d, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e,
	0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x65, 0x48,
	0x00, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x65, 0x42, 0x09, 0x0a, 0x07,
	0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x42, 0xce, 0x01, 0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e,
	0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x73, 0x42, 0x0e, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x65,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x6d, 0x74, 0x70, 0x2f, 0x78, 0x6d, 0x74, 0x70, 0x64, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0xa2, 0x02, 0x03, 0x58, 0x4d, 0x58, 0xaa,
	0x02, 0x14, 0x58, 0x6d, 0x74, 0x70, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0xca, 0x02, 0x14, 0x58, 0x6d, 0x74, 0x70, 0x5c, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0xe2, 0x02, 0x20,
	0x58, 0x6d, 0x74, 0x70, 0x5c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x15, 0x58, 0x6d, 0x74, 0x70, 0x3a, 0x3a, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_contents_composite_proto_rawDescOnce sync.Once
	file_message_contents_composite_proto_rawDescData = file_message_contents_composite_proto_rawDesc
)

func file_message_contents_composite_proto_rawDescGZIP() []byte {
	file_message_contents_composite_proto_rawDescOnce.Do(func() {
		file_message_contents_composite_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_contents_composite_proto_rawDescData)
	})
	return file_message_contents_composite_proto_rawDescData
}

var file_message_contents_composite_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_message_contents_composite_proto_goTypes = []any{
	(*Composite)(nil),      // 0: xmtp.message_contents.Composite
	(*Composite_Part)(nil), // 1: xmtp.message_contents.Composite.Part
	(*EncodedContent)(nil), // 2: xmtp.message_contents.EncodedContent
}
var file_message_contents_composite_proto_depIdxs = []int32{
	1, // 0: xmtp.message_contents.Composite.parts:type_name -> xmtp.message_contents.Composite.Part
	2, // 1: xmtp.message_contents.Composite.Part.part:type_name -> xmtp.message_contents.EncodedContent
	0, // 2: xmtp.message_contents.Composite.Part.composite:type_name -> xmtp.message_contents.Composite
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_message_contents_composite_proto_init() }
func file_message_contents_composite_proto_init() {
	if File_message_contents_composite_proto != nil {
		return
	}
	file_message_contents_content_proto_init()
	file_message_contents_composite_proto_msgTypes[1].OneofWrappers = []any{
		(*Composite_Part_Part)(nil),
		(*Composite_Part_Composite)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_contents_composite_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_contents_composite_proto_goTypes,
		DependencyIndexes: file_message_contents_composite_proto_depIdxs,
		MessageInfos:      file_message_contents_composite_proto_msgTypes,
	}.Build()
	File_message_contents_composite_proto = out.File
	file_message_contents_composite_proto_rawDesc = nil
	file_message_contents_composite_proto_goTypes = nil
	file_message_contents_composite_proto_depIdxs = nil
}
