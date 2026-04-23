package sync

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	messageApiMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/message_api"
)

func TestEnvelopesStreamAdapter_ReturnsEnvelopes(t *testing.T) {
	mockStream := messageApiMocks.NewMockReplicationApi_SubscribeEnvelopesClient[*message_api.SubscribeEnvelopesResponse](
		t,
	)
	env := &envelopes.OriginatorEnvelope{}
	mockStream.EXPECT().Recv().Return(&message_api.SubscribeEnvelopesResponse{
		Envelopes: []*envelopes.OriginatorEnvelope{env},
	}, nil).Once()

	adapter := &envelopesStreamAdapter{stream: mockStream}
	got, err := adapter.Recv()
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Same(t, env, got[0])
}

func TestEnvelopesStreamAdapter_PropagatesError(t *testing.T) {
	mockStream := messageApiMocks.NewMockReplicationApi_SubscribeEnvelopesClient[*message_api.SubscribeEnvelopesResponse](
		t,
	)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Once()

	adapter := &envelopesStreamAdapter{stream: mockStream}
	_, err := adapter.Recv()
	require.ErrorIs(t, err, io.EOF)
}

func TestEnvelopesStreamAdapter_CloseSendDelegates(t *testing.T) {
	mockStream := messageApiMocks.NewMockReplicationApi_SubscribeEnvelopesClient[*message_api.SubscribeEnvelopesResponse](
		t,
	)
	wantErr := errors.New("close failed")
	mockStream.EXPECT().CloseSend().Return(wantErr).Once()

	adapter := &envelopesStreamAdapter{stream: mockStream}
	require.ErrorIs(t, adapter.CloseSend(), wantErr)
}

func TestOriginatorsStreamAdapter_ReturnsEnvelopes(t *testing.T) {
	mockStream := messageApiMocks.NewMockReplicationApi_SubscribeOriginatorsClient[*message_api.SubscribeOriginatorsResponse](
		t,
	)
	env := &envelopes.OriginatorEnvelope{}
	mockStream.EXPECT().Recv().Return(&message_api.SubscribeOriginatorsResponse{
		Response: &message_api.SubscribeOriginatorsResponse_Envelopes_{
			Envelopes: &message_api.SubscribeOriginatorsResponse_Envelopes{
				Envelopes: []*envelopes.OriginatorEnvelope{env},
			},
		},
	}, nil).Once()

	adapter := &originatorsStreamAdapter{stream: mockStream}
	got, err := adapter.Recv()
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Same(t, env, got[0])
}

func TestOriginatorsStreamAdapter_KeepaliveReturnsEmpty(t *testing.T) {
	mockStream := messageApiMocks.NewMockReplicationApi_SubscribeOriginatorsClient[*message_api.SubscribeOriginatorsResponse](
		t,
	)
	// Keepalive: empty response with no oneof set.
	mockStream.EXPECT().Recv().Return(&message_api.SubscribeOriginatorsResponse{}, nil).Once()

	adapter := &originatorsStreamAdapter{stream: mockStream}
	got, err := adapter.Recv()
	require.NoError(t, err)
	require.Empty(t, got)
}

func TestOriginatorsStreamAdapter_PropagatesError(t *testing.T) {
	mockStream := messageApiMocks.NewMockReplicationApi_SubscribeOriginatorsClient[*message_api.SubscribeOriginatorsResponse](
		t,
	)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Once()

	adapter := &originatorsStreamAdapter{stream: mockStream}
	_, err := adapter.Recv()
	require.ErrorIs(t, err, io.EOF)
}

func TestOriginatorsStreamAdapter_CloseSendDelegates(t *testing.T) {
	mockStream := messageApiMocks.NewMockReplicationApi_SubscribeOriginatorsClient[*message_api.SubscribeOriginatorsResponse](
		t,
	)
	wantErr := errors.New("close failed")
	mockStream.EXPECT().CloseSend().Return(wantErr).Once()

	adapter := &originatorsStreamAdapter{stream: mockStream}
	require.ErrorIs(t, adapter.CloseSend(), wantErr)
}
