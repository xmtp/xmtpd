// Signature is a generic structure for public key signatures.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: message_contents/frames.proto

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

// The message that will be signed by the Client and returned inside the
// `action_body` field of the FrameAction message
type FrameActionBody struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The URL of the frame that was clicked
	// May be different from `post_url`
	FrameUrl string `protobuf:"bytes,1,opt,name=frame_url,json=frameUrl,proto3" json:"frame_url,omitempty"`
	// The 1-indexed button that was clicked
	ButtonIndex int32 `protobuf:"varint,2,opt,name=button_index,json=buttonIndex,proto3" json:"button_index,omitempty"`
	// Timestamp of the click in milliseconds since the epoch
	//
	// Deprecated: Marked as deprecated in message_contents/frames.proto.
	Timestamp uint64 `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// A unique identifier for the conversation, not tied to anything on the
	// network. Will not match the topic or conversation_id
	OpaqueConversationIdentifier string `protobuf:"bytes,4,opt,name=opaque_conversation_identifier,json=opaqueConversationIdentifier,proto3" json:"opaque_conversation_identifier,omitempty"`
	// Unix timestamp
	UnixTimestamp uint32 `protobuf:"varint,5,opt,name=unix_timestamp,json=unixTimestamp,proto3" json:"unix_timestamp,omitempty"`
	// Input text from a text input field
	InputText string `protobuf:"bytes,6,opt,name=input_text,json=inputText,proto3" json:"input_text,omitempty"`
	// A state serialized to a string (for example via JSON.stringify()). Maximum 4096 bytes.
	State string `protobuf:"bytes,7,opt,name=state,proto3" json:"state,omitempty"`
	// A 0x wallet address
	Address string `protobuf:"bytes,8,opt,name=address,proto3" json:"address,omitempty"`
	// A hash from a transaction
	TransactionId string `protobuf:"bytes,9,opt,name=transaction_id,json=transactionId,proto3" json:"transaction_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FrameActionBody) Reset() {
	*x = FrameActionBody{}
	mi := &file_message_contents_frames_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FrameActionBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FrameActionBody) ProtoMessage() {}

func (x *FrameActionBody) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_frames_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FrameActionBody.ProtoReflect.Descriptor instead.
func (*FrameActionBody) Descriptor() ([]byte, []int) {
	return file_message_contents_frames_proto_rawDescGZIP(), []int{0}
}

func (x *FrameActionBody) GetFrameUrl() string {
	if x != nil {
		return x.FrameUrl
	}
	return ""
}

func (x *FrameActionBody) GetButtonIndex() int32 {
	if x != nil {
		return x.ButtonIndex
	}
	return 0
}

// Deprecated: Marked as deprecated in message_contents/frames.proto.
func (x *FrameActionBody) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *FrameActionBody) GetOpaqueConversationIdentifier() string {
	if x != nil {
		return x.OpaqueConversationIdentifier
	}
	return ""
}

func (x *FrameActionBody) GetUnixTimestamp() uint32 {
	if x != nil {
		return x.UnixTimestamp
	}
	return 0
}

func (x *FrameActionBody) GetInputText() string {
	if x != nil {
		return x.InputText
	}
	return ""
}

func (x *FrameActionBody) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *FrameActionBody) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *FrameActionBody) GetTransactionId() string {
	if x != nil {
		return x.TransactionId
	}
	return ""
}

// The outer payload that will be sent as the `messageBytes` in the
// `trusted_data` part of the Frames message
type FrameAction struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Deprecated: Marked as deprecated in message_contents/frames.proto.
	Signature *Signature `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	// The SignedPublicKeyBundle of the signer, used to link the XMTP signature
	// with a blockchain account through a chain of signatures.
	//
	// Deprecated: Marked as deprecated in message_contents/frames.proto.
	SignedPublicKeyBundle *SignedPublicKeyBundle `protobuf:"bytes,2,opt,name=signed_public_key_bundle,json=signedPublicKeyBundle,proto3" json:"signed_public_key_bundle,omitempty"`
	// Serialized FrameActionBody message, so that the signature verification can
	// happen on a byte-perfect representation of the message
	ActionBody []byte `protobuf:"bytes,3,opt,name=action_body,json=actionBody,proto3" json:"action_body,omitempty"`
	// The installation signature
	InstallationSignature []byte `protobuf:"bytes,4,opt,name=installation_signature,json=installationSignature,proto3" json:"installation_signature,omitempty"`
	// The public installation id used to sign.
	InstallationId []byte `protobuf:"bytes,5,opt,name=installation_id,json=installationId,proto3" json:"installation_id,omitempty"`
	// The inbox id of the installation used to sign.
	InboxId       string `protobuf:"bytes,6,opt,name=inbox_id,json=inboxId,proto3" json:"inbox_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FrameAction) Reset() {
	*x = FrameAction{}
	mi := &file_message_contents_frames_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FrameAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FrameAction) ProtoMessage() {}

func (x *FrameAction) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_frames_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FrameAction.ProtoReflect.Descriptor instead.
func (*FrameAction) Descriptor() ([]byte, []int) {
	return file_message_contents_frames_proto_rawDescGZIP(), []int{1}
}

// Deprecated: Marked as deprecated in message_contents/frames.proto.
func (x *FrameAction) GetSignature() *Signature {
	if x != nil {
		return x.Signature
	}
	return nil
}

// Deprecated: Marked as deprecated in message_contents/frames.proto.
func (x *FrameAction) GetSignedPublicKeyBundle() *SignedPublicKeyBundle {
	if x != nil {
		return x.SignedPublicKeyBundle
	}
	return nil
}

func (x *FrameAction) GetActionBody() []byte {
	if x != nil {
		return x.ActionBody
	}
	return nil
}

func (x *FrameAction) GetInstallationSignature() []byte {
	if x != nil {
		return x.InstallationSignature
	}
	return nil
}

func (x *FrameAction) GetInstallationId() []byte {
	if x != nil {
		return x.InstallationId
	}
	return nil
}

func (x *FrameAction) GetInboxId() string {
	if x != nil {
		return x.InboxId
	}
	return ""
}

var File_message_contents_frames_proto protoreflect.FileDescriptor

var file_message_contents_frames_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x73, 0x2f, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x15, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x21, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f,
	0x6b, 0x65, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd6, 0x02, 0x0a, 0x0f,
	0x46, 0x72, 0x61, 0x6d, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x6f, 0x64, 0x79, 0x12,
	0x1b, 0x0a, 0x09, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x21, 0x0a, 0x0c,
	0x62, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0b, 0x62, 0x75, 0x74, 0x74, 0x6f, 0x6e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12,
	0x20, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x04, 0x42, 0x02, 0x18, 0x01, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x44, 0x0a, 0x1e, 0x6f, 0x70, 0x61, 0x71, 0x75, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x76,
	0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1c, 0x6f, 0x70, 0x61, 0x71, 0x75,
	0x65, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x0e, 0x75, 0x6e, 0x69, 0x78, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0d, 0x75, 0x6e, 0x69, 0x78, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1d,
	0x0a, 0x0a, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x54, 0x65, 0x78, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x25, 0x0a,
	0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x22, 0xd8, 0x02, 0x0a, 0x0b, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x42, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x2e,
	0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x42, 0x02, 0x18, 0x01, 0x52, 0x09, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x69, 0x0a, 0x18, 0x73, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x62, 0x75,
	0x6e, 0x64, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x78, 0x6d, 0x74,
	0x70, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x73, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b,
	0x65, 0x79, 0x42, 0x75, 0x6e, 0x64, 0x6c, 0x65, 0x42, 0x02, 0x18, 0x01, 0x52, 0x15, 0x73, 0x69,
	0x67, 0x6e, 0x65, 0x64, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x42, 0x75, 0x6e,
	0x64, 0x6c, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x62, 0x6f,
	0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x42, 0x6f, 0x64, 0x79, 0x12, 0x35, 0x0a, 0x16, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x15, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x69,
	0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0e, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x6e, 0x62, 0x6f, 0x78, 0x5f, 0x69, 0x64,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x69, 0x6e, 0x62, 0x6f, 0x78, 0x49, 0x64, 0x42,
	0xcb, 0x01, 0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x42, 0x0b, 0x46,
	0x72, 0x61, 0x6d, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x30, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x6d, 0x74, 0x70, 0x2f, 0x78, 0x6d,
	0x74, 0x70, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0xa2, 0x02,
	0x03, 0x58, 0x4d, 0x58, 0xaa, 0x02, 0x14, 0x58, 0x6d, 0x74, 0x70, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0xca, 0x02, 0x14, 0x58, 0x6d,
	0x74, 0x70, 0x5c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x73, 0xe2, 0x02, 0x20, 0x58, 0x6d, 0x74, 0x70, 0x5c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x15, 0x58, 0x6d, 0x74, 0x70, 0x3a, 0x3a, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_contents_frames_proto_rawDescOnce sync.Once
	file_message_contents_frames_proto_rawDescData = file_message_contents_frames_proto_rawDesc
)

func file_message_contents_frames_proto_rawDescGZIP() []byte {
	file_message_contents_frames_proto_rawDescOnce.Do(func() {
		file_message_contents_frames_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_contents_frames_proto_rawDescData)
	})
	return file_message_contents_frames_proto_rawDescData
}

var file_message_contents_frames_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_message_contents_frames_proto_goTypes = []any{
	(*FrameActionBody)(nil),       // 0: xmtp.message_contents.FrameActionBody
	(*FrameAction)(nil),           // 1: xmtp.message_contents.FrameAction
	(*Signature)(nil),             // 2: xmtp.message_contents.Signature
	(*SignedPublicKeyBundle)(nil), // 3: xmtp.message_contents.SignedPublicKeyBundle
}
var file_message_contents_frames_proto_depIdxs = []int32{
	2, // 0: xmtp.message_contents.FrameAction.signature:type_name -> xmtp.message_contents.Signature
	3, // 1: xmtp.message_contents.FrameAction.signed_public_key_bundle:type_name -> xmtp.message_contents.SignedPublicKeyBundle
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_message_contents_frames_proto_init() }
func file_message_contents_frames_proto_init() {
	if File_message_contents_frames_proto != nil {
		return
	}
	file_message_contents_public_key_proto_init()
	file_message_contents_signature_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_message_contents_frames_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_contents_frames_proto_goTypes,
		DependencyIndexes: file_message_contents_frames_proto_depIdxs,
		MessageInfos:      file_message_contents_frames_proto_msgTypes,
	}.Build()
	File_message_contents_frames_proto = out.File
	file_message_contents_frames_proto_rawDesc = nil
	file_message_contents_frames_proto_goTypes = nil
	file_message_contents_frames_proto_depIdxs = nil
}
