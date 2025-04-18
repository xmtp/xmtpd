// Invitation is used by an initiator to invite participants
// into a new conversation. Invitation carries the chosen topic name
// and encryption scheme and key material to be used for message encryption.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: message_contents/invitation.proto

package message_contents

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

// Version of consent proof payload
type ConsentProofPayloadVersion int32

const (
	ConsentProofPayloadVersion_CONSENT_PROOF_PAYLOAD_VERSION_UNSPECIFIED ConsentProofPayloadVersion = 0
	ConsentProofPayloadVersion_CONSENT_PROOF_PAYLOAD_VERSION_1           ConsentProofPayloadVersion = 1
)

// Enum value maps for ConsentProofPayloadVersion.
var (
	ConsentProofPayloadVersion_name = map[int32]string{
		0: "CONSENT_PROOF_PAYLOAD_VERSION_UNSPECIFIED",
		1: "CONSENT_PROOF_PAYLOAD_VERSION_1",
	}
	ConsentProofPayloadVersion_value = map[string]int32{
		"CONSENT_PROOF_PAYLOAD_VERSION_UNSPECIFIED": 0,
		"CONSENT_PROOF_PAYLOAD_VERSION_1":           1,
	}
)

func (x ConsentProofPayloadVersion) Enum() *ConsentProofPayloadVersion {
	p := new(ConsentProofPayloadVersion)
	*p = x
	return p
}

func (x ConsentProofPayloadVersion) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConsentProofPayloadVersion) Descriptor() protoreflect.EnumDescriptor {
	return file_message_contents_invitation_proto_enumTypes[0].Descriptor()
}

func (ConsentProofPayloadVersion) Type() protoreflect.EnumType {
	return &file_message_contents_invitation_proto_enumTypes[0]
}

func (x ConsentProofPayloadVersion) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConsentProofPayloadVersion.Descriptor instead.
func (ConsentProofPayloadVersion) EnumDescriptor() ([]byte, []int) {
	return file_message_contents_invitation_proto_rawDescGZIP(), []int{0}
}

// Unsealed invitation V1
type InvitationV1 struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// topic name chosen for this conversation.
	// It MUST be randomly generated bytes (length >= 32),
	// then base64 encoded without padding
	Topic string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	// A context object defining metadata
	Context *InvitationV1_Context `protobuf:"bytes,2,opt,name=context,proto3" json:"context,omitempty"`
	// message encryption scheme and keys for this conversation.
	//
	// Types that are valid to be assigned to Encryption:
	//
	//	*InvitationV1_Aes256GcmHkdfSha256
	Encryption isInvitationV1_Encryption `protobuf_oneof:"encryption"`
	// The user's consent proof
	ConsentProof  *ConsentProofPayload `protobuf:"bytes,4,opt,name=consent_proof,json=consentProof,proto3" json:"consent_proof,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InvitationV1) Reset() {
	*x = InvitationV1{}
	mi := &file_message_contents_invitation_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InvitationV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvitationV1) ProtoMessage() {}

func (x *InvitationV1) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_invitation_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvitationV1.ProtoReflect.Descriptor instead.
func (*InvitationV1) Descriptor() ([]byte, []int) {
	return file_message_contents_invitation_proto_rawDescGZIP(), []int{0}
}

func (x *InvitationV1) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *InvitationV1) GetContext() *InvitationV1_Context {
	if x != nil {
		return x.Context
	}
	return nil
}

func (x *InvitationV1) GetEncryption() isInvitationV1_Encryption {
	if x != nil {
		return x.Encryption
	}
	return nil
}

func (x *InvitationV1) GetAes256GcmHkdfSha256() *InvitationV1_Aes256GcmHkdfsha256 {
	if x != nil {
		if x, ok := x.Encryption.(*InvitationV1_Aes256GcmHkdfSha256); ok {
			return x.Aes256GcmHkdfSha256
		}
	}
	return nil
}

func (x *InvitationV1) GetConsentProof() *ConsentProofPayload {
	if x != nil {
		return x.ConsentProof
	}
	return nil
}

type isInvitationV1_Encryption interface {
	isInvitationV1_Encryption()
}

type InvitationV1_Aes256GcmHkdfSha256 struct {
	// Specify the encryption method to process the key material properly.
	Aes256GcmHkdfSha256 *InvitationV1_Aes256GcmHkdfsha256 `protobuf:"bytes,3,opt,name=aes256_gcm_hkdf_sha256,json=aes256GcmHkdfSha256,proto3,oneof"`
}

func (*InvitationV1_Aes256GcmHkdfSha256) isInvitationV1_Encryption() {}

// Sealed Invitation V1 Header
// Header carries information that is unencrypted, thus readable by the network
// it is however authenticated as associated data with the AEAD scheme used
// to encrypt the invitation body, thus providing tamper evidence.
type SealedInvitationHeaderV1 struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Sender        *SignedPublicKeyBundle `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Recipient     *SignedPublicKeyBundle `protobuf:"bytes,2,opt,name=recipient,proto3" json:"recipient,omitempty"`
	CreatedNs     uint64                 `protobuf:"varint,3,opt,name=created_ns,json=createdNs,proto3" json:"created_ns,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SealedInvitationHeaderV1) Reset() {
	*x = SealedInvitationHeaderV1{}
	mi := &file_message_contents_invitation_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SealedInvitationHeaderV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SealedInvitationHeaderV1) ProtoMessage() {}

func (x *SealedInvitationHeaderV1) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_invitation_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SealedInvitationHeaderV1.ProtoReflect.Descriptor instead.
func (*SealedInvitationHeaderV1) Descriptor() ([]byte, []int) {
	return file_message_contents_invitation_proto_rawDescGZIP(), []int{1}
}

func (x *SealedInvitationHeaderV1) GetSender() *SignedPublicKeyBundle {
	if x != nil {
		return x.Sender
	}
	return nil
}

func (x *SealedInvitationHeaderV1) GetRecipient() *SignedPublicKeyBundle {
	if x != nil {
		return x.Recipient
	}
	return nil
}

func (x *SealedInvitationHeaderV1) GetCreatedNs() uint64 {
	if x != nil {
		return x.CreatedNs
	}
	return 0
}

// Sealed Invitation V1
// Invitation encrypted with key material derived from the sender's and
// recipient's public key bundles using simplified X3DH where
// the sender's ephemeral key is replaced with sender's pre-key.
type SealedInvitationV1 struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// encoded SealedInvitationHeaderV1 used as associated data for Ciphertext
	HeaderBytes []byte `protobuf:"bytes,1,opt,name=header_bytes,json=headerBytes,proto3" json:"header_bytes,omitempty"`
	// Ciphertext.payload MUST contain encrypted InvitationV1.
	Ciphertext    *Ciphertext `protobuf:"bytes,2,opt,name=ciphertext,proto3" json:"ciphertext,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SealedInvitationV1) Reset() {
	*x = SealedInvitationV1{}
	mi := &file_message_contents_invitation_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SealedInvitationV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SealedInvitationV1) ProtoMessage() {}

func (x *SealedInvitationV1) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_invitation_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SealedInvitationV1.ProtoReflect.Descriptor instead.
func (*SealedInvitationV1) Descriptor() ([]byte, []int) {
	return file_message_contents_invitation_proto_rawDescGZIP(), []int{2}
}

func (x *SealedInvitationV1) GetHeaderBytes() []byte {
	if x != nil {
		return x.HeaderBytes
	}
	return nil
}

func (x *SealedInvitationV1) GetCiphertext() *Ciphertext {
	if x != nil {
		return x.Ciphertext
	}
	return nil
}

// Versioned Sealed Invitation
type SealedInvitation struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Version:
	//
	//	*SealedInvitation_V1
	Version       isSealedInvitation_Version `protobuf_oneof:"version"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SealedInvitation) Reset() {
	*x = SealedInvitation{}
	mi := &file_message_contents_invitation_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SealedInvitation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SealedInvitation) ProtoMessage() {}

func (x *SealedInvitation) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_invitation_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SealedInvitation.ProtoReflect.Descriptor instead.
func (*SealedInvitation) Descriptor() ([]byte, []int) {
	return file_message_contents_invitation_proto_rawDescGZIP(), []int{3}
}

func (x *SealedInvitation) GetVersion() isSealedInvitation_Version {
	if x != nil {
		return x.Version
	}
	return nil
}

func (x *SealedInvitation) GetV1() *SealedInvitationV1 {
	if x != nil {
		if x, ok := x.Version.(*SealedInvitation_V1); ok {
			return x.V1
		}
	}
	return nil
}

type isSealedInvitation_Version interface {
	isSealedInvitation_Version()
}

type SealedInvitation_V1 struct {
	V1 *SealedInvitationV1 `protobuf:"bytes,1,opt,name=v1,proto3,oneof"`
}

func (*SealedInvitation_V1) isSealedInvitation_Version() {}

// Payload for user's consent proof to be set in the invitation
// Signifying the conversation should be preapproved for the user on receipt
type ConsentProofPayload struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the user's signature in hex format
	Signature string `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	// approximate time when the user signed
	Timestamp uint64 `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// version of the payload
	PayloadVersion ConsentProofPayloadVersion `protobuf:"varint,3,opt,name=payload_version,json=payloadVersion,proto3,enum=xmtp.message_contents.ConsentProofPayloadVersion" json:"payload_version,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *ConsentProofPayload) Reset() {
	*x = ConsentProofPayload{}
	mi := &file_message_contents_invitation_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConsentProofPayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsentProofPayload) ProtoMessage() {}

func (x *ConsentProofPayload) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_invitation_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsentProofPayload.ProtoReflect.Descriptor instead.
func (*ConsentProofPayload) Descriptor() ([]byte, []int) {
	return file_message_contents_invitation_proto_rawDescGZIP(), []int{4}
}

func (x *ConsentProofPayload) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

func (x *ConsentProofPayload) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *ConsentProofPayload) GetPayloadVersion() ConsentProofPayloadVersion {
	if x != nil {
		return x.PayloadVersion
	}
	return ConsentProofPayloadVersion_CONSENT_PROOF_PAYLOAD_VERSION_UNSPECIFIED
}

// Supported encryption schemes
// AES256-GCM-HKDF-SHA256
type InvitationV1_Aes256GcmHkdfsha256 struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	KeyMaterial   []byte                 `protobuf:"bytes,1,opt,name=key_material,json=keyMaterial,proto3" json:"key_material,omitempty"` // randomly generated key material (32 bytes)
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InvitationV1_Aes256GcmHkdfsha256) Reset() {
	*x = InvitationV1_Aes256GcmHkdfsha256{}
	mi := &file_message_contents_invitation_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InvitationV1_Aes256GcmHkdfsha256) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvitationV1_Aes256GcmHkdfsha256) ProtoMessage() {}

func (x *InvitationV1_Aes256GcmHkdfsha256) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_invitation_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvitationV1_Aes256GcmHkdfsha256.ProtoReflect.Descriptor instead.
func (*InvitationV1_Aes256GcmHkdfsha256) Descriptor() ([]byte, []int) {
	return file_message_contents_invitation_proto_rawDescGZIP(), []int{0, 0}
}

func (x *InvitationV1_Aes256GcmHkdfsha256) GetKeyMaterial() []byte {
	if x != nil {
		return x.KeyMaterial
	}
	return nil
}

// The context type
type InvitationV1_Context struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Expected to be a URI (ie xmtp.org/convo1)
	ConversationId string `protobuf:"bytes,1,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	// Key value map of additional metadata that would be exposed to
	// application developers and could be used for filtering
	Metadata      map[string]string `protobuf:"bytes,2,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InvitationV1_Context) Reset() {
	*x = InvitationV1_Context{}
	mi := &file_message_contents_invitation_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InvitationV1_Context) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvitationV1_Context) ProtoMessage() {}

func (x *InvitationV1_Context) ProtoReflect() protoreflect.Message {
	mi := &file_message_contents_invitation_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvitationV1_Context.ProtoReflect.Descriptor instead.
func (*InvitationV1_Context) Descriptor() ([]byte, []int) {
	return file_message_contents_invitation_proto_rawDescGZIP(), []int{0, 1}
}

func (x *InvitationV1_Context) GetConversationId() string {
	if x != nil {
		return x.ConversationId
	}
	return ""
}

func (x *InvitationV1_Context) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

var File_message_contents_invitation_proto protoreflect.FileDescriptor

const file_message_contents_invitation_proto_rawDesc = "" +
	"\n" +
	"!message_contents/invitation.proto\x12\x15xmtp.message_contents\x1a!message_contents/ciphertext.proto\x1a!message_contents/public_key.proto\"\xbd\x04\n" +
	"\fInvitationV1\x12\x14\n" +
	"\x05topic\x18\x01 \x01(\tR\x05topic\x12E\n" +
	"\acontext\x18\x02 \x01(\v2+.xmtp.message_contents.InvitationV1.ContextR\acontext\x12n\n" +
	"\x16aes256_gcm_hkdf_sha256\x18\x03 \x01(\v27.xmtp.message_contents.InvitationV1.Aes256gcmHkdfsha256H\x00R\x13aes256GcmHkdfSha256\x12O\n" +
	"\rconsent_proof\x18\x04 \x01(\v2*.xmtp.message_contents.ConsentProofPayloadR\fconsentProof\x1a8\n" +
	"\x13Aes256gcmHkdfsha256\x12!\n" +
	"\fkey_material\x18\x01 \x01(\fR\vkeyMaterial\x1a\xc6\x01\n" +
	"\aContext\x12'\n" +
	"\x0fconversation_id\x18\x01 \x01(\tR\x0econversationId\x12U\n" +
	"\bmetadata\x18\x02 \x03(\v29.xmtp.message_contents.InvitationV1.Context.MetadataEntryR\bmetadata\x1a;\n" +
	"\rMetadataEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01B\f\n" +
	"\n" +
	"encryption\"\xcb\x01\n" +
	"\x18SealedInvitationHeaderV1\x12D\n" +
	"\x06sender\x18\x01 \x01(\v2,.xmtp.message_contents.SignedPublicKeyBundleR\x06sender\x12J\n" +
	"\trecipient\x18\x02 \x01(\v2,.xmtp.message_contents.SignedPublicKeyBundleR\trecipient\x12\x1d\n" +
	"\n" +
	"created_ns\x18\x03 \x01(\x04R\tcreatedNs\"z\n" +
	"\x12SealedInvitationV1\x12!\n" +
	"\fheader_bytes\x18\x01 \x01(\fR\vheaderBytes\x12A\n" +
	"\n" +
	"ciphertext\x18\x02 \x01(\v2!.xmtp.message_contents.CiphertextR\n" +
	"ciphertext\"`\n" +
	"\x10SealedInvitation\x12;\n" +
	"\x02v1\x18\x01 \x01(\v2).xmtp.message_contents.SealedInvitationV1H\x00R\x02v1B\t\n" +
	"\aversionJ\x04\b\x02\x10\x03\"\xad\x01\n" +
	"\x13ConsentProofPayload\x12\x1c\n" +
	"\tsignature\x18\x01 \x01(\tR\tsignature\x12\x1c\n" +
	"\ttimestamp\x18\x02 \x01(\x04R\ttimestamp\x12Z\n" +
	"\x0fpayload_version\x18\x03 \x01(\x0e21.xmtp.message_contents.ConsentProofPayloadVersionR\x0epayloadVersion*p\n" +
	"\x1aConsentProofPayloadVersion\x12-\n" +
	")CONSENT_PROOF_PAYLOAD_VERSION_UNSPECIFIED\x10\x00\x12#\n" +
	"\x1fCONSENT_PROOF_PAYLOAD_VERSION_1\x10\x01B\xcf\x01\n" +
	"\x19com.xmtp.message_contentsB\x0fInvitationProtoP\x01Z0github.com/xmtp/xmtpd/pkg/proto/message_contents\xa2\x02\x03XMX\xaa\x02\x14Xmtp.MessageContents\xca\x02\x14Xmtp\\MessageContents\xe2\x02 Xmtp\\MessageContents\\GPBMetadata\xea\x02\x15Xmtp::MessageContentsb\x06proto3"

var (
	file_message_contents_invitation_proto_rawDescOnce sync.Once
	file_message_contents_invitation_proto_rawDescData []byte
)

func file_message_contents_invitation_proto_rawDescGZIP() []byte {
	file_message_contents_invitation_proto_rawDescOnce.Do(func() {
		file_message_contents_invitation_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_message_contents_invitation_proto_rawDesc), len(file_message_contents_invitation_proto_rawDesc)))
	})
	return file_message_contents_invitation_proto_rawDescData
}

var file_message_contents_invitation_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_message_contents_invitation_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_message_contents_invitation_proto_goTypes = []any{
	(ConsentProofPayloadVersion)(0),          // 0: xmtp.message_contents.ConsentProofPayloadVersion
	(*InvitationV1)(nil),                     // 1: xmtp.message_contents.InvitationV1
	(*SealedInvitationHeaderV1)(nil),         // 2: xmtp.message_contents.SealedInvitationHeaderV1
	(*SealedInvitationV1)(nil),               // 3: xmtp.message_contents.SealedInvitationV1
	(*SealedInvitation)(nil),                 // 4: xmtp.message_contents.SealedInvitation
	(*ConsentProofPayload)(nil),              // 5: xmtp.message_contents.ConsentProofPayload
	(*InvitationV1_Aes256GcmHkdfsha256)(nil), // 6: xmtp.message_contents.InvitationV1.Aes256gcmHkdfsha256
	(*InvitationV1_Context)(nil),             // 7: xmtp.message_contents.InvitationV1.Context
	nil,                                      // 8: xmtp.message_contents.InvitationV1.Context.MetadataEntry
	(*SignedPublicKeyBundle)(nil),            // 9: xmtp.message_contents.SignedPublicKeyBundle
	(*Ciphertext)(nil),                       // 10: xmtp.message_contents.Ciphertext
}
var file_message_contents_invitation_proto_depIdxs = []int32{
	7,  // 0: xmtp.message_contents.InvitationV1.context:type_name -> xmtp.message_contents.InvitationV1.Context
	6,  // 1: xmtp.message_contents.InvitationV1.aes256_gcm_hkdf_sha256:type_name -> xmtp.message_contents.InvitationV1.Aes256gcmHkdfsha256
	5,  // 2: xmtp.message_contents.InvitationV1.consent_proof:type_name -> xmtp.message_contents.ConsentProofPayload
	9,  // 3: xmtp.message_contents.SealedInvitationHeaderV1.sender:type_name -> xmtp.message_contents.SignedPublicKeyBundle
	9,  // 4: xmtp.message_contents.SealedInvitationHeaderV1.recipient:type_name -> xmtp.message_contents.SignedPublicKeyBundle
	10, // 5: xmtp.message_contents.SealedInvitationV1.ciphertext:type_name -> xmtp.message_contents.Ciphertext
	3,  // 6: xmtp.message_contents.SealedInvitation.v1:type_name -> xmtp.message_contents.SealedInvitationV1
	0,  // 7: xmtp.message_contents.ConsentProofPayload.payload_version:type_name -> xmtp.message_contents.ConsentProofPayloadVersion
	8,  // 8: xmtp.message_contents.InvitationV1.Context.metadata:type_name -> xmtp.message_contents.InvitationV1.Context.MetadataEntry
	9,  // [9:9] is the sub-list for method output_type
	9,  // [9:9] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_message_contents_invitation_proto_init() }
func file_message_contents_invitation_proto_init() {
	if File_message_contents_invitation_proto != nil {
		return
	}
	file_message_contents_ciphertext_proto_init()
	file_message_contents_public_key_proto_init()
	file_message_contents_invitation_proto_msgTypes[0].OneofWrappers = []any{
		(*InvitationV1_Aes256GcmHkdfSha256)(nil),
	}
	file_message_contents_invitation_proto_msgTypes[3].OneofWrappers = []any{
		(*SealedInvitation_V1)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_message_contents_invitation_proto_rawDesc), len(file_message_contents_invitation_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_contents_invitation_proto_goTypes,
		DependencyIndexes: file_message_contents_invitation_proto_depIdxs,
		EnumInfos:         file_message_contents_invitation_proto_enumTypes,
		MessageInfos:      file_message_contents_invitation_proto_msgTypes,
	}.Build()
	File_message_contents_invitation_proto = out.File
	file_message_contents_invitation_proto_goTypes = nil
	file_message_contents_invitation_proto_depIdxs = nil
}
