package sync

import (
	"io"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	envUtils "github.com/xmtp/xmtpd/pkg/envelopes"
	messageApiMocks "github.com/xmtp/xmtpd/pkg/mocks/message_api"
	payerreportMocks "github.com/xmtp/xmtpd/pkg/mocks/payerreport"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	feesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/fees"
	registryTestUtils "github.com/xmtp/xmtpd/pkg/testutils/registry"
)

var payerReportDomainSeparator = testutils.RandomDomainSeparator()

func mockSubscriptionOnePage(
	t *testing.T,
	envs []*envelopes.OriginatorEnvelope,
) message_api.ReplicationApi_SubscribeEnvelopesClient {
	stream := messageApiMocks.NewMockReplicationApi_SubscribeEnvelopesClient(t)
	hasSent := false
	stream.EXPECT().
		Recv().
		RunAndReturn(func() (*message_api.SubscribeEnvelopesResponse, error) {
			if hasSent {
				return nil, io.EOF
			}
			hasSent = true
			return &message_api.SubscribeEnvelopesResponse{
				Envelopes: envs,
			}, nil
		})
	return stream
}

func newTestOriginatorStream(
	t *testing.T,
	node *registry.Node,
	stream message_api.ReplicationApi_SubscribeEnvelopesClient,
	lastEnvelope *envUtils.OriginatorEnvelope,
) *originatorStream {
	log := testutils.NewLog(t)
	calculator := feesTestUtils.NewTestFeeCalculator()
	db, _ := testutils.NewDB(t, t.Context())

	return newOriginatorStream(
		t.Context(),
		db,
		log,
		node,
		lastEnvelope,
		stream,
		calculator,
		payerreportMocks.NewMockIPayerReportStore(t),
		payerReportDomainSeparator,
	)
}

func getAllMessagesForOriginator(
	t *testing.T,
	originatorStream *originatorStream,
) []queries.GatewayEnvelope {
	envs, err := originatorStream.queries.SelectGatewayEnvelopes(
		t.Context(),
		queries.SelectGatewayEnvelopesParams{
			OriginatorNodeIds: []int32{int32(originatorStream.node.NodeID)},
		},
	)
	require.NoError(t, err)
	return envs
}

func TestSyncWorkerSuccess(t *testing.T) {
	nodeID := uint32(200)
	sequenceID := uint64(100)
	envelope := envelopeTestUtils.CreateOriginatorEnvelope(t, nodeID, sequenceID)
	stream := mockSubscriptionOnePage(t, []*envelopes.OriginatorEnvelope{envelope})

	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	origStream := newTestOriginatorStream(t, &node, stream, nil)

	err := origStream.listen()
	var retryAfter *backoff.RetryAfterError
	require.ErrorAs(t, err, &retryAfter)
	require.Equal(t, retryAfter.Duration.Seconds(), float64(1))

	require.Eventually(t, func() bool {
		envs := getAllMessagesForOriginator(t, origStream)
		return len(envs) == 1 && envs[0].OriginatorSequenceID == int64(sequenceID)
	}, 1*time.Second, 50*time.Millisecond)
}

func TestSyncWorkerIgnoresInvalidEnvelopes(t *testing.T) {
	nodeID := uint32(200)
	sequenceID := uint64(100)
	envelope := envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(
		t,
		nodeID,
		sequenceID,
		[]byte("broken"),
	)

	stream := mockSubscriptionOnePage(t, []*envelopes.OriginatorEnvelope{envelope})
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	origStream := newTestOriginatorStream(t, &node, stream, nil)

	err := origStream.listen()
	var retryAfter *backoff.RetryAfterError
	require.ErrorAs(t, err, &retryAfter)

	// Give the write worker a chance to save the envelope
	time.Sleep(50 * time.Millisecond)
	envs := getAllMessagesForOriginator(t, origStream)
	require.Len(t, envs, 0)
}
