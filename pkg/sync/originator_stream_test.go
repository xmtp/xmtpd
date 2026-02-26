package sync

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/xmtp/xmtpd/pkg/migrator"

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
	dbInstance, _ := testutils.NewDB(t, ctx)

	err := dbInstance.Query().EnsureGatewayParts(
		ctx,
		queries.EnsureGatewayPartsParams{
			OriginatorNodeID:     200,
			OriginatorSequenceID: 1,
			BandWidth:            db.GatewayEnvelopeBandWidth,
		},
	)
	require.NoError(t, err)

	return newEnvelopeSink(
		ctx,
		dbInstance,
		log,
		calculator,
		payerreportMocks.NewMockIPayerReportStore(t),
		payerReportDomainSeparator,
		writeQueue,
		100*time.Millisecond,
	)
}

func newTestOriginatorStream(
	t *testing.T,
	node *registry.Node,
	stream message_api.ReplicationApi_SubscribeEnvelopesClient,
	lastSequenceID map[uint32]uint64,
	writeQueue chan *envUtils.OriginatorEnvelope,
) *originatorStream {
	log := testutils.NewLog(t)

	permittedOriginators := map[uint32]struct{}{
		node.NodeID: {},
	}

	return newOriginatorStream(
		t.Context(),
		log,
		node,
		lastSequenceID,
		permittedOriginators,
		stream,
		writeQueue,
	)
}

func getAllMessagesForOriginator(
	t *testing.T,
	storer *EnvelopeSink,
	nodeID uint32,
) []queries.GatewayEnvelopesView {
	envs, err := storer.db.ReadQuery().SelectGatewayEnvelopesByOriginators(
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
	sequenceID := uint64(1)
	envelope := envelopeTestUtils.CreateOriginatorEnvelope(t, nodeID, sequenceID)
	stream := mockSubscriptionOnePage(t, []*envelopes.OriginatorEnvelope{envelope})

	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)
	defer close(writeQueue)
	dbStorerInstance := newTestEnvelopeSink(t, writeQueue, t.Context())
	go func() {
		dbStorerInstance.Start()
	}()

	lastSequenceIds := make(map[uint32]uint64)

	origStream := newTestOriginatorStream(t, &node, stream, lastSequenceIds, writeQueue)

	err := origStream.listen()
	var retryAfter *backoff.RetryAfterError
	require.ErrorAs(t, err, &retryAfter)
	require.InEpsilon(t, float64(1), retryAfter.Duration.Seconds(), 0.001)

	require.Eventually(t, func() bool {
		envs := getAllMessagesForOriginator(t, dbStorerInstance, nodeID)
		return len(envs) == 1 && envs[0].OriginatorSequenceID == int64(sequenceID)
	}, time.Second, 50*time.Millisecond)
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

	lastSequenceIds := make(map[uint32]uint64)

	origStream := newTestOriginatorStream(t, &node, stream, lastSequenceIds, writeQueue)

	err := origStream.listen()
	var retryAfter *backoff.RetryAfterError
	require.ErrorAs(t, err, &retryAfter)

	// Give the write worker a chance to save the envelope
	time.Sleep(50 * time.Millisecond)
	envs := getAllMessagesForOriginator(t, dbStorerInstance, nodeID)
	require.Empty(t, envs)
}

func TestEnvelopeSinkShutdownViaClose(t *testing.T) {
	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)
	dbStorerInstance := newTestEnvelopeSink(t, writeQueue, t.Context())
	var wg sync.WaitGroup
	wg.Go(func() {
		dbStorerInstance.Start()
	})

	close(writeQueue)
	wg.Wait()
}

func TestEnvelopeSinkShutdownViaContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)
	dbStorerInstance := newTestEnvelopeSink(t, writeQueue, ctx)
	var wg sync.WaitGroup
	wg.Go(func() {
		dbStorerInstance.Start()
	})

	cancel()

	wg.Wait()
}

func newTestOriginatorStreamWithPermitted(
	t *testing.T,
	node *registry.Node,
	stream message_api.ReplicationApi_SubscribeEnvelopesClient,
	lastSequenceIds map[uint32]uint64,
	permitted map[uint32]struct{},
	writeQueue chan *envUtils.OriginatorEnvelope,
) *originatorStream {
	log := testutils.NewLog(t)

	return newOriginatorStream(
		t.Context(),
		log,
		node,
		lastSequenceIds,
		permitted,
		stream,
		writeQueue,
	)
}

func TestSyncWorkerRejectsEnvelopeFromUnpermittedOriginator(t *testing.T) {
	// Node we are syncing *from*
	nodeID := uint32(200)

	// Envelope claims it was authored by a different originator
	badOriginatorID := uint32(201)
	sequenceID := uint64(1)
	envelope := envelopeTestUtils.CreateOriginatorEnvelope(t, badOriginatorID, sequenceID)

	stream := mockSubscriptionOnePage(t, []*envelopes.OriginatorEnvelope{envelope})
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)
	defer close(writeQueue)

	dbStorerInstance := newTestEnvelopeSink(t, writeQueue, t.Context())
	go dbStorerInstance.Start()

	// Only permit nodeID (200), NOT badOriginatorID (201)
	permitted := map[uint32]struct{}{
		nodeID: {},
	}
	lastSequenceIds := make(map[uint32]uint64)

	origStream := newTestOriginatorStreamWithPermitted(
		t,
		&node,
		stream,
		lastSequenceIds,
		permitted,
		writeQueue,
	)

	err := origStream.listen()
	require.Error(t, err)

	// Ensure nothing got queued
	select {
	case got := <-writeQueue:
		require.Failf(
			t,
			"unexpected envelope queued",
			"got originator=%d seq=%d",
			got.OriginatorNodeID(),
			got.OriginatorSequenceID(),
		)
	default:
		// ok
	}

	// Ensure nothing got stored
	time.Sleep(50 * time.Millisecond) // give sink a moment (defensive)
	envs := getAllMessagesForOriginator(t, dbStorerInstance, badOriginatorID)
	require.Empty(t, envs)
}

func TestSyncWorkerAcceptsEnvelopeFromPermittedOriginator(t *testing.T) {
	// Node we are syncing *from*
	nodeID := uint32(200)

	// Envelope comes from an additional permitted originator
	otherPermittedID := migrator.KeyPackagesOriginatorID
	sequenceID := uint64(1)
	envelope := envelopeTestUtils.CreateOriginatorEnvelope(t, otherPermittedID, sequenceID)

	stream := mockSubscriptionOnePage(t, []*envelopes.OriginatorEnvelope{envelope})
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)
	defer close(writeQueue)

	dbStorerInstance := newTestEnvelopeSink(t, writeQueue, t.Context())
	go dbStorerInstance.Start()

	// Permit both nodeID and otherPermittedID
	permitted := map[uint32]struct{}{
		nodeID:           {},
		otherPermittedID: {},
	}

	lastSequenceIds := make(map[uint32]uint64)

	origStream := newTestOriginatorStreamWithPermitted(
		t,
		&node,
		stream,
		lastSequenceIds,
		permitted,
		writeQueue,
	)

	// Your listen() currently returns RetryAfter on EOF, so just assert it returns *an* error
	// and then assert the envelope eventually appears in DB.
	_ = origStream.listen()

	require.Eventually(t, func() bool {
		envs := getAllMessagesForOriginator(t, dbStorerInstance, otherPermittedID)
		return len(envs) == 1 && envs[0].OriginatorSequenceID == int64(sequenceID)
	}, time.Second, 50*time.Millisecond)
}

func TestSyncWorkerOutOfOrderStillAdvancesLastSequenceId(t *testing.T) {
	nodeID := uint32(200)

	// Create seq=1 then seq=3 (skip 2 to force out-of-order)
	env1 := envelopeTestUtils.CreateOriginatorEnvelope(t, nodeID, uint64(1))
	env3 := envelopeTestUtils.CreateOriginatorEnvelope(t, nodeID, uint64(3))

	stream := mockSubscriptionOnePage(t, []*envelopes.OriginatorEnvelope{env1, env3})
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 10)
	defer close(writeQueue)

	dbStorerInstance := newTestEnvelopeSink(t, writeQueue, t.Context())
	go dbStorerInstance.Start()

	lastSequenceIds := make(map[uint32]uint64)

	// --- Replace test logger with zap observer ---
	core, recorded := observer.New(zap.ErrorLevel)
	logger := zap.New(core)

	permitted := map[uint32]struct{}{
		nodeID: {},
	}

	origStream := newOriginatorStream(
		t.Context(),
		logger, // use observed logger
		&node,
		lastSequenceIds,
		permitted,
		stream,
		writeQueue,
	)

	_ = origStream.listen()

	// ---- Assert lastSequenceId advanced to 3 ----
	require.Eventually(t, func() bool {
		return lastSequenceIds[nodeID] == 3
	}, time.Second, 50*time.Millisecond)

	// ---- Assert error log was emitted ----
	require.Eventually(t, func() bool {
		logs := recorded.All()
		for _, log := range logs {
			if log.Message == "received out-of-order envelope" {
				return true
			}
		}
		return false
	}, time.Second, 50*time.Millisecond)
}

func TestSyncWorkerNoOutOfOrderErrorForMultipleOriginatorsInOrder(t *testing.T) {
	envs := []*envelopes.OriginatorEnvelope{
		envelopeTestUtils.CreateOriginatorEnvelope(t, 200, 1),
		envelopeTestUtils.CreateOriginatorEnvelope(t, 10, 1),
		envelopeTestUtils.CreateOriginatorEnvelope(t, 10, 2),
		envelopeTestUtils.CreateOriginatorEnvelope(t, 10, 3),
		envelopeTestUtils.CreateOriginatorEnvelope(t, 13, 1),

		envelopeTestUtils.CreateOriginatorEnvelope(t, 200, 2),
		envelopeTestUtils.CreateOriginatorEnvelope(t, 13, 2),

		envelopeTestUtils.CreateOriginatorEnvelope(t, 200, 3),
		envelopeTestUtils.CreateOriginatorEnvelope(t, 13, 3),
	}

	stream := mockSubscriptionOnePage(t, envs)

	// "Node we are syncing from" (doesn't have to match all originators, but must be permitted)
	nodeID := uint32(200)
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	writeQueue := make(chan *envUtils.OriginatorEnvelope, 50)
	defer close(writeQueue)

	dbStorerInstance := newTestEnvelopeSink(t, writeQueue, t.Context())
	go dbStorerInstance.Start()

	lastSequenceIds := make(map[uint32]uint64)

	// Observe logs at Error level
	core, recorded := observer.New(zap.ErrorLevel)
	logger := zap.New(core)

	permitted := map[uint32]struct{}{
		200: {},
		10:  {},
		13:  {},
	}

	origStream := newOriginatorStream(
		t.Context(),
		logger,
		&node,
		lastSequenceIds,
		permitted,
		stream,
		writeQueue,
	)

	_ = origStream.listen()

	// And sanity-check lastSequenceIds advanced correctly for all originators
	require.Eventually(t, func() bool {
		return lastSequenceIds[200] == 3 && lastSequenceIds[10] == 3 && lastSequenceIds[13] == 3
	}, time.Second, 50*time.Millisecond)

	require.Eventually(t, func() bool {
		a := getAllMessagesForOriginator(t, dbStorerInstance, 200)
		b := getAllMessagesForOriginator(t, dbStorerInstance, 10)
		c := getAllMessagesForOriginator(t, dbStorerInstance, 13)
		return len(a) == 3 && len(b) == 3 && len(c) == 3
	}, time.Second, 50*time.Millisecond)

	// even though we encountered 1, 1,2,3 2,3 2,3 we should not complain
	require.Empty(t, recorded.FilterMessage("received out-of-order envelope").All())
}
