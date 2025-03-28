package sync

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"testing"

	"github.com/cenkalti/backoff/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	messageApiMocks "github.com/xmtp/xmtpd/pkg/mocks/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	feesTestUtils "github.com/xmtp/xmtpd/pkg/testutils/fees"
	registryTestUtils "github.com/xmtp/xmtpd/pkg/testutils/registry"
)

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

func newMinimalSyncWorker(t *testing.T) *syncWorker {
	ctx := context.Background()
	log := testutils.NewLog(t)
	calculator := feesTestUtils.NewTestFeeCalculator()
	db, _, dbCleanup := testutils.NewDB(t, ctx)
	t.Cleanup(dbCleanup)

	return &syncWorker{
		ctx:           ctx,
		log:           log,
		store:         db,
		feeCalculator: calculator,
	}
}

func createBrokenDB(t *testing.T) *sql.DB {
	config, err := pgxpool.ParseConfig("postgres://foo:5432")
	require.NoError(t, err)

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	require.NoError(t, err)

	return stdlib.OpenDBFromPool(db)
}

func TestSyncWorkerSuccess(t *testing.T) {
	worker := newMinimalSyncWorker(t)
	nodeID := uint32(200)
	sequenceID := uint64(100)
	envelope := envelopeTestUtils.CreateOriginatorEnvelope(t, nodeID, sequenceID)

	stream := mockSubscriptionOnePage(t, []*envelopes.OriginatorEnvelope{envelope})
	origStream := &originatorStream{
		nodeID:              nodeID,
		stream:              stream,
		lastEnvelope:        nil,
		messageRetryBackoff: backoff.NewExponentialBackOff(),
	}
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	err := worker.listenToStream(context.Background(), node, origStream)
	var retryAfter *backoff.RetryAfterError
	require.ErrorAs(t, err, &retryAfter)
	require.Equal(t, retryAfter.Duration.Seconds(), float64(1))
}

func TestSyncWorkerPermanentError(t *testing.T) {
	worker := newMinimalSyncWorker(t)
	nodeID := uint32(200)
	sequenceID := uint64(100)
	envelope := envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(
		t,
		nodeID,
		sequenceID,
		[]byte("broken"),
	)

	stream := mockSubscriptionOnePage(t, []*envelopes.OriginatorEnvelope{envelope})
	origStream := &originatorStream{
		nodeID:              nodeID,
		stream:              stream,
		lastEnvelope:        nil,
		messageRetryBackoff: backoff.NewExponentialBackOff(),
	}
	node := registryTestUtils.CreateNode(nodeID, 999, testutils.RandomPrivateKey(t))

	err := worker.listenToStream(context.Background(), node, origStream)
	var retryAfter *backoff.RetryAfterError
	require.ErrorAs(t, err, &retryAfter)
}

func TestSyncWorkerRetryableError(t *testing.T) {
	worker := newMinimalSyncWorker(t)
	nodeID := uint32(200)
	sequenceID := uint64(100)
	envelope := envelopeTestUtils.CreateOriginatorEnvelope(t, nodeID, sequenceID)
	origStream := &originatorStream{
		nodeID:              nodeID,
		stream:              nil,
		lastEnvelope:        nil,
		messageRetryBackoff: backoff.NewExponentialBackOff(),
	}
	// Create a totally broken DB connection and replace the one in the store
	worker.store = createBrokenDB(t)

	err := worker.validateAndInsertEnvelope(origStream, envelope)
	// This should be a retryable error, and not permanent
	var permanent *backoff.PermanentError
	require.False(t, errors.As(err, &permanent))
}
