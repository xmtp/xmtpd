package utils

import (
	"context"
	"net"
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
