// Package clientip extracts client IP addresses from incoming connections
// with optional trusted-proxy X-Forwarded-For peeling and IPv6 normalization.
// It is consumed by both the rate-limit interceptor and the SubscribeTopics
// handler, so it lives in a neutral utility package rather than in either
// caller.
package clientip

import (
	"fmt"
	"net"
	"strings"
)

// InvalidClientIPKey is the bucket key returned by ExtractClientIP when the
// client IP cannot be parsed as a valid IP address. All clients with malformed
// peer addresses or junk in X-Forwarded-For share this single bucket, which
// rate-limits the aggregate of all malformed traffic without leaking arbitrary
// strings into any downstream keyspace.
const InvalidClientIPKey = "invalid"

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

// Extract returns the bucket-key IP for an incoming request.
//
// The algorithm:
//  1. Parse the immediate peer address (host:port).
//  2. While the peer is in any trusted CIDR and X-Forwarded-For has remaining
//     entries, peel the rightmost XFF entry and treat it as the new peer.
//  3. For IPv4, return the dotted-quad string. For IPv6, return the /64 prefix
//     in CIDR notation (clients within a /64 share a bucket).
//  4. If the resulting host is not a parseable IP (junk peer addr or malformed
//     XFF entry), return InvalidClientIPKey rather than leaking arbitrary
//     strings downstream.
func Extract(peerAddr, xff string, trusted []*net.IPNet) string {
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
		return InvalidClientIPKey
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
