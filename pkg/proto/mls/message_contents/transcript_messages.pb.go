// Message content encoding structures

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: mls/message_contents/transcript_messages.proto

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

// A group member and affected installation IDs
type MembershipChange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InstallationIds           [][]byte `protobuf:"bytes,1,rep,name=installation_ids,json=installationIds,proto3" json:"installation_ids,omitempty"`
	AccountAddress            string   `protobuf:"bytes,2,opt,name=account_address,json=accountAddress,proto3" json:"account_address,omitempty"`
	InitiatedByAccountAddress string   `protobuf:"bytes,3,opt,name=initiated_by_account_address,json=initiatedByAccountAddress,proto3" json:"initiated_by_account_address,omitempty"`
}

func (x *MembershipChange) Reset() {
	*x = MembershipChange{}
	mi := &file_mls_message_contents_transcript_messages_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MembershipChange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MembershipChange) ProtoMessage() {}

func (x *MembershipChange) ProtoReflect() protoreflect.Message {
	mi := &file_mls_message_contents_transcript_messages_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MembershipChange.ProtoReflect.Descriptor instead.
func (*MembershipChange) Descriptor() ([]byte, []int) {
	return file_mls_message_contents_transcript_messages_proto_rawDescGZIP(), []int{0}
}

func (x *MembershipChange) GetInstallationIds() [][]byte {
	if x != nil {
		return x.InstallationIds
	}
	return nil
}

func (x *MembershipChange) GetAccountAddress() string {
	if x != nil {
		return x.AccountAddress
	}
	return ""
}

func (x *MembershipChange) GetInitiatedByAccountAddress() string {
	if x != nil {
		return x.InitiatedByAccountAddress
	}
	return ""
}

// The group membership change proto
type GroupMembershipChanges struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Members that have been added in the commit
	MembersAdded []*MembershipChange `protobuf:"bytes,1,rep,name=members_added,json=membersAdded,proto3" json:"members_added,omitempty"`
	// Members that have been removed in the commit
	MembersRemoved []*MembershipChange `protobuf:"bytes,2,rep,name=members_removed,json=membersRemoved,proto3" json:"members_removed,omitempty"`
	// Installations that have been added in the commit, grouped by member
	InstallationsAdded []*MembershipChange `protobuf:"bytes,3,rep,name=installations_added,json=installationsAdded,proto3" json:"installations_added,omitempty"`
	// Installations removed in the commit, grouped by member
	InstallationsRemoved []*MembershipChange `protobuf:"bytes,4,rep,name=installations_removed,json=installationsRemoved,proto3" json:"installations_removed,omitempty"`
}

func (x *GroupMembershipChanges) Reset() {
	*x = GroupMembershipChanges{}
	mi := &file_mls_message_contents_transcript_messages_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GroupMembershipChanges) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupMembershipChanges) ProtoMessage() {}

func (x *GroupMembershipChanges) ProtoReflect() protoreflect.Message {
	mi := &file_mls_message_contents_transcript_messages_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupMembershipChanges.ProtoReflect.Descriptor instead.
func (*GroupMembershipChanges) Descriptor() ([]byte, []int) {
	return file_mls_message_contents_transcript_messages_proto_rawDescGZIP(), []int{1}
}

func (x *GroupMembershipChanges) GetMembersAdded() []*MembershipChange {
	if x != nil {
		return x.MembersAdded
	}
	return nil
}

func (x *GroupMembershipChanges) GetMembersRemoved() []*MembershipChange {
	if x != nil {
		return x.MembersRemoved
	}
	return nil
}

func (x *GroupMembershipChanges) GetInstallationsAdded() []*MembershipChange {
	if x != nil {
		return x.InstallationsAdded
	}
	return nil
}

func (x *GroupMembershipChanges) GetInstallationsRemoved() []*MembershipChange {
	if x != nil {
		return x.InstallationsRemoved
	}
	return nil
}

// A summary of the changes in a commit.
// Includes added/removed inboxes and changes to metadata
type GroupUpdated struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InitiatedByInboxId string `protobuf:"bytes,1,opt,name=initiated_by_inbox_id,json=initiatedByInboxId,proto3" json:"initiated_by_inbox_id,omitempty"`
	// The inboxes added in the commit
	AddedInboxes []*GroupUpdated_Inbox `protobuf:"bytes,2,rep,name=added_inboxes,json=addedInboxes,proto3" json:"added_inboxes,omitempty"`
	// The inboxes removed in the commit
	RemovedInboxes []*GroupUpdated_Inbox `protobuf:"bytes,3,rep,name=removed_inboxes,json=removedInboxes,proto3" json:"removed_inboxes,omitempty"`
	// The metadata changes in the commit
	MetadataFieldChanges []*GroupUpdated_MetadataFieldChange `protobuf:"bytes,4,rep,name=metadata_field_changes,json=metadataFieldChanges,proto3" json:"metadata_field_changes,omitempty"`
}

func (x *GroupUpdated) Reset() {
	*x = GroupUpdated{}
	mi := &file_mls_message_contents_transcript_messages_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GroupUpdated) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupUpdated) ProtoMessage() {}

func (x *GroupUpdated) ProtoReflect() protoreflect.Message {
	mi := &file_mls_message_contents_transcript_messages_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupUpdated.ProtoReflect.Descriptor instead.
func (*GroupUpdated) Descriptor() ([]byte, []int) {
	return file_mls_message_contents_transcript_messages_proto_rawDescGZIP(), []int{2}
}

func (x *GroupUpdated) GetInitiatedByInboxId() string {
	if x != nil {
		return x.InitiatedByInboxId
	}
	return ""
}

func (x *GroupUpdated) GetAddedInboxes() []*GroupUpdated_Inbox {
	if x != nil {
		return x.AddedInboxes
	}
	return nil
}

func (x *GroupUpdated) GetRemovedInboxes() []*GroupUpdated_Inbox {
	if x != nil {
		return x.RemovedInboxes
	}
	return nil
}

func (x *GroupUpdated) GetMetadataFieldChanges() []*GroupUpdated_MetadataFieldChange {
	if x != nil {
		return x.MetadataFieldChanges
	}
	return nil
}

// An inbox that was added or removed in this commit
type GroupUpdated_Inbox struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InboxId string `protobuf:"bytes,1,opt,name=inbox_id,json=inboxId,proto3" json:"inbox_id,omitempty"`
}

func (x *GroupUpdated_Inbox) Reset() {
	*x = GroupUpdated_Inbox{}
	mi := &file_mls_message_contents_transcript_messages_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GroupUpdated_Inbox) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupUpdated_Inbox) ProtoMessage() {}

func (x *GroupUpdated_Inbox) ProtoReflect() protoreflect.Message {
	mi := &file_mls_message_contents_transcript_messages_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupUpdated_Inbox.ProtoReflect.Descriptor instead.
func (*GroupUpdated_Inbox) Descriptor() ([]byte, []int) {
	return file_mls_message_contents_transcript_messages_proto_rawDescGZIP(), []int{2, 0}
}

func (x *GroupUpdated_Inbox) GetInboxId() string {
	if x != nil {
		return x.InboxId
	}
	return ""
}

// A summary of a change to the mutable metadata
type GroupUpdated_MetadataFieldChange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The field that was changed
	FieldName string `protobuf:"bytes,1,opt,name=field_name,json=fieldName,proto3" json:"field_name,omitempty"`
	// The previous value
	OldValue *string `protobuf:"bytes,2,opt,name=old_value,json=oldValue,proto3,oneof" json:"old_value,omitempty"`
	// The updated value
	NewValue *string `protobuf:"bytes,3,opt,name=new_value,json=newValue,proto3,oneof" json:"new_value,omitempty"`
}

func (x *GroupUpdated_MetadataFieldChange) Reset() {
	*x = GroupUpdated_MetadataFieldChange{}
	mi := &file_mls_message_contents_transcript_messages_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GroupUpdated_MetadataFieldChange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupUpdated_MetadataFieldChange) ProtoMessage() {}

func (x *GroupUpdated_MetadataFieldChange) ProtoReflect() protoreflect.Message {
	mi := &file_mls_message_contents_transcript_messages_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupUpdated_MetadataFieldChange.ProtoReflect.Descriptor instead.
func (*GroupUpdated_MetadataFieldChange) Descriptor() ([]byte, []int) {
	return file_mls_message_contents_transcript_messages_proto_rawDescGZIP(), []int{2, 1}
}

func (x *GroupUpdated_MetadataFieldChange) GetFieldName() string {
	if x != nil {
		return x.FieldName
	}
	return ""
}

func (x *GroupUpdated_MetadataFieldChange) GetOldValue() string {
	if x != nil && x.OldValue != nil {
		return *x.OldValue
	}
	return ""
}

func (x *GroupUpdated_MetadataFieldChange) GetNewValue() string {
	if x != nil && x.NewValue != nil {
		return *x.NewValue
	}
	return ""
}

var File_mls_message_contents_transcript_messages_proto protoreflect.FileDescriptor

var file_mls_message_contents_transcript_messages_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x6d, 0x6c, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x19, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x6c, 0x73, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x22, 0xa7, 0x01, 0x0a, 0x10,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x12, 0x29, 0x0a, 0x10, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0f, 0x69, 0x6e, 0x73, 0x74,
	0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x3f, 0x0a, 0x1c, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x62, 0x79, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x19, 0x69, 0x6e, 0x69, 0x74,
	0x69, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x80, 0x03, 0x0a, 0x16, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73,
	0x12, 0x50, 0x0a, 0x0d, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x5f, 0x61, 0x64, 0x64, 0x65,
	0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d,
	0x6c, 0x73, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x52, 0x0c, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x41, 0x64, 0x64,
	0x65, 0x64, 0x12, 0x54, 0x0a, 0x0f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x5f, 0x72, 0x65,
	0x6d, 0x6f, 0x76, 0x65, 0x64, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x78, 0x6d,
	0x74, 0x70, 0x2e, 0x6d, 0x6c, 0x73, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68,
	0x69, 0x70, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x0e, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x73, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x12, 0x5c, 0x0a, 0x13, 0x69, 0x6e, 0x73, 0x74,
	0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x61, 0x64, 0x64, 0x65, 0x64, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x6c, 0x73,
	0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x52, 0x12, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x41, 0x64, 0x64, 0x65, 0x64, 0x12, 0x60, 0x0a, 0x15, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x6c, 0x73,
	0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x52, 0x14, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x22, 0x9b, 0x04, 0x0a, 0x0c, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x31, 0x0a, 0x15, 0x69, 0x6e, 0x69,
	0x74, 0x69, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x5f, 0x69, 0x6e, 0x62, 0x6f, 0x78, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61,
	0x74, 0x65, 0x64, 0x42, 0x79, 0x49, 0x6e, 0x62, 0x6f, 0x78, 0x49, 0x64, 0x12, 0x52, 0x0a, 0x0d,
	0x61, 0x64, 0x64, 0x65, 0x64, 0x5f, 0x69, 0x6e, 0x62, 0x6f, 0x78, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x6c, 0x73, 0x2e, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x2e,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x49, 0x6e, 0x62,
	0x6f, 0x78, 0x52, 0x0c, 0x61, 0x64, 0x64, 0x65, 0x64, 0x49, 0x6e, 0x62, 0x6f, 0x78, 0x65, 0x73,
	0x12, 0x56, 0x0a, 0x0f, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x5f, 0x69, 0x6e, 0x62, 0x6f,
	0x78, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x78, 0x6d, 0x74, 0x70,
	0x2e, 0x6d, 0x6c, 0x73, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x2e, 0x49, 0x6e, 0x62, 0x6f, 0x78, 0x52, 0x0e, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x64, 0x49, 0x6e, 0x62, 0x6f, 0x78, 0x65, 0x73, 0x12, 0x71, 0x0a, 0x16, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3b, 0x2e, 0x78, 0x6d, 0x74, 0x70, 0x2e,
	0x6d, 0x6c, 0x73, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x73, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x14, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x1a, 0x22, 0x0a, 0x05, 0x49,
	0x6e, 0x62, 0x6f, 0x78, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x6e, 0x62, 0x6f, 0x78, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x69, 0x6e, 0x62, 0x6f, 0x78, 0x49, 0x64, 0x1a,
	0x94, 0x01, 0x0a, 0x13, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x46, 0x69, 0x65, 0x6c,
	0x64, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x09, 0x6f, 0x6c, 0x64, 0x5f, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x6f, 0x6c, 0x64,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x09, 0x6e, 0x65, 0x77, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x6e,
	0x65, 0x77, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6f,
	0x6c, 0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6e, 0x65, 0x77,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0xf0, 0x01, 0x0a, 0x1d, 0x63, 0x6f, 0x6d, 0x2e, 0x78,
	0x6d, 0x74, 0x70, 0x2e, 0x6d, 0x6c, 0x73, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x42, 0x17, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x78, 0x6d, 0x74, 0x70, 0x2f, 0x78, 0x6d, 0x74, 0x70, 0x64, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6c, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0xa2, 0x02, 0x03, 0x58, 0x4d, 0x4d, 0xaa,
	0x02, 0x18, 0x58, 0x6d, 0x74, 0x70, 0x2e, 0x4d, 0x6c, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0xca, 0x02, 0x18, 0x58, 0x6d, 0x74,
	0x70, 0x5c, 0x4d, 0x6c, 0x73, 0x5c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x73, 0xe2, 0x02, 0x24, 0x58, 0x6d, 0x74, 0x70, 0x5c, 0x4d, 0x6c, 0x73,
	0x5c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1a, 0x58,
	0x6d, 0x74, 0x70, 0x3a, 0x3a, 0x4d, 0x6c, 0x73, 0x3a, 0x3a, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_mls_message_contents_transcript_messages_proto_rawDescOnce sync.Once
	file_mls_message_contents_transcript_messages_proto_rawDescData = file_mls_message_contents_transcript_messages_proto_rawDesc
)

func file_mls_message_contents_transcript_messages_proto_rawDescGZIP() []byte {
	file_mls_message_contents_transcript_messages_proto_rawDescOnce.Do(func() {
		file_mls_message_contents_transcript_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_mls_message_contents_transcript_messages_proto_rawDescData)
	})
	return file_mls_message_contents_transcript_messages_proto_rawDescData
}

var file_mls_message_contents_transcript_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_mls_message_contents_transcript_messages_proto_goTypes = []any{
	(*MembershipChange)(nil),                 // 0: xmtp.mls.message_contents.MembershipChange
	(*GroupMembershipChanges)(nil),           // 1: xmtp.mls.message_contents.GroupMembershipChanges
	(*GroupUpdated)(nil),                     // 2: xmtp.mls.message_contents.GroupUpdated
	(*GroupUpdated_Inbox)(nil),               // 3: xmtp.mls.message_contents.GroupUpdated.Inbox
	(*GroupUpdated_MetadataFieldChange)(nil), // 4: xmtp.mls.message_contents.GroupUpdated.MetadataFieldChange
}
var file_mls_message_contents_transcript_messages_proto_depIdxs = []int32{
	0, // 0: xmtp.mls.message_contents.GroupMembershipChanges.members_added:type_name -> xmtp.mls.message_contents.MembershipChange
	0, // 1: xmtp.mls.message_contents.GroupMembershipChanges.members_removed:type_name -> xmtp.mls.message_contents.MembershipChange
	0, // 2: xmtp.mls.message_contents.GroupMembershipChanges.installations_added:type_name -> xmtp.mls.message_contents.MembershipChange
	0, // 3: xmtp.mls.message_contents.GroupMembershipChanges.installations_removed:type_name -> xmtp.mls.message_contents.MembershipChange
	3, // 4: xmtp.mls.message_contents.GroupUpdated.added_inboxes:type_name -> xmtp.mls.message_contents.GroupUpdated.Inbox
	3, // 5: xmtp.mls.message_contents.GroupUpdated.removed_inboxes:type_name -> xmtp.mls.message_contents.GroupUpdated.Inbox
	4, // 6: xmtp.mls.message_contents.GroupUpdated.metadata_field_changes:type_name -> xmtp.mls.message_contents.GroupUpdated.MetadataFieldChange
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_mls_message_contents_transcript_messages_proto_init() }
func file_mls_message_contents_transcript_messages_proto_init() {
	if File_mls_message_contents_transcript_messages_proto != nil {
		return
	}
	file_mls_message_contents_transcript_messages_proto_msgTypes[4].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mls_message_contents_transcript_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_mls_message_contents_transcript_messages_proto_goTypes,
		DependencyIndexes: file_mls_message_contents_transcript_messages_proto_depIdxs,
		MessageInfos:      file_mls_message_contents_transcript_messages_proto_msgTypes,
	}.Build()
	File_mls_message_contents_transcript_messages_proto = out.File
	file_mls_message_contents_transcript_messages_proto_rawDesc = nil
	file_mls_message_contents_transcript_messages_proto_goTypes = nil
	file_mls_message_contents_transcript_messages_proto_depIdxs = nil
}
