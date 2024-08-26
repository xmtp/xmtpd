package api

import (
	"context"
	"database/sql"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type PublishWorker struct {
	ctx          context.Context
	log          *zap.Logger
	listener     <-chan []queries.StagedOriginatorEnvelope
	notifier     chan<- bool
	registrant   *registrant.Registrant
	store        *sql.DB
	subscription db.DBSubscription[queries.StagedOriginatorEnvelope]
}

func StartPublishWorker(
	ctx context.Context,
	log *zap.Logger,
	reg *registrant.Registrant,
	store *sql.DB,
) (*PublishWorker, error) {
	q := queries.New(store)
	query := func(ctx context.Context, lastSeenID int64, numRows int32) ([]queries.StagedOriginatorEnvelope, int64, error) {
		results, err := q.SelectStagedOriginatorEnvelopes(
			ctx,
			queries.SelectStagedOriginatorEnvelopesParams{
				LastSeenID: lastSeenID,
				NumRows:    numRows,
			},
		)
		if err != nil {
			return nil, 0, err
		}
		if len(results) > 0 {
			lastSeenID = results[len(results)-1].ID
		}
		return results, lastSeenID, nil
	}
	notifier := make(chan bool, 1)
	subscription := db.NewDBSubscription(
		ctx,
		log,
		query,
		0, // lastSeenID
		db.PollingOptions{Interval: time.Second, Notifier: notifier, NumRows: 100},
	)
	listener, err := subscription.Start()
	if err != nil {
		return nil, err
	}

	worker := &PublishWorker{
		ctx:          ctx,
		log:          log,
		notifier:     notifier,
		subscription: *subscription,
		listener:     listener,
		registrant:   reg,
		store:        store,
	}
	go worker.start()

	return worker, nil
}

func (p *PublishWorker) NotifyStagedPublish() {
	select {
	case p.notifier <- true:
	default:
	}
}

func (p *PublishWorker) start() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case new_batch := <-p.listener:
			for _, stagedEnv := range new_batch {
				for !p.publishStagedEnvelope(stagedEnv) {
					// Infinite retry on failure to publish; we cannot
					// continue to the next envelope until this one is processed
					time.Sleep(time.Second)
				}
			}
		}
	}
}

func (p *PublishWorker) publishStagedEnvelope(stagedEnv queries.StagedOriginatorEnvelope) bool {
	logger := p.log.With(zap.Int64("sequenceID", stagedEnv.ID))
	originatorEnv, err := p.registrant.SignStagedEnvelope(stagedEnv)
	if err != nil {
		logger.Error(
			"Failed to sign staged envelope",
			zap.Error(err),
		)
		return false
	}
	originatorBytes, err := proto.Marshal(originatorEnv)
	if err != nil {
		logger.Error("Failed to marshal originator envelope", zap.Error(err))
		return false
	}

	q := queries.New(p.store)

	// On unique constraint conflicts, no error is thrown, but numRows is 0
	inserted, err := q.InsertGatewayEnvelope(
		p.ctx,
		queries.InsertGatewayEnvelopeParams{
			OriginatorID:         int32(p.registrant.NodeID()),
			OriginatorSequenceID: stagedEnv.ID,
			Topic:                stagedEnv.Topic,
			OriginatorEnvelope:   originatorBytes,
		},
	)
	if err != nil {
		logger.Error("Failed to insert gateway envelope", zap.Error(err))
		return false
	} else if inserted == 0 {
		// Envelope was already inserted by another worker
		logger.Debug("Envelope already inserted")
	}

	// Try to delete the row regardless of if the gateway envelope was inserted elsewhere
	deleted, err := q.DeleteStagedOriginatorEnvelope(context.Background(), stagedEnv.ID)
	if err != nil {
		logger.Error("Failed to delete staged envelope", zap.Error(err))
		// Envelope is already inserted, so it is safe to continue
		return true
	} else if deleted == 0 {
		// Envelope was already deleted by another worker
		logger.Debug("Envelope already deleted")
	}

	return true
}
