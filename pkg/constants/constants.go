// Package constants implements the constants for the xmtpd server.
package constants

const (
	JWTDomainSeparationLabel              = "jwt|"
	PayerDomainSeparationLabel            = "payer|"
	TargetOriginatorDomainSeparationLabel = "target|"
	OriginatorDomainSeparationLabel       = "originator|"
	NodeAuthorizationHeaderName           = "node-authorization"
	DefaultStorageDurationDays            = 60

	// Indexer originator IDs for group messages and identity updates.
	GroupMessageOriginatorID   = 0
	IdentityUpdateOriginatorID = 1

	// TODO: Revert to 25 * 1024 * 1024 after node sync is migrated to use the new API with pagination.
	// Has to be in sync with xmtp/libxmtp/crates/xmtp_configuration/src/common/api.rs GRPC_PAYLOAD_LIMIT.
	GRPCPayloadLimit = 50 * 1024 * 1024
)

type VerifiedNodeRequestCtxKey struct{}
