package worker

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	DefaultFillThreshold = 0.7
	DefaultCheckInterval = 30 * time.Minute
)

var defaultConfig = Config{
	Interval: DefaultCheckInterval,
	Partition: PartitionConfig{
		FillThreshold: DefaultFillThreshold,
	},
}

type Config struct {
	Interval  time.Duration
	Partition PartitionConfig
}

// partition config controls when should the database worker create new partitions.
// TODO: Perhaps use a lis of originators and setup partitions even for nodes that do not have any yet.

type PartitionConfig struct {
	FillThreshold float64
}

type Worker struct {
	cfg Config

	lock *sync.RWMutex

	log *zap.Logger
	db  *db.Handler

	createdPartitions map[partitionHeader]struct{}
}

func NewWorker(log *zap.Logger, db *db.Handler) *Worker {
	return newWorkerWithConfig(defaultConfig, log, db)
}

func newWorkerWithConfig(cfg Config, log *zap.Logger, db *db.Handler) *Worker {
	worker := &Worker{
		cfg:               cfg,
		log:               log.Named(utils.DatabaseWorkerLoggerName),
		db:                db,
		lock:              &sync.RWMutex{},
		createdPartitions: make(map[partitionHeader]struct{}),
	}

	return worker
}

type partitionHeader struct {
	nodeID     uint32
	startSeqID uint64
}

func (w *Worker) Start(ctx context.Context) error {
	err := w.runDBCheck(ctx)
	if err != nil {
		w.log.Error("database check failed", zap.Error(err))
		// Not stopping on this error.
	}

	go func() {
		ticker := time.NewTicker(w.cfg.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				w.log.Info("context done, stopping")
				return

			case <-ticker.C:
				w.log.Debug("running db check")
				err := w.runDBCheck(ctx)
				if err != nil {
					w.log.Error("database check failed", zap.Error(err))
				}

				// On error - do nothing.
			}
		}
	}()

	return nil
}

func (w *Worker) runDBCheck(ctx context.Context) error {
	le, err := w.db.ReadQuery().SelectVectorClock(ctx)
	if err != nil {
		return fmt.Errorf("could not retrieve vector clock: %w", err)
	}

	vc := db.ToVectorClock(le)

	var errs []error

	for nodeID, seqID := range vc {

		// NOTE: We're doing uint64 => int64 conversion and arithmetic, so let's be pedantic.
		if seqID > math.MaxInt64 {
			w.log.Warn("sequence ID value larger than int64 range", zap.Uint64("value", seqID))
			continue
		}

		err = w.runPartitionCheck(ctx, nodeID, int64(seqID))
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf("partition check failed for node (id: %v): %w", nodeID, err),
			)
		}
	}

	return errors.Join(errs...)
}

func (w *Worker) runPartitionCheck(ctx context.Context, nodeID uint32, seqID int64) error {
	partitionSize := db.GatewayEnvelopeBandWidth

	fillRatio := float64(seqID%partitionSize) / float64(partitionSize)

	w.log.Debug("partition fill ratio",
		utils.OriginatorIDField(nodeID),
		utils.SequenceIDField(seqID),
		zap.String("filled_percentage", fmt.Sprintf("%.2f", fillRatio)))

	// Partition has enough room left, continue.
	if fillRatio <= w.cfg.Partition.FillThreshold {
		return nil
	}

	targetSeqID := seqID + partitionSize
	if targetSeqID < 0 {
		return fmt.Errorf("sequence ID overflow (value: %v)", seqID)
	}

	startSeqID := uint64((targetSeqID / partitionSize) * partitionSize)

	w.lock.RLock()
	_, exists := w.createdPartitions[partitionHeader{nodeID: nodeID, startSeqID: startSeqID}]
	w.lock.RUnlock()

	if exists {
		w.log.Info("partition already created, skipping",
			utils.OriginatorIDField(nodeID),
			zap.Uint64("start_sequence_id", startSeqID),
		)
		return nil
	}

	params := queries.EnsureGatewayPartsParams{
		OriginatorNodeID:     int32(nodeID),
		OriginatorSequenceID: targetSeqID,
		BandWidth:            partitionSize,
	}
	err := w.db.WriteQuery().EnsureGatewayParts(ctx, params)
	if err != nil {
		return fmt.Errorf("could not create gateway partitions: %w", err)
	}

	w.log.Info("created partition for node",
		utils.OriginatorIDField(nodeID),
		zap.Uint64("start_sequence_id", startSeqID),
	)

	w.lock.Lock()
	defer w.lock.Unlock()
	w.createdPartitions[partitionHeader{nodeID: nodeID, startSeqID: startSeqID}] = struct{}{}

	return nil
}
