// Message API for XMTP V4

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: xmtpv4/envelopes/envelopes.proto

package envelopes

import (
	associations "github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	v1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
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

// The last seen entry per originator. Originators that have not been seen are omitted.
type Cursor struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	NodeIdToSequenceId map[uint32]uint64      `protobuf:"bytes,1,rep,name=node_id_to_sequence_id,json=nodeIdToSequenceId,proto3" json:"node_id_to_sequence_id,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *Cursor) Reset() {
	*x = Cursor{}
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cursor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cursor) ProtoMessage() {}

func (x *Cursor) ProtoReflect() protoreflect.Message {
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cursor.ProtoReflect.Descriptor instead.
func (*Cursor) Descriptor() ([]byte, []int) {
	return file_xmtpv4_envelopes_envelopes_proto_rawDescGZIP(), []int{0}
}

func (x *Cursor) GetNodeIdToSequenceId() map[uint32]uint64 {
	if x != nil {
		return x.NodeIdToSequenceId
	}
	return nil
}

// Data visible to the server that has been authenticated by the client.
type AuthenticatedData struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Deprecated: Marked as deprecated in xmtpv4/envelopes/envelopes.proto.
	TargetOriginator *uint32 `protobuf:"varint,1,opt,name=target_originator,json=targetOriginator,proto3,oneof" json:"target_originator,omitempty"`
	TargetTopic      []byte  `protobuf:"bytes,2,opt,name=target_topic,json=targetTopic,proto3" json:"target_topic,omitempty"`
	DependsOn        *Cursor `protobuf:"bytes,3,opt,name=depends_on,json=dependsOn,proto3" json:"depends_on,omitempty"`
	IsCommit         bool    `protobuf:"varint,4,opt,name=is_commit,json=isCommit,proto3" json:"is_commit,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *AuthenticatedData) Reset() {
	*x = AuthenticatedData{}
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthenticatedData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticatedData) ProtoMessage() {}

func (x *AuthenticatedData) ProtoReflect() protoreflect.Message {
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticatedData.ProtoReflect.Descriptor instead.
func (*AuthenticatedData) Descriptor() ([]byte, []int) {
	return file_xmtpv4_envelopes_envelopes_proto_rawDescGZIP(), []int{1}
}

// Deprecated: Marked as deprecated in xmtpv4/envelopes/envelopes.proto.
func (x *AuthenticatedData) GetTargetOriginator() uint32 {
	if x != nil && x.TargetOriginator != nil {
		return *x.TargetOriginator
	}
	return 0
}

func (x *AuthenticatedData) GetTargetTopic() []byte {
	if x != nil {
		return x.TargetTopic
	}
	return nil
}

func (x *AuthenticatedData) GetDependsOn() *Cursor {
	if x != nil {
		return x.DependsOn
	}
	return nil
}

func (x *AuthenticatedData) GetIsCommit() bool {
	if x != nil {
		return x.IsCommit
	}
	return false
}

type ClientEnvelope struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Aad   *AuthenticatedData     `protobuf:"bytes,1,opt,name=aad,proto3" json:"aad,omitempty"`
	// Types that are valid to be assigned to Payload:
	//
	//	*ClientEnvelope_GroupMessage
	//	*ClientEnvelope_WelcomeMessage
	//	*ClientEnvelope_UploadKeyPackage
	//	*ClientEnvelope_IdentityUpdate
	Payload       isClientEnvelope_Payload `protobuf_oneof:"payload"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ClientEnvelope) Reset() {
	*x = ClientEnvelope{}
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClientEnvelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientEnvelope) ProtoMessage() {}

func (x *ClientEnvelope) ProtoReflect() protoreflect.Message {
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientEnvelope.ProtoReflect.Descriptor instead.
func (*ClientEnvelope) Descriptor() ([]byte, []int) {
	return file_xmtpv4_envelopes_envelopes_proto_rawDescGZIP(), []int{2}
}

func (x *ClientEnvelope) GetAad() *AuthenticatedData {
	if x != nil {
		return x.Aad
	}
	return nil
}

func (x *ClientEnvelope) GetPayload() isClientEnvelope_Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *ClientEnvelope) GetGroupMessage() *v1.GroupMessageInput {
	if x != nil {
		if x, ok := x.Payload.(*ClientEnvelope_GroupMessage); ok {
			return x.GroupMessage
		}
	}
	return nil
}

func (x *ClientEnvelope) GetWelcomeMessage() *v1.WelcomeMessageInput {
	if x != nil {
		if x, ok := x.Payload.(*ClientEnvelope_WelcomeMessage); ok {
			return x.WelcomeMessage
		}
	}
	return nil
}

func (x *ClientEnvelope) GetUploadKeyPackage() *v1.UploadKeyPackageRequest {
	if x != nil {
		if x, ok := x.Payload.(*ClientEnvelope_UploadKeyPackage); ok {
			return x.UploadKeyPackage
		}
	}
	return nil
}

func (x *ClientEnvelope) GetIdentityUpdate() *associations.IdentityUpdate {
	if x != nil {
		if x, ok := x.Payload.(*ClientEnvelope_IdentityUpdate); ok {
			return x.IdentityUpdate
		}
	}
	return nil
}

type isClientEnvelope_Payload interface {
	isClientEnvelope_Payload()
}

type ClientEnvelope_GroupMessage struct {
	GroupMessage *v1.GroupMessageInput `protobuf:"bytes,2,opt,name=group_message,json=groupMessage,proto3,oneof"`
}

type ClientEnvelope_WelcomeMessage struct {
	WelcomeMessage *v1.WelcomeMessageInput `protobuf:"bytes,3,opt,name=welcome_message,json=welcomeMessage,proto3,oneof"`
}

type ClientEnvelope_UploadKeyPackage struct {
	UploadKeyPackage *v1.UploadKeyPackageRequest `protobuf:"bytes,4,opt,name=upload_key_package,json=uploadKeyPackage,proto3,oneof"`
}

type ClientEnvelope_IdentityUpdate struct {
	IdentityUpdate *associations.IdentityUpdate `protobuf:"bytes,5,opt,name=identity_update,json=identityUpdate,proto3,oneof"`
}

func (*ClientEnvelope_GroupMessage) isClientEnvelope_Payload() {}

func (*ClientEnvelope_WelcomeMessage) isClientEnvelope_Payload() {}

func (*ClientEnvelope_UploadKeyPackage) isClientEnvelope_Payload() {}

func (*ClientEnvelope_IdentityUpdate) isClientEnvelope_Payload() {}

// Wraps client envelope with payer signature
type PayerEnvelope struct {
	state                  protoimpl.MessageState                  `protogen:"open.v1"`
	UnsignedClientEnvelope []byte                                  `protobuf:"bytes,1,opt,name=unsigned_client_envelope,json=unsignedClientEnvelope,proto3" json:"unsigned_client_envelope,omitempty"` // Protobuf serialized
	PayerSignature         *associations.RecoverableEcdsaSignature `protobuf:"bytes,2,opt,name=payer_signature,json=payerSignature,proto3" json:"payer_signature,omitempty"`
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *PayerEnvelope) Reset() {
	*x = PayerEnvelope{}
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PayerEnvelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayerEnvelope) ProtoMessage() {}

func (x *PayerEnvelope) ProtoReflect() protoreflect.Message {
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayerEnvelope.ProtoReflect.Descriptor instead.
func (*PayerEnvelope) Descriptor() ([]byte, []int) {
	return file_xmtpv4_envelopes_envelopes_proto_rawDescGZIP(), []int{3}
}

func (x *PayerEnvelope) GetUnsignedClientEnvelope() []byte {
	if x != nil {
		return x.UnsignedClientEnvelope
	}
	return nil
}

func (x *PayerEnvelope) GetPayerSignature() *associations.RecoverableEcdsaSignature {
	if x != nil {
		return x.PayerSignature
	}
	return nil
}

// For blockchain envelopes, these fields are set by the smart contract
type UnsignedOriginatorEnvelope struct {
	state                protoimpl.MessageState `protogen:"open.v1"`
	OriginatorNodeId     uint32                 `protobuf:"varint,1,opt,name=originator_node_id,json=originatorNodeId,proto3" json:"originator_node_id,omitempty"`
	OriginatorSequenceId uint64                 `protobuf:"varint,2,opt,name=originator_sequence_id,json=originatorSequenceId,proto3" json:"originator_sequence_id,omitempty"`
	OriginatorNs         int64                  `protobuf:"varint,3,opt,name=originator_ns,json=originatorNs,proto3" json:"originator_ns,omitempty"`
	PayerEnvelope        *PayerEnvelope         `protobuf:"bytes,4,opt,name=payer_envelope,json=payerEnvelope,proto3" json:"payer_envelope,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *UnsignedOriginatorEnvelope) Reset() {
	*x = UnsignedOriginatorEnvelope{}
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnsignedOriginatorEnvelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnsignedOriginatorEnvelope) ProtoMessage() {}

func (x *UnsignedOriginatorEnvelope) ProtoReflect() protoreflect.Message {
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnsignedOriginatorEnvelope.ProtoReflect.Descriptor instead.
func (*UnsignedOriginatorEnvelope) Descriptor() ([]byte, []int) {
	return file_xmtpv4_envelopes_envelopes_proto_rawDescGZIP(), []int{4}
}

func (x *UnsignedOriginatorEnvelope) GetOriginatorNodeId() uint32 {
	if x != nil {
		return x.OriginatorNodeId
	}
	return 0
}

func (x *UnsignedOriginatorEnvelope) GetOriginatorSequenceId() uint64 {
	if x != nil {
		return x.OriginatorSequenceId
	}
	return 0
}

func (x *UnsignedOriginatorEnvelope) GetOriginatorNs() int64 {
	if x != nil {
		return x.OriginatorNs
	}
	return 0
}

func (x *UnsignedOriginatorEnvelope) GetPayerEnvelope() *PayerEnvelope {
	if x != nil {
		return x.PayerEnvelope
	}
	return nil
}

// An alternative to a signature for blockchain payloads
type BlockchainProof struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	TransactionHash []byte                 `protobuf:"bytes,1,opt,name=transaction_hash,json=transactionHash,proto3" json:"transaction_hash,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *BlockchainProof) Reset() {
	*x = BlockchainProof{}
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BlockchainProof) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockchainProof) ProtoMessage() {}

func (x *BlockchainProof) ProtoReflect() protoreflect.Message {
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockchainProof.ProtoReflect.Descriptor instead.
func (*BlockchainProof) Descriptor() ([]byte, []int) {
	return file_xmtpv4_envelopes_envelopes_proto_rawDescGZIP(), []int{5}
}

func (x *BlockchainProof) GetTransactionHash() []byte {
	if x != nil {
		return x.TransactionHash
	}
	return nil
}

// Signed originator envelope
type OriginatorEnvelope struct {
	state                      protoimpl.MessageState `protogen:"open.v1"`
	UnsignedOriginatorEnvelope []byte                 `protobuf:"bytes,1,opt,name=unsigned_originator_envelope,json=unsignedOriginatorEnvelope,proto3" json:"unsigned_originator_envelope,omitempty"` // Protobuf serialized
	// Types that are valid to be assigned to Proof:
	//
	//	*OriginatorEnvelope_OriginatorSignature
	//	*OriginatorEnvelope_BlockchainProof
	Proof         isOriginatorEnvelope_Proof `protobuf_oneof:"proof"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OriginatorEnvelope) Reset() {
	*x = OriginatorEnvelope{}
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OriginatorEnvelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OriginatorEnvelope) ProtoMessage() {}

func (x *OriginatorEnvelope) ProtoReflect() protoreflect.Message {
	mi := &file_xmtpv4_envelopes_envelopes_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OriginatorEnvelope.ProtoReflect.Descriptor instead.
func (*OriginatorEnvelope) Descriptor() ([]byte, []int) {
	return file_xmtpv4_envelopes_envelopes_proto_rawDescGZIP(), []int{6}
}

func (x *OriginatorEnvelope) GetUnsignedOriginatorEnvelope() []byte {
	if x != nil {
		return x.UnsignedOriginatorEnvelope
	}
	return nil
}

func (x *OriginatorEnvelope) GetProof() isOriginatorEnvelope_Proof {
	if x != nil {
		return x.Proof
	}
	return nil
}

func (x *OriginatorEnvelope) GetOriginatorSignature() *associations.RecoverableEcdsaSignature {
	if x != nil {
		if x, ok := x.Proof.(*OriginatorEnvelope_OriginatorSignature); ok {
			return x.OriginatorSignature
		}
	}
	return nil
}

func (x *OriginatorEnvelope) GetBlockchainProof() *BlockchainProof {
	if x != nil {
		if x, ok := x.Proof.(*OriginatorEnvelope_BlockchainProof); ok {
			return x.BlockchainProof
		}
	}
	return nil
}

type isOriginatorEnvelope_Proof interface {
	isOriginatorEnvelope_Proof()
}

type OriginatorEnvelope_OriginatorSignature struct {
	OriginatorSignature *associations.RecoverableEcdsaSignature `protobuf:"bytes,2,opt,name=originator_signature,json=originatorSignature,proto3,oneof"`
}

type OriginatorEnvelope_BlockchainProof struct {
	BlockchainProof *BlockchainProof `protobuf:"bytes,3,opt,name=blockchain_proof,json=blockchainProof,proto3,oneof"`
}

func (*OriginatorEnvelope_OriginatorSignature) isOriginatorEnvelope_Proof() {}

func (*OriginatorEnvelope_BlockchainProof) isOriginatorEnvelope_Proof() {}

var File_xmtpv4_envelopes_envelopes_proto protoreflect.FileDescriptor

var file_xmtpv4_envelopes_envelopes_proto_rawDesc = string([]byte{
	0x0a, 0x20, 0x78, 0x6d, 0x74, 0x70, 0x76, 0x34, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70,
	0x65, 0x73, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x15, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x76, 0x34, 0x2e,
	0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x73, 0x1a, 0x27, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x2f, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x25, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x61, 0x73, 0x73,
	0x6f, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x6d, 0x6c, 0x73, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xba, 0x01, 0x0a, 0x06, 0x43, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x12, 0x69, 0x0a, 0x16, 0x6e, 0x6f,
	0x64, 0x65, 0x5f, 0x69, 0x64, 0x5f, 0x74, 0x6f, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x78, 0x6d, 0x74,
	0x70, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x76, 0x34, 0x2e, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70,
	0x65, 0x73, 0x2e, 0x43, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64,
	0x54, 0x6f, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x12, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x54, 0x6f, 0x53, 0x65, 0x71, 0x75, 0x65,
	0x6e, 0x63, 0x65, 0x49, 0x64, 0x1a, 0x45, 0x0a, 0x17, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x54,
	0x6f, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xdd, 0x01, 0x0a,
	0x11, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x34, 0x0a, 0x11, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x6f, 0x72, 0x69,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x02, 0x18,
	0x01, 0x48, 0x00, 0x52, 0x10, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x5f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x3c, 0x0a, 0x0a, 0x64,
	0x65, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x5f, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x76, 0x34, 0x2e, 0x65, 0x6e,
	0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x73, 0x2e, 0x43, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x52, 0x09,
	0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x4f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f,
	0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73,
	0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x42, 0x14, 0x0a, 0x12, 0x5f, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x5f, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x22, 0xa4, 0x03, 0x0a,
	0x0e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12,
	0x3a, 0x0a, 0x03, 0x61, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x78,
	0x6d, 0x74, 0x70, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x76, 0x34, 0x2e, 0x65, 0x6e, 0x76, 0x65, 0x6c,
	0x6f, 0x70, 0x65, 0x73, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x03, 0x61, 0x61, 0x64, 0x12, 0x49, 0x0a, 0x0d, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x22, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x6c, 0x73, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x48, 0x00, 0x52, 0x0c, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x4f, 0x0a, 0x0f, 0x77, 0x65, 0x6c, 0x63, 0x6f, 0x6d,
	0x65, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x24, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x6c, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x57, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x48, 0x00, 0x52, 0x0e, 0x77, 0x65, 0x6c, 0x63, 0x6f, 0x6d, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x58, 0x0a, 0x12, 0x75, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x6c, 0x73, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4b, 0x65, 0x79, 0x50,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52,
	0x10, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4b, 0x65, 0x79, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67,
	0x65, 0x12, 0x55, 0x0a, 0x0f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x78, 0x6d, 0x74,
	0x70, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x61, 0x73, 0x73, 0x6f, 0x63,
	0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x0e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x22, 0xa9, 0x01, 0x0a, 0x0d, 0x50, 0x61, 0x79, 0x65, 0x72, 0x45, 0x6e, 0x76,
	0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x38, 0x0a, 0x18, 0x75, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x16, 0x75, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x65,
	0x64, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12,
	0x5e, 0x0a, 0x0f, 0x70, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x61, 0x62, 0x6c,
	0x65, 0x45, 0x63, 0x64, 0x73, 0x61, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x52,
	0x0e, 0x70, 0x61, 0x79, 0x65, 0x72, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22,
	0xf2, 0x01, 0x0a, 0x1a, 0x55, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4f, 0x72, 0x69, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x2c,
	0x0a, 0x12, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x6e, 0x6f, 0x64,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x10, 0x6f, 0x72, 0x69, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x16,
	0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65,
	0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x14, 0x6f, 0x72,
	0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65,
	0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72,
	0x5f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x6f, 0x72, 0x4e, 0x73, 0x12, 0x4b, 0x0a, 0x0e, 0x70, 0x61, 0x79, 0x65, 0x72,
	0x5f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x24, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x76, 0x34, 0x2e, 0x65, 0x6e,
	0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x73, 0x2e, 0x50, 0x61, 0x79, 0x65, 0x72, 0x45, 0x6e, 0x76,
	0x65, 0x6c, 0x6f, 0x70, 0x65, 0x52, 0x0d, 0x70, 0x61, 0x79, 0x65, 0x72, 0x45, 0x6e, 0x76, 0x65,
	0x6c, 0x6f, 0x70, 0x65, 0x22, 0x3c, 0x0a, 0x0f, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x61,
	0x73, 0x68, 0x22, 0xa0, 0x02, 0x0a, 0x12, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x6f,
	0x72, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x40, 0x0a, 0x1c, 0x75, 0x6e, 0x73,
	0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72,
	0x5f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x1a, 0x75, 0x6e, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x6f, 0x72, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x6a, 0x0a, 0x14, 0x6f,
	0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x78, 0x6d, 0x74, 0x70,
	0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x61, 0x62,
	0x6c, 0x65, 0x45, 0x63, 0x64, 0x73, 0x61, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x48, 0x00, 0x52, 0x13, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x6f, 0x72, 0x53, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x53, 0x0a, 0x10, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x26, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x76, 0x34, 0x2e,
	0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x73, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x63,
	0x68, 0x61, 0x69, 0x6e, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x48, 0x00, 0x52, 0x0f, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x42, 0x07, 0x0a, 0x05,
	0x70, 0x72, 0x6f, 0x6f, 0x66, 0x42, 0xd3, 0x01, 0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x78, 0x6d,
	0x74, 0x70, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x76, 0x34, 0x2e, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f,
	0x70, 0x65, 0x73, 0x42, 0x0e, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x73, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x78, 0x6d, 0x74, 0x70, 0x2f, 0x78, 0x6d, 0x74, 0x70, 0x64, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x78, 0x6d, 0x74, 0x70, 0x76, 0x34, 0x2f, 0x65, 0x6e,
	0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x73, 0xa2, 0x02, 0x03, 0x58, 0x58, 0x45, 0xaa, 0x02, 0x15,
	0x58, 0x6d, 0x74, 0x70, 0x2e, 0x58, 0x6d, 0x74, 0x70, 0x76, 0x34, 0x2e, 0x45, 0x6e, 0x76, 0x65,
	0x6c, 0x6f, 0x70, 0x65, 0x73, 0xca, 0x02, 0x15, 0x58, 0x6d, 0x74, 0x70, 0x5c, 0x58, 0x6d, 0x74,
	0x70, 0x76, 0x34, 0x5c, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x73, 0xe2, 0x02, 0x21,
	0x58, 0x6d, 0x74, 0x70, 0x5c, 0x58, 0x6d, 0x74, 0x70, 0x76, 0x34, 0x5c, 0x45, 0x6e, 0x76, 0x65,
	0x6c, 0x6f, 0x70, 0x65, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x17, 0x58, 0x6d, 0x74, 0x70, 0x3a, 0x3a, 0x58, 0x6d, 0x74, 0x70, 0x76, 0x34,
	0x3a, 0x3a, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
})

var (
	file_xmtpv4_envelopes_envelopes_proto_rawDescOnce sync.Once
	file_xmtpv4_envelopes_envelopes_proto_rawDescData []byte
)

func file_xmtpv4_envelopes_envelopes_proto_rawDescGZIP() []byte {
	file_xmtpv4_envelopes_envelopes_proto_rawDescOnce.Do(func() {
		file_xmtpv4_envelopes_envelopes_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_xmtpv4_envelopes_envelopes_proto_rawDesc), len(file_xmtpv4_envelopes_envelopes_proto_rawDesc)))
	})
	return file_xmtpv4_envelopes_envelopes_proto_rawDescData
}

var file_xmtpv4_envelopes_envelopes_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_xmtpv4_envelopes_envelopes_proto_goTypes = []any{
	(*Cursor)(nil),                                 // 0: xmtp.xmtpv4.envelopes.Cursor
	(*AuthenticatedData)(nil),                      // 1: xmtp.xmtpv4.envelopes.AuthenticatedData
	(*ClientEnvelope)(nil),                         // 2: xmtp.xmtpv4.envelopes.ClientEnvelope
	(*PayerEnvelope)(nil),                          // 3: xmtp.xmtpv4.envelopes.PayerEnvelope
	(*UnsignedOriginatorEnvelope)(nil),             // 4: xmtp.xmtpv4.envelopes.UnsignedOriginatorEnvelope
	(*BlockchainProof)(nil),                        // 5: xmtp.xmtpv4.envelopes.BlockchainProof
	(*OriginatorEnvelope)(nil),                     // 6: xmtp.xmtpv4.envelopes.OriginatorEnvelope
	nil,                                            // 7: xmtp.xmtpv4.envelopes.Cursor.NodeIdToSequenceIdEntry
	(*v1.GroupMessageInput)(nil),                   // 8: xmtp.mls.api.v1.GroupMessageInput
	(*v1.WelcomeMessageInput)(nil),                 // 9: xmtp.mls.api.v1.WelcomeMessageInput
	(*v1.UploadKeyPackageRequest)(nil),             // 10: xmtp.mls.api.v1.UploadKeyPackageRequest
	(*associations.IdentityUpdate)(nil),            // 11: xmtp.identity.associations.IdentityUpdate
	(*associations.RecoverableEcdsaSignature)(nil), // 12: xmtp.identity.associations.RecoverableEcdsaSignature
}
var file_xmtpv4_envelopes_envelopes_proto_depIdxs = []int32{
	7,  // 0: xmtp.xmtpv4.envelopes.Cursor.node_id_to_sequence_id:type_name -> xmtp.xmtpv4.envelopes.Cursor.NodeIdToSequenceIdEntry
	0,  // 1: xmtp.xmtpv4.envelopes.AuthenticatedData.depends_on:type_name -> xmtp.xmtpv4.envelopes.Cursor
	1,  // 2: xmtp.xmtpv4.envelopes.ClientEnvelope.aad:type_name -> xmtp.xmtpv4.envelopes.AuthenticatedData
	8,  // 3: xmtp.xmtpv4.envelopes.ClientEnvelope.group_message:type_name -> xmtp.mls.api.v1.GroupMessageInput
	9,  // 4: xmtp.xmtpv4.envelopes.ClientEnvelope.welcome_message:type_name -> xmtp.mls.api.v1.WelcomeMessageInput
	10, // 5: xmtp.xmtpv4.envelopes.ClientEnvelope.upload_key_package:type_name -> xmtp.mls.api.v1.UploadKeyPackageRequest
	11, // 6: xmtp.xmtpv4.envelopes.ClientEnvelope.identity_update:type_name -> xmtp.identity.associations.IdentityUpdate
	12, // 7: xmtp.xmtpv4.envelopes.PayerEnvelope.payer_signature:type_name -> xmtp.identity.associations.RecoverableEcdsaSignature
	3,  // 8: xmtp.xmtpv4.envelopes.UnsignedOriginatorEnvelope.payer_envelope:type_name -> xmtp.xmtpv4.envelopes.PayerEnvelope
	12, // 9: xmtp.xmtpv4.envelopes.OriginatorEnvelope.originator_signature:type_name -> xmtp.identity.associations.RecoverableEcdsaSignature
	5,  // 10: xmtp.xmtpv4.envelopes.OriginatorEnvelope.blockchain_proof:type_name -> xmtp.xmtpv4.envelopes.BlockchainProof
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_xmtpv4_envelopes_envelopes_proto_init() }
func file_xmtpv4_envelopes_envelopes_proto_init() {
	if File_xmtpv4_envelopes_envelopes_proto != nil {
		return
	}
	file_xmtpv4_envelopes_envelopes_proto_msgTypes[1].OneofWrappers = []any{}
	file_xmtpv4_envelopes_envelopes_proto_msgTypes[2].OneofWrappers = []any{
		(*ClientEnvelope_GroupMessage)(nil),
		(*ClientEnvelope_WelcomeMessage)(nil),
		(*ClientEnvelope_UploadKeyPackage)(nil),
		(*ClientEnvelope_IdentityUpdate)(nil),
	}
	file_xmtpv4_envelopes_envelopes_proto_msgTypes[6].OneofWrappers = []any{
		(*OriginatorEnvelope_OriginatorSignature)(nil),
		(*OriginatorEnvelope_BlockchainProof)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_xmtpv4_envelopes_envelopes_proto_rawDesc), len(file_xmtpv4_envelopes_envelopes_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_xmtpv4_envelopes_envelopes_proto_goTypes,
		DependencyIndexes: file_xmtpv4_envelopes_envelopes_proto_depIdxs,
		MessageInfos:      file_xmtpv4_envelopes_envelopes_proto_msgTypes,
	}.Build()
	File_xmtpv4_envelopes_envelopes_proto = out.File
	file_xmtpv4_envelopes_envelopes_proto_goTypes = nil
	file_xmtpv4_envelopes_envelopes_proto_depIdxs = nil
}
