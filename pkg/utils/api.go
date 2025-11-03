// Package utils implements different utilities to be used by production code.
package utils

import (
	"context"
	"net"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func ClientIPFromContext(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	vals := md.Get("x-forwarded-for")
	if len(vals) == 0 {
		p, ok := peer.FromContext(ctx)
		if ok {
			host, _, err := net.SplitHostPort(p.Addr.String())
			if err != nil {
				// If SplitHostPort fails, fall back to the original string splitting
				ipAndPort := strings.Split(p.Addr.String(), ":")
				return ipAndPort[0]
			}
			return host
		} else {
			return ""
		}
	}
	// There are potentially multiple comma separated IPs bundled in that first value
	ips := strings.Split(vals[0], ",")
	return strings.TrimSpace(ips[0])
}

func ClientIPFromHeaderOrPeer(headers http.Header, peer string) string {
	xForwardedFor := headers.Get("x-forwarded-for")

	// Try peer string if x-forwarded-for is not set.
	if len(xForwardedFor) == 0 {
		if peer != "" {
			host, _, err := net.SplitHostPort(peer)
			if err != nil {
				return ""
			}
			return host
		}

		return ""
	}

	ips := strings.Split(xForwardedFor, ",")

	return strings.TrimSpace(ips[0])
}
