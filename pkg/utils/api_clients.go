package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api/payer_apiconnect"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// File api_clients.go implements the API clients for the different APIs.
// The replication API client is provided with all the different protocols supported by the replication API.
// For metadata and gateway APIs, only the base client is provided, and the consumer has to specify the api options.

const (
	maxMessageSize  = 25 * 1024 * 1024
	readIdleTimeout = 10 * time.Second
	pingTimeout     = 30 * time.Second
	clientTimeout   = 10 * time.Second
)

/* Connect-based clients */

// NewConnectReplicationAPIClient builds a Connect (default protocol) Replication API client.
//   - Uses connect-go default (Connect protocol) over HTTP/1.1 or HTTP/2.
//   - Requires a schemeful base URL (http:// or https://).
//   - For HTTP/2 (TLS) or h2c (plaintext), pass an http.Client configured appropriately
//     (e.g., via utils.BuildHTTP2Client).
func NewConnectReplicationAPIClient(
	ctx context.Context,
	httpAddress string,
	extraDialOpts ...connect.ClientOption,
) (message_apiconnect.ReplicationApiClient, error) {
	target, isTLS, err := HTTPAddressToConnectProtocolTarget(httpAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid http address: %w", err)
	}

	httpClient, err := BuildHTTP2Client(ctx, isTLS)
	if err != nil {
		return nil, fmt.Errorf("failed to build http client: %w", err)
	}

	opts := BuildConnectProtocolDialOptions(extraDialOpts...)

	return message_apiconnect.NewReplicationApiClient(httpClient, target, opts...), nil
}

// NewConnectGRPCReplicationAPIClient builds a Connect-based client configured to speak classic gRPC.
// - Uses connect.WithGRPC() (wire-compatible gRPC over an http.Client).
// - Requires a schemeful base URL: "http(s)://host[:port]" ("host:port" will fail).
// - The http.Client must speak HTTP/2 (TLS) or h2c (plaintext AllowHTTP+DialTLS) for classic gRPC.
func NewConnectGRPCReplicationAPIClient(
	ctx context.Context,
	httpAddress string,
	extraDialOpts ...connect.ClientOption,
) (message_apiconnect.ReplicationApiClient, error) {
	target, isTLS, err := HTTPAddressToConnectProtocolTarget(httpAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid http address: %w", err)
	}

	// Classic gRPC requires HTTP/2; use TLS (h2) or plaintext h2c depending on the URL scheme.
	httpClient, err := BuildHTTP2Client(ctx, isTLS)
	if err != nil {
		return nil, fmt.Errorf("failed to build http client: %w", err)
	}

	opts := BuildGRPCDialOptions(extraDialOpts...)

	return message_apiconnect.NewReplicationApiClient(httpClient, target, opts...), nil
}

// NewConnectGRPCWebReplicationAPIClient builds a Connect-based gRPC-Web Replication API client.
// - Use connect.WithGRPCWeb() (not WithGRPC()).
// - gRPC-Web works over HTTP/1.1 and HTTP/2; a standard http.Client is sufficient.
// - Base URL must be schemeful: "http(s)://host[:port]" ("host:port" will fail).
func NewConnectGRPCWebReplicationAPIClient(
	ctx context.Context,
	httpAddress string,
	extraDialOpts ...connect.ClientOption,
) (message_apiconnect.ReplicationApiClient, error) {
	target, isTLS, err := HTTPAddressToConnectProtocolTarget(httpAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid http address: %w", err)
	}

	// gRPC-Web works over HTTP/1.1 and HTTP/2; use an HTTP/2-capable client for consistency.
	httpClient, err := BuildHTTP2Client(ctx, isTLS)
	if err != nil {
		return nil, fmt.Errorf("failed to build http client: %w", err)
	}

	opts := BuildGRPCWebDialOptions(extraDialOpts...)

	return message_apiconnect.NewReplicationApiClient(httpClient, target, opts...), nil
}

// NewConnectGatewayAPIClient builds a Connect-based Gateway API client.
// The consumer is responsible of passing any extra dial options, to make the
// connection gRPC or gRPC-Web.
func NewConnectGatewayAPIClient(
	ctx context.Context,
	httpAddress string,
	extraDialOpts ...connect.ClientOption,
) (payer_apiconnect.PayerApiClient, error) {
	target, isTLS, err := HTTPAddressToConnectProtocolTarget(httpAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid http address: %w", err)
	}

	httpClient, err := BuildHTTP2Client(ctx, isTLS)
	if err != nil {
		return nil, fmt.Errorf("failed to build http client: %w", err)
	}

	opts := BuildConnectProtocolDialOptions(extraDialOpts...)

	return payer_apiconnect.NewPayerApiClient(httpClient, target, opts...), nil
}

// NewConnectMetadataAPIClient builds a Connect-based Metadata API client.
// The consumer is responsible of passing any extra dial options, to make the
// connection gRPC or gRPC-Web.
func NewConnectMetadataAPIClient(
	ctx context.Context,
	httpAddress string,
	extraDialOpts ...connect.ClientOption,
) (metadata_apiconnect.MetadataApiClient, error) {
	target, isTLS, err := HTTPAddressToConnectProtocolTarget(httpAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid http address: %w", err)
	}

	httpClient, err := BuildHTTP2Client(ctx, isTLS)
	if err != nil {
		return nil, fmt.Errorf("failed to build http client: %w", err)
	}

	return metadata_apiconnect.NewMetadataApiClient(httpClient, target, extraDialOpts...), nil
}

/* Native gRPC clients */

// NewGRPCConn builds a native grpc-go client for the given HTTP address.
//   - Uses the standard grpc-go library (not connect-go).
//   - Requires a schemeful base URL (http:// or https://).
//   - For Connect-based gRPC clients, use NewConnectGRPCReplicationAPIClient instead.
//
// Developer Note: Upstream caller is responsible for closing the returned connection.
func NewGRPCConn(
	httpAddress string,
	extraDialOpts ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	target, isTLS, err := HTTPAddressToGRPCTarget(httpAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP address: %w", err)
	}

	dialOptions := append([]grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(maxMessageSize),
			grpc.MaxCallRecvMsgSize(maxMessageSize),
		),
	}, extraDialOpts...)

	if isTLS {
		tlsConfig, err := buildTLSConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to build TLS config: %w", err)
		}

		dialOptions = append(
			dialOptions,
			grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		)
	} else {
		// h2c: plaintext HTTP/2
		dialOptions = append(dialOptions,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, "tcp", addr)
			}),
		)
	}

	conn, err := grpc.NewClient(target, dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc connection: %w", err)
	}

	return conn, nil
}

// BuildHTTP2Client creates an HTTP/2 client for Connect-based RPC or other HTTP/2 services.
func BuildHTTP2Client(ctx context.Context, isTLS bool) (*http.Client, error) {
	dialer := &net.Dialer{
		Timeout: clientTimeout,
	}

	if isTLS {
		tlsConfig, err := buildTLSConfig()
		if err != nil {
			return nil, err
		}

		return &http.Client{
			Transport: &http.Transport{
				ForceAttemptHTTP2: true,
				TLSClientConfig:   tlsConfig,
			},
			Timeout: clientTimeout,
		}, nil
	}

	// h2c for plaintext HTTP/2
	transport := &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		},
		ReadIdleTimeout: readIdleTimeout,
		PingTimeout:     pingTimeout,
	}
	return &http.Client{Transport: transport, Timeout: clientTimeout}, nil
}

/* Utils */

// BuildConnectProtocolDialOptions builds the default dial options for a Connect-Go, gRPC or gRPC-Web connection.
// Internal node <-> node communication can rely on this protocol.
func BuildConnectProtocolDialOptions(extraDialOpts ...connect.ClientOption) []connect.ClientOption {
	return getBaseDialOptions(extraDialOpts...)
}

// BuildGRPCDialOptions instructs the client to use the gRPC transport.
// Ideal for client <-> node communication, where the client only implements gRPC (i.e., Tonic).
func BuildGRPCDialOptions(extraDialOpts ...connect.ClientOption) []connect.ClientOption {
	options := []connect.ClientOption{
		connect.WithGRPC(),
	}
	options = append(options, extraDialOpts...)
	return getBaseDialOptions(options...)
}

// BuildGRPCWebDialOptions instructs the client to use the gRPC-Web transport.
// Ideal for WASM clients that need to use the gRPC-Web protocol.
func BuildGRPCWebDialOptions(extraDialOpts ...connect.ClientOption) []connect.ClientOption {
	options := []connect.ClientOption{
		connect.WithGRPCWeb(),
	}
	options = append(options, extraDialOpts...)
	return getBaseDialOptions(options...)
}

// HTTPAddressToGRPCTarget maps from a URL, as defined in https://pkg.go.dev/net/url#URL, to a gRPC target,
// as defined in https://github.com/grpc/grpc/blob/master/doc/naming.md
// Use only with clients with classic gRPC bindings.
func HTTPAddressToGRPCTarget(httpAddress string) (target string, isTLS bool, err error) {
	parsedURL, err := url.Parse(httpAddress)
	if err != nil {
		return "", false, err
	}

	switch parsedURL.Scheme {
	case "https":
		isTLS = true
	case "http", "":
		isTLS = false
	default:
		return "", false, fmt.Errorf("unknown connection schema %s", parsedURL.Scheme)
	}

	host := parsedURL.Hostname()
	if host == "" {
		return "", false, fmt.Errorf("missing host in address %q", httpAddress)
	}

	if parsedURL.Port() != "" {
		return fmt.Sprintf("%s:%s", host, parsedURL.Port()), isTLS, nil
	}

	return host, isTLS, nil
}

// HTTPAddressToConnectProtocolTarget maps from a URL to a Connect-Go target.
func HTTPAddressToConnectProtocolTarget(httpAddress string) (target string, isTLS bool, err error) {
	parsedURL, err := url.Parse(httpAddress)
	if err != nil {
		return "", false, err
	}

	host := parsedURL.Hostname()
	if host == "" {
		return "", false, fmt.Errorf("missing host in address %q", httpAddress)
	}

	scheme := parsedURL.Scheme
	if scheme == "" {
		scheme = "http"
	}

	switch scheme {
	case "https":
		isTLS = true
	case "http":
		isTLS = false
	default:
		return "", false, fmt.Errorf("unknown connection scheme: %s", parsedURL.Scheme)
	}

	if parsedURL.Port() != "" {
		return fmt.Sprintf(
			"%s://%s:%s",
			scheme,
			host,
			parsedURL.Port(),
		), isTLS, nil
	}

	return fmt.Sprintf("%s://%s", scheme, host), isTLS, nil
}

// getBaseDialOptions builds the default dial options for a Connect-Go, gRPC or gRPC-Web connection.
func getBaseDialOptions(extraDialOpts ...connect.ClientOption) []connect.ClientOption {
	// TODO: Extend with compression options?
	return append([]connect.ClientOption{
		connect.WithReadMaxBytes(maxMessageSize),
		connect.WithSendMaxBytes(maxMessageSize),
		connect.WithSendGzip(),
	}, extraDialOpts...)
}

// buildTLSConfig generates a TLS config.
// Note: If it's needed to use mutual TLS later, extend buildTLSConfig with client certs.
func buildTLSConfig() (*tls.Config, error) {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system CA certificates: %w", err)
	}

	if certPool == nil {
		return nil, errors.New("no system CA certificates available")
	}

	return &tls.Config{
		RootCAs: certPool,
	}, nil
}
