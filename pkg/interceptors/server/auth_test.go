package server_test

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/xmtp/xmtpd/pkg/authn"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/interceptors/client"
	"github.com/xmtp/xmtpd/pkg/interceptors/server"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/utils"
)

func TestServer_EnforcesAuth(t *testing.T) {
	var (
		ctx          = t.Context()
		log          = testutils.NewLog(t)
		serverNodeID = uint32(200)

		verifier = newSimpleVerifier()

		serverAuth = server.NewServerAuthInterceptor(log, verifier, server.RequireToken(true))
		server     = &replicationServer{}

		emptyQueryResponse = &connect.Response[message_api.QueryEnvelopesResponse]{
			Msg: &message_api.QueryEnvelopesResponse{
				Envelopes: []*envelopes.OriginatorEnvelope{},
			},
		}
	)

	server.query = func(ctx context.Context, req *connect.Request[message_api.QueryEnvelopesRequest],
	) (*connect.Response[message_api.QueryEnvelopesResponse], error) {
		// Server should have the information about authentication embedded.
		verified, ok := ctx.Value(constants.VerifiedNodeRequestCtxKey{}).(bool)
		require.True(t, ok)
		require.True(t, verified)

		return emptyQueryResponse, nil
	}

	httpServer, addr := newReplicationHTTPServer(server, serverAuth)
	defer httpServer.Close()

	// NOTE: Do not paralelize these tests as they modify the same server.

	t.Run("nominal case", func(t *testing.T) {
		// NOTE: Server will require and validate the JWT token.

		var (
			clientID     = uint32(100)
			clientKey    = testutils.RandomPrivateKey(t)
			tokenFactory = authn.NewTokenFactory(clientKey, clientID, nil)
			clientAuth   = client.NewClientAuthInterceptor(tokenFactory, serverNodeID)

			start = time.Now()
		)

		// Have the verifier actually process the token and validate the content.
		verifier.verify = func(token string) (uint32, authn.CloseFunc, error) {
			subject := strconv.FormatUint(uint64(clientID), 10)
			validateJWT(t, token, start, serverNodeID, subject, &clientKey.PublicKey)
			return clientID, func() {}, nil
		}

		client, err := utils.NewConnectReplicationAPIClient(ctx, addr,
			connect.WithInterceptors(clientAuth))
		require.NoError(t, err)

		_, err = client.QueryEnvelopes(
			ctx,
			connect.NewRequest(&message_api.QueryEnvelopesRequest{}))
		require.NoError(t, err)
	})
	t.Run("request without a token is rejected", func(t *testing.T) {
		client, err := utils.NewConnectReplicationAPIClient(ctx, addr)
		require.NoError(t, err)

		_, err = client.QueryEnvelopes(
			ctx,
			connect.NewRequest(&message_api.QueryEnvelopesRequest{}))

		var connErr *connect.Error
		require.ErrorAs(t, err, &connErr)
		require.Equal(t, connect.CodeUnauthenticated, connErr.Code())
	})
	t.Run("request not approved by verifier is rejected", func(t *testing.T) {
		// NOTE: Client has proper authentication token but the verifier rejects it.

		var (
			clientNodeID = uint32(100)
			clientKey    = testutils.RandomPrivateKey(t)
			tokenFactory = authn.NewTokenFactory(clientKey, clientNodeID, nil)
			clientAuth   = client.NewClientAuthInterceptor(tokenFactory, serverNodeID)

			verifierCalls = 0
		)

		client, err := utils.NewConnectReplicationAPIClient(ctx, addr,
			connect.WithInterceptors(clientAuth))
		require.NoError(t, err)

		// Baseline - verifier approves the request.
		verifier.verify = func(string) (uint32, authn.CloseFunc, error) {
			verifierCalls += 1
			return 0, func() {}, nil
		}

		_, err = client.QueryEnvelopes(
			ctx,
			connect.NewRequest(&message_api.QueryEnvelopesRequest{}))
		require.NoError(t, err)
		require.Equal(t, 1, verifierCalls)

		// Now, tell the verifier to not approve this request.
		verifier.verify = func(string) (uint32, authn.CloseFunc, error) {
			verifierCalls += 1
			return 0, func() {}, errors.New("rejected")
		}

		_, err = client.QueryEnvelopes(
			ctx,
			connect.NewRequest(&message_api.QueryEnvelopesRequest{}))
		var connErr *connect.Error
		require.ErrorAs(t, err, &connErr)
		require.Equal(t, connect.CodeUnauthenticated, connErr.Code())

		require.Equal(t, 2, verifierCalls)
	})
}

func TestServer_RelaxedAuth(t *testing.T) {
	var (
		ctx          = t.Context()
		log          = testutils.NewLog(t)
		serverNodeID = uint32(200)

		verifier = newSimpleVerifier()

		serverAuth = server.NewServerAuthInterceptor(log, verifier)
		server     = &replicationServer{}

		emptyQueryResponse = &connect.Response[message_api.QueryEnvelopesResponse]{
			Msg: &message_api.QueryEnvelopesResponse{
				Envelopes: []*envelopes.OriginatorEnvelope{},
			},
		}
	)

	server.query = func(ctx context.Context, req *connect.Request[message_api.QueryEnvelopesRequest],
	) (*connect.Response[message_api.QueryEnvelopesResponse], error) {
		return emptyQueryResponse, nil
	}

	httpServer, addr := newReplicationHTTPServer(server, serverAuth)
	defer httpServer.Close()

	// NOTE: Do not paralelize these tests as they modify the same server.
	t.Run("token presence is not enforced", func(t *testing.T) {
		server.query = func(ctx context.Context, req *connect.Request[message_api.QueryEnvelopesRequest],
		) (*connect.Response[message_api.QueryEnvelopesResponse], error) {
			return emptyQueryResponse, nil
		}

		client, err := utils.NewConnectReplicationAPIClient(ctx, addr)
		require.NoError(t, err)

		_, err = client.QueryEnvelopes(
			ctx,
			connect.NewRequest(&message_api.QueryEnvelopesRequest{}))
		require.NoError(t, err)
	})
	t.Run("if token is present, it must be validated", func(t *testing.T) {
		// NOTE: If the token is present, requests not passing validation should be rejected.
		var (
			clientNodeID = uint32(100)
			clientKey    = testutils.RandomPrivateKey(t)
			tokenFactory = authn.NewTokenFactory(clientKey, clientNodeID, nil)
			clientAuth   = client.NewClientAuthInterceptor(tokenFactory, serverNodeID)
		)

		client, err := utils.NewConnectReplicationAPIClient(ctx, addr,
			connect.WithInterceptors(clientAuth))
		require.NoError(t, err)

		// Now, tell the verifier to not approve this request.
		verifier.verify = func(string) (uint32, authn.CloseFunc, error) {
			return 0, func() {}, errors.New("rejected")
		}

		_, err = client.QueryEnvelopes(
			ctx,
			connect.NewRequest(&message_api.QueryEnvelopesRequest{}))
		var connErr *connect.Error
		require.ErrorAs(t, err, &connErr)
		require.Equal(t, connect.CodeUnauthenticated, connErr.Code())
	})
}

type queryEnvelopesFunc func(
	context.Context,
	*connect.Request[message_api.QueryEnvelopesRequest],
) (*connect.Response[message_api.QueryEnvelopesResponse], error)

type replicationServer struct {
	message_apiconnect.ReplicationApiHandler

	query queryEnvelopesFunc
}

func (s *replicationServer) QueryEnvelopes(
	ctx context.Context,
	req *connect.Request[message_api.QueryEnvelopesRequest],
) (*connect.Response[message_api.QueryEnvelopesResponse], error) {
	return s.query(ctx, req)
}

func newReplicationHTTPServer(
	server *replicationServer,
	interceptors ...connect.Interceptor,
) (*httptest.Server, string) {
	path, handler := message_apiconnect.NewReplicationApiHandler(
		server,
		connect.WithInterceptors(interceptors...),
	)

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	h2cHandler := h2c.NewHandler(mux, &http2.Server{
		IdleTimeout: 5 * time.Minute,
	})

	httpServer := httptest.NewServer(h2cHandler)

	return httpServer, httpServer.URL
}

type simpleVerifier struct {
	verify func(string) (uint32, authn.CloseFunc, error)
}

func newSimpleVerifier() *simpleVerifier {
	return &simpleVerifier{}
}

func (s *simpleVerifier) Verify(token string) (uint32, authn.CloseFunc, error) {
	return s.verify(token)
}

func verifyTokenAudience(t *testing.T, token *jwt.Token, nodeID uint32) {
	t.Helper()

	audience, err := token.Claims.GetAudience()
	require.NoError(t, err)
	require.Len(t, audience, 1)
	require.Equal(t, strconv.FormatUint(uint64(nodeID), 10), audience[0])
}

func validateJWT(
	t *testing.T,
	token string,
	start time.Time,
	serverID uint32,
	subject string,
	clientKey *ecdsa.PublicKey,
) {
	t.Helper()

	keyfunc := func(_ *jwt.Token) (any, error) {
		return clientKey, nil
	}

	parsed, err := jwt.ParseWithClaims(token, &authn.XmtpdClaims{}, keyfunc)
	require.NoError(t, err)

	require.True(t, parsed.Valid)

	// Verify token audience.
	verifyTokenAudience(t, parsed, serverID)

	// Verify that IssuedAt is the same or larger than test start.
	iat, err := parsed.Claims.GetIssuedAt()
	require.NoError(t, err)

	cmp := start.Round(time.Second).Compare(iat.Time)
	require.LessOrEqual(t, cmp, 0)

	// Verify that Expiration time is after test start.
	exp, err := parsed.Claims.GetExpirationTime()
	require.NoError(t, err)
	require.True(t, exp.After(start))

	// Verify Subject is the client node ID.
	tokenSubject, err := parsed.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, subject, tokenSubject)
}
