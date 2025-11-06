package sync

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"

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
	stream := messageApiMocks.NewMockReplicationApi_SubscribeEnvelopesClient[*message_api.SubscribeEnvelopesResponse](
		t,
	)
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

func newTestEnvelopeSink(
	t *testing.T,
	writeQueue chan *envUtils.OriginatorEnvelope,
	ctx context.Context,
) *EnvelopeSink {
	log := testutils.NewLog(t)
	calculator := feesTestUtils.NewTestFeeCalculator()
	db, _ := testutils.NewDB(t, ctx)

	return newEnvelopeSink(
		ctx,
		db,
		log,
		calculator,
		payerreportMocks.NewMockIPayerReportStore(t),
		payerReportDomainSeparator,
		writeQueue,
	)
}

func newTestOriginatorStream(
	t *testing.T,
	node *registry.Node,
	stream message_api.ReplicationApi_SubscribeEnvelopesClient,
	cursor *cursor,
	writeQueue chan *envUtils.OriginatorEnvelope,
) *originatorStream {
	log := testutils.NewLog(t)

	return newOriginatorStream(
		t.Context(),
		log,
		node,
		cursor,
		stream,
		writeQueue,
	)
}

func getAllMessagesForOriginator(
	t *testing.T,
	storer *EnvelopeSink,
	nodeID uint32,
) []queries.GatewayEnvelopesView {
	envs, err := storer.queries.SelectGatewayEnvelopesByOriginators(
		t.Context(),
		queries.SelectGatewayEnvelopesByOriginatorsParams{
			OriginatorNodeIds: []int32{int32(nodeID)},
		},
	)
	require.NoError(t, err)
	return db.TransformRowsByOriginator(envs)
}

func TestSyncWorkerSuccess(t *testing.T) {
	nodeID := uint32(200)
	sequenceID := uint64(100)
	envelope := envelopeTestUtils.CreateOriginatorEnvelope(t, nodeID, sequenceID)
	stream := mockSubscriptionOnePage(t, []*envelopes.OriginatorEnvelope{envelope})

	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)
	defer close(writeQueue)
	dbStorerInstance := newTestEnvelopeSink(t, writeQueue, t.Context())
	go func() {
		dbStorerInstance.Start()
	}()
	origStream := newTestOriginatorStream(t, &node, stream, nil, writeQueue)

	err := origStream.listen()
	var retryAfter *backoff.RetryAfterError
	require.ErrorAs(t, err, &retryAfter)
	require.Equal(t, retryAfter.Duration.Seconds(), float64(1))

	require.Eventually(t, func() bool {
		envs := getAllMessagesForOriginator(t, dbStorerInstance, nodeID)
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

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)
	defer close(writeQueue)
	dbStorerInstance := newTestEnvelopeSink(t, writeQueue, t.Context())
	go func() {
		dbStorerInstance.Start()
	}()
	origStream := newTestOriginatorStream(t, &node, stream, nil, writeQueue)

	err := origStream.listen()
	var retryAfter *backoff.RetryAfterError
	require.ErrorAs(t, err, &retryAfter)

	// Give the write worker a chance to save the envelope
	time.Sleep(50 * time.Millisecond)
	envs := getAllMessagesForOriginator(t, dbStorerInstance, nodeID)
	require.Len(t, envs, 0)
}

func TestEnvelopeSinkShutdownViaClose(t *testing.T) {
	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)
	dbStorerInstance := newTestEnvelopeSink(t, writeQueue, t.Context())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		dbStorerInstance.Start()
	}()

	close(writeQueue)
	wg.Wait()
}

func TestEnvelopeSinkShutdownViaContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)
	dbStorerInstance := newTestEnvelopeSink(t, writeQueue, ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		dbStorerInstance.Start()
	}()

	cancel()

	wg.Wait()
}
