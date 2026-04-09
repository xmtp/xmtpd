package ratelimiter

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/constants"
)

// Tier identifies which rate-limit policy applies to a request.
//
// Note: Tier1 (authenticated fan-out / app-operator clients) is intentionally
// reserved for a future phase. v1 only distinguishes Tier 0 (node JWT bypass)
// and Tier 2 (everyone else, IP-keyed). Do not introduce Tier1 without
// updating the spec at tasks/rate-limiting-spec.md first.
type Tier int

const (
	// Tier0 is authenticated node-to-node traffic. Bypasses all limits.
	Tier0 Tier = iota
	// Tier2 is unauthenticated edge-client traffic. Subject to limits.
	// (Tier1 is reserved — see Tier doc comment.)
	Tier2
)

// ClassifyTier inspects the request context for the auth interceptor's
// verified-node flag. If present and true, the request is Tier 0. Otherwise
// it is Tier 2. Tier 0 verification (JWT validation, signer-in-registry check)
// is performed by the auth interceptor upstream — the classifier never
// re-verifies and never falls back to Tier 2 on JWT failure (the auth
// interceptor returns Unauthenticated directly in that case).
func ClassifyTier(ctx context.Context) Tier {
	if v, ok := ctx.Value(constants.VerifiedNodeRequestCtxKey{}).(bool); ok && v {
		return Tier0
	}
	return Tier2
}
