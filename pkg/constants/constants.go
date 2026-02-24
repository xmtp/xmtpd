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

	GRPCPayloadLimit = 1024 * 1024 * 25 // 25MB
)

type VerifiedNodeRequestCtxKey struct{}
