package constants

const (
	JWT_DOMAIN_SEPARATION_LABEL        = "jwt|"
	PAYER_DOMAIN_SEPARATION_LABEL      = "payer|"
	TARGET_ORIGINATOR_SEPARATION_LABEL = "target|"
	ORIGINATOR_DOMAIN_SEPARATION_LABEL = "originator|"
	NODE_AUTHORIZATION_HEADER_NAME     = "node-authorization"
	MAX_BLOCKCHAIN_ORIGINATOR_ID       = 100
)

type VerifiedNodeRequestCtxKey struct{}
