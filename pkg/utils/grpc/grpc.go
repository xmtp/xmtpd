// Package grpc provides utility functions for working with gRPC.
package grpc

import "strings"

// ParseProcedure extracts full service name, service name and method from a gRPC path.
// The procedure path format is: /package.service/method
// Input: /xmtp.xmtpv4.message_api.ReplicationApi/QueryEnvelopes
// Output: fullServiceName="xmtp.xmtpv4.message_api.ReplicationApi", serviceName="ReplicationApi", method="QueryEnvelopes"
func ParseProcedure(path string) (fullServiceName string, service string, method string) {
	// Trim leading slash without allocating.
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	normalizedPath := strings.TrimSuffix(path, "/")

	// Find last slash and return the service and method.
	if i := strings.LastIndexByte(normalizedPath, '/'); i != -1 {
		serviceParts := strings.Split(normalizedPath[:i], ".")
		if len(serviceParts) > 0 && serviceParts[len(serviceParts)-1] != "" {
			return normalizedPath[:i], serviceParts[len(serviceParts)-1], normalizedPath[i+1:]
		}
	}

	return "unknown", "unknown", normalizedPath
}
