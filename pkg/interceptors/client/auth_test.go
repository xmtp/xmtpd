package client

import (
	"context"
	"net"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// Create a mock implementation of the ReplicationApiServer interface
// but that embeds `message_apiconnect.ReplicationApiHandler` (which mockery won't do for us)
type mockReplicationAPIServer struct {
	message_apiconnect.ReplicationApiHandler
	expectedToken string
}

func (s *mockReplicationAPIServer) QueryEnvelopes(
	ctx context.Context,
	req *connect.Request[message_api.QueryEnvelopesRequest],
) (*connect.Response[message_api.QueryEnvelopesResponse], error) {
	// Extract and verify the token
	token := req.Header().Get(constants.NodeAuthorizationHeaderName)
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	if token != s.expectedToken {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization token")
	}

	// You can add more assertions here to verify the token's content
	// For example, you might want to decode the token and check its claims.
	return &connect.Response[message_api.QueryEnvelopesResponse]{
		Msg: &message_api.QueryEnvelopesResponse{},
	}, nil
}

func TestAuthInterceptor(t *testing.T) {
	var (
		privateKey   = testutils.RandomPrivateKey(t)
		myNodeID     = uint32(100)
		targetNodeID = uint32(200)
		tokenFactory = authn.NewTokenFactory(privateKey, myNodeID, nil)
		interceptor  = NewClientAuthInterceptor(tokenFactory, targetNodeID)
	)

	token, err := interceptor.getToken()
	require.NoError(t, err)

	// Use a bufconn listener to simulate a gRPC connection without actually dialing
	listener := bufconn.Listen(1024 * 1024)

	_, handler := message_apiconnect.NewReplicationApiHandler(
		&mockReplicationAPIServer{expectedToken: token.SignedString},
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	// Connect to the fake server and set the right interceptors
	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
	)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Create a client with the connection
	client := message_api.NewReplicationApiClient(conn)

	// Call the unary method and check the response
	_, err = client.QueryEnvelopes(context.Background(), &message_api.QueryEnvelopesRequest{})
	require.NoError(t, err)

	// Create another client without the interceptor
	connWithoutInterceptor, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer func() {
		_ = connWithoutInterceptor.Close()
	}()

	client = message_api.NewReplicationApiClient(connWithoutInterceptor)

	// Call the unary method and check the response
	_, err = client.QueryEnvelopes(context.Background(), &message_api.QueryEnvelopesRequest{})
	require.Error(t, err)
}
