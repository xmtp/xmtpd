package ratelimiter

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/xmtp/xmtpd/pkg/constants"
)

// Tier identifies which rate-limit policy applies to a request.
type Tier int

const (
	// Tier0 is authenticated node-to-node traffic. Bypasses all limits.
	Tier0 Tier = iota
	// Tier2 is unauthenticated edge-client traffic. Subject to limits.
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

// ParseTrustedProxyCIDRs parses a comma-separated CIDR list. Empty entries are
// ignored. Returns an error on the first malformed entry.
func ParseTrustedProxyCIDRs(s string) ([]*net.IPNet, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	out := make([]*net.IPNet, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		_, n, err := net.ParseCIDR(p)
		if err != nil {
			return nil, fmt.Errorf("invalid trusted-proxy CIDR %q: %w", p, err)
		}
		out = append(out, n)
	}
	return out, nil
}

// ExtractClientIP returns the bucket-key IP for an incoming request.
//
// The algorithm:
//  1. Parse the immediate peer address (host:port).
//  2. While the peer is in any trusted CIDR and X-Forwarded-For has remaining
//     entries, peel the rightmost XFF entry and treat it as the new peer.
//  3. For IPv4, return the dotted-quad string. For IPv6, return the /64 prefix
//     in CIDR notation (clients within a /64 share a bucket).
func ExtractClientIP(peerAddr, xff string, trusted []*net.IPNet) string {
	host, _, err := net.SplitHostPort(peerAddr)
	if err != nil {
		host = peerAddr
	}

	xffParts := splitXFF(xff)
	for ipInTrusted(host, trusted) && len(xffParts) > 0 {
		host = strings.TrimSpace(xffParts[len(xffParts)-1])
		xffParts = xffParts[:len(xffParts)-1]
	}

	parsed := net.ParseIP(host)
	if parsed == nil {
		return host
	}
	if v4 := parsed.To4(); v4 != nil {
		return v4.String()
	}
	mask := net.CIDRMask(64, 128)
	return (&net.IPNet{IP: parsed.Mask(mask), Mask: mask}).String()
}

func splitXFF(xff string) []string {
	if xff == "" {
		return nil
	}
	return strings.Split(xff, ",")
}

func ipInTrusted(host string, trusted []*net.IPNet) bool {
	if len(trusted) == 0 {
		return false
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	for _, n := range trusted {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}
