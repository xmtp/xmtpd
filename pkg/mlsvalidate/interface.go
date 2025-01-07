package mlsvalidate

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	identity_proto "github.com/xmtp/xmtpd/pkg/proto/identity"
	associations "github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	mlsv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
)

type KeyPackageValidationResult struct {
	IsOk            bool
	InstallationKey []byte
	Credential      *identity_proto.MlsCredential
	Expiration      uint64
}

type GroupMessageValidationResult struct {
	GroupId string
}

type AssociationStateResult struct {
	AssociationState *associations.AssociationState     `protobuf:"bytes,1,opt,name=association_state,json=associationState,proto3" json:"association_state,omitempty"`
	StateDiff        *associations.AssociationStateDiff `protobuf:"bytes,2,opt,name=state_diff,json=stateDiff,proto3"               json:"state_diff,omitempty"`
}

type MLSValidationService interface {
	ValidateKeyPackages(
		ctx context.Context,
		keyPackages [][]byte,
	) ([]KeyPackageValidationResult, error)
	ValidateGroupMessages(
		ctx context.Context,
		groupMessages []*mlsv1.GroupMessageInput,
	) ([]GroupMessageValidationResult, error)
	GetAssociationState(
		ctx context.Context,
		oldUpdates []*associations.IdentityUpdate,
		newUpdates []*associations.IdentityUpdate,
	) (*AssociationStateResult, error)
	GetAssociationStateFromEnvelopes(
		ctx context.Context,
		oldUpdates []queries.GatewayEnvelope,
		newIdentityUpdate *associations.IdentityUpdate,
	) (*AssociationStateResult, error)
}
