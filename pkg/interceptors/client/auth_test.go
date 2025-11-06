package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/testutils"
	apiTestUtils "github.com/xmtp/xmtpd/pkg/testutils/api"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func newMockReplicationAPIServer(
	token *authn.Token,
) (server *httptest.Server, addr string) {
	// Mock handler for the replication API.
	path, handler := message_apiconnect.NewReplicationApiHandler(
		&mockReplicationAPIServer{expectedToken: token.SignedString},
	)

	// Create a new HTTP mux to serve the API handlers.
	mux := http.NewServeMux()
	mux.Handle(path, handler)

	// Allow HTTP/2 and HTTP/1.1 connections.
	h2cHandler := h2c.NewHandler(mux, &http2.Server{
		IdleTimeout: 5 * time.Minute,
	})

	// Create the HTTP server to serve the API handlers.
	server = httptest.NewServer(h2cHandler)

	return server, strings.TrimPrefix(server.URL, "http://")
}

func TestAuthInterceptor(t *testing.T) {
	var (
		privateKey        = testutils.RandomPrivateKey(t)
		myNodeID          = uint32(100)
		targetNodeID      = uint32(200)
		wrongTargetNodeID = uint32(300)
		tokenFactory      = authn.NewTokenFactory(privateKey, myNodeID, nil)
		interceptorHappy  = NewClientAuthInterceptor(tokenFactory, targetNodeID)
		interceptorFail   = NewClientAuthInterceptor(tokenFactory, wrongTargetNodeID)
	)

	token, err := interceptorHappy.getToken()
	require.NoError(t, err)

	// Create a mock server to serve the API handlers.
	server, addr := newMockReplicationAPIServer(token)
	defer server.Close()

	// Happy path: Create client with interceptor, should succeed its queries.
	client := apiTestUtils.NewTestGRPCReplicationAPIClient(
		t,
		addr,
		connect.WithInterceptors(interceptorHappy),
	)

	// Call the unary method and check the response.
	_, err = client.QueryEnvelopes(
		t.Context(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{}),
	)
	require.NoError(t, err)

	// Sad path: Create another client without the interceptor, should fail its queries.
	client = apiTestUtils.NewTestGRPCReplicationAPIClient(
		t,
		addr,
	)

	_, err = client.QueryEnvelopes(
		t.Context(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{}),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "authorization token is not provided")

	// Sad path: Create another client with the wrong target node ID, should fail its queries.
	client = apiTestUtils.NewTestGRPCReplicationAPIClient(
		t,
		addr,
		connect.WithInterceptors(interceptorFail),
	)

	_, err = client.QueryEnvelopes(
		t.Context(),
		connect.NewRequest(&message_api.QueryEnvelopesRequest{}),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid authorization token")
}
