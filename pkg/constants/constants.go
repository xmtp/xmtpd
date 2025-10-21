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
)

type VerifiedNodeRequestCtxKey struct{}
