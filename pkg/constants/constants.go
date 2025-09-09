// Package constants implements the constants for the xmtpd server.
package constants

const (
	JWTDomainSeparationLabel              = "jwt|"
	PayerDomainSeparationLabel            = "payer|"
	TargetOriginatorDomainSeparationLabel = "target|"
	OriginatorDomainSeparationLabel       = "originator|"
	NodeAuthorizationHeaderName           = "node-authorization"
	DefaultStorageDurationDays            = 60
)

type VerifiedNodeRequestCtxKey struct{}
