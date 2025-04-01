package client

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/api/message"
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// Create a mock implementation of the ReplicationApiServer interface
// but that embeds `UnimplementedReplicationApiServer` (which mockery won't do for us)
type mockReplicationApiServer struct {
	message.Service
	expectedToken string
}

func (s *mockReplicationApiServer) QueryEnvelopes(
	ctx context.Context,
	req *message_api.QueryEnvelopesRequest,
) (*message_api.QueryEnvelopesResponse, error) {
	// Get metadata from the context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	// Extract and verify the token
	tokens := md.Get(constants.NODE_AUTHORIZATION_HEADER_NAME)
	if len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}
	token := tokens[0]
	if token != s.expectedToken {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization token")
	}

	// You can add more assertions here to verify the token's content
	// For example, you might want to decode the token and check its claims
	return &message_api.QueryEnvelopesResponse{}, nil
}

func TestAuthInterceptor(t *testing.T) {
	privateKey := testutils.RandomPrivateKey(t)
	myNodeID := uint32(100)
	targetNodeID := uint32(200)
	tokenFactory := authn.NewTokenFactory(privateKey, myNodeID, nil)
	interceptor := NewAuthInterceptor(tokenFactory, targetNodeID)
	token, err := interceptor.getToken()
	require.NoError(t, err)

	// Use a bufconn listener to simulate a gRPC connection without actually dialing
	listener := bufconn.Listen(1024 * 1024)

	// Register the mock service on the server
	server := grpc.NewServer()
	message_api.RegisterReplicationApiServer(
		server,
		&mockReplicationApiServer{expectedToken: token.SignedString},
	)

	// Start the gRPC server in a goroutine
	go func() {
		if err := server.Serve(listener); err != nil {
			t.Fail()
		}
	}()

	t.Cleanup(func() {
		server.Stop()
		_ = listener.Close()
	})

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
