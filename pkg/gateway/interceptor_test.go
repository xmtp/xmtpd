package gateway

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestGatewayInterceptor(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Identity Injection", func(t *testing.T) {
		expectedIdentity := Identity{
			Kind:     identityKindIP,
			Identity: "192.168.1.1",
		}

		interceptor := NewGatewayInterceptor(
			logger,
			func(ctx context.Context) (Identity, error) {
				return expectedIdentity, nil
			},
			nil,
		)

		unaryInterceptor := interceptor.Unary()

		var capturedCtx context.Context
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			capturedCtx = ctx
			return nil, nil
		}

		info := &grpc.UnaryServerInfo{
			FullMethod: "/test/method",
		}

		_, err := unaryInterceptor(context.Background(), nil, info, handler)
		require.NoError(t, err)

		// Check that identity was injected into context
		identity, ok := IdentityFromContext(capturedCtx)
		assert.True(t, ok)
		assert.Equal(t, expectedIdentity, identity)
	})

	t.Run("Authorization Success", func(t *testing.T) {
		identity := Identity{
			Kind:     identityKindIP,
			Identity: "127.0.0.1",
		}

		authorizerCalled := false
		authorizer := func(ctx context.Context, id Identity, req PublishRequestSummary) (bool, error) {
			authorizerCalled = true
			assert.Equal(t, identity, id)
			assert.Equal(t, 2, req.TotalEnvelopes)
			return true, nil
		}

		interceptor := NewGatewayInterceptor(
			logger,
			func(ctx context.Context) (Identity, error) {
				return identity, nil
			},
			[]AuthorizePublishFn{authorizer},
		)

		unaryInterceptor := interceptor.Unary()

		req := &payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopes.ClientEnvelope{{}, {}},
		}

		info := &grpc.UnaryServerInfo{
			FullMethod: "/xmtp.xmtpv4.payer_api.PayerApi/PublishClientEnvelopes",
		}

		handlerCalled := false
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return nil, nil
		}

		_, err := unaryInterceptor(context.Background(), req, info, handler)
		require.NoError(t, err)
		assert.True(t, authorizerCalled)
		assert.True(t, handlerCalled)
	})

	t.Run("Authorization Denied", func(t *testing.T) {
		identity := Identity{
			Kind:     identityKindIP,
			Identity: "192.168.1.1",
		}

		authorizer := func(ctx context.Context, id Identity, req PublishRequestSummary) (bool, error) {
			return false, nil
		}

		interceptor := NewGatewayInterceptor(
			logger,
			func(ctx context.Context) (Identity, error) {
				return identity, nil
			},
			[]AuthorizePublishFn{authorizer},
		)

		unaryInterceptor := interceptor.Unary()

		req := &payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopes.ClientEnvelope{{}},
		}

		info := &grpc.UnaryServerInfo{
			FullMethod: "/xmtp.xmtpv4.payer_api.PayerApi/PublishClientEnvelopes",
		}

		handlerCalled := false
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return nil, nil
		}

		_, err := unaryInterceptor(context.Background(), req, info, handler)
		require.Error(t, err)
		assert.False(t, handlerCalled) // Handler should not be called
		assert.Contains(t, err.Error(), "PermissionDenied")
	})

	t.Run("Authorization Not Called For Other Methods", func(t *testing.T) {
		identity := Identity{
			Kind:     identityKindIP,
			Identity: "192.168.1.1",
		}

		authorizerCalled := false
		authorizer := func(ctx context.Context, id Identity, req PublishRequestSummary) (bool, error) {
			authorizerCalled = true
			return true, nil
		}

		interceptor := NewGatewayInterceptor(
			logger,
			func(ctx context.Context) (Identity, error) {
				return identity, nil
			},
			[]AuthorizePublishFn{authorizer},
		)

		unaryInterceptor := interceptor.Unary()

		// Test with a different method
		req := &payer_api.PublishClientEnvelopesRequest{
			Envelopes: []*envelopes.ClientEnvelope{{}},
		}

		info := &grpc.UnaryServerInfo{
			FullMethod: "/xmtp.xmtpv4.payer_api.PayerApi/SomeOtherMethod",
		}

		handlerCalled := false
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return nil, nil
		}

		_, err := unaryInterceptor(context.Background(), req, info, handler)
		require.NoError(t, err)
		assert.False(t, authorizerCalled) // Authorizer should not be called
		assert.True(t, handlerCalled)     // Handler should still be called
	})
}
