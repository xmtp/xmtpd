package worker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/utils"
)

const (
	DefaultPartitionSize = 1_000_000
	DefaultFillThreshold = 0.7
	DefaultCheckInterval = 30 * time.Minute
)

var defaultConfig = Config{
	Interval: DefaultCheckInterval,
	Partition: PartitionConfig{
		PartitionSize: DefaultPartitionSize,
		FillThreshold: DefaultFillThreshold,
	},
}

type Config struct {
	Interval  time.Duration
	Partition PartitionConfig
}

// partition config controls when should the database worker create new partitions.
// given the partition size, it will create new partitions when the the partition
// is fill more than the specified fill threshold.
// For example: partition size is 1000, fill threshold is 70% => when the partition
// has over 700 entries it will create the next partition.
// TODO: Perhaps use a list of originators and setup partitions even for nodes that do not have any yet.

type PartitionConfig struct {
	PartitionSize uint64
	FillThreshold float64
}

type Worker struct {
	cfg Config

	log *zap.Logger
	db  *db.Handler
}

func NewWorker(log *zap.Logger, db *db.Handler) *Worker {
	return newWorkerWithConfig(defaultConfig, log, db)
}

func newWorkerWithConfig(cfg Config, log *zap.Logger, db *db.Handler) *Worker {
	worker := &Worker{
		cfg: cfg,
		log: log.Named(utils.DatabaseWorkerLoggerName),
		db:  db,
	}

	return worker
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
	partitions, err := w.getPartitionList(ctx)
	if err != nil {
		return fmt.Errorf("could not retrieve list of database partitions: %w", err)
	}

	if len(partitions) == 0 {
		return errors.New("could not identify any partition tables")
	}

	np := sortPartitions(partitions)
	err = np.validate()
	if err != nil {
		return fmt.Errorf("invalid partition chain(s) found: %w", err)
	}

	// NOTE: We do not validate that partitions are the width that we expect
	// from the worker; we will only enforce that for newly created partitions

	var errs []error

	for nodeID, partitions := range np.partitions {
		w.log.Debug("processing partitions for node",
			zap.Uint32("node_id", nodeID))

		// Should not happen.
		if len(partitions) == 0 {
			continue
		}

		// Only check the last partition as it is the one we will be extending.
		last := partitions[len(partitions)-1]
		count, err := w.getLastSequenceID(ctx, last.name)
		if err != nil {
			// Try to process as many as possible - return the error at the end
			errs = append(errs, fmt.Errorf("could not get last sequence ID for table: %v: %w",
				last.name,
				err),
			)
			continue
		}

		// Use the value we have from the partition name to determine fill ratio.
		fillRatio := float64(count) / float64(last.end)

		w.log.Info("partition fill ratio",
			zap.String("name", last.name),
			zap.Float64("fill", fillRatio))

		// Partition has enough room left, continue.
		if fillRatio <= w.cfg.Partition.FillThreshold {
			continue
		}

		err = w.createPartition(ctx, nodeID, last.end)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not create partition for table: %v: %w",
				last.name,
				err))
		}
	}

	return errors.Join(errs...)
}

func (w *Worker) query() Querier {
	return New(w.db.DB())
}

func (w *Worker) getPartitionList(ctx context.Context) ([]partitionTableInfo, error) {
	tables, err := w.query().ListPartitions(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve list of database tables: %w", err)
	}

	var partitions []partitionTableInfo
	for _, table := range tables {

		info, err := parsePartitionInfo(table)
		if err != nil {
			w.log.Warn("could not parse partition info", zap.String("table", table), zap.Error(err))
			continue
		}

		partitions = append(partitions, info)
	}

	return partitions, nil
}

func (w *Worker) getLastSequenceID(ctx context.Context, table string) (int64, error) {
	query := fmt.Sprintf("SELECT MAX(originator_sequence_id) FROM %s", table)

	var count sql.NullInt64
	err := w.db.DB().QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("could not execute query: %w", err)
	}

	return count.Int64, nil
}

func (w *Worker) createPartition(ctx context.Context, nodeID uint32, sequenceID uint64) error {
	params := queries.EnsureGatewayPartsParams{
		OriginatorNodeID:     int32(nodeID),
		OriginatorSequenceID: int64(sequenceID),
		BandWidth:            int64(w.cfg.Partition.PartitionSize),
	}

	err := w.db.WriteQuery().EnsureGatewayParts(ctx, params)
	if err != nil {
		return fmt.Errorf(
			"could not create partition for node (id: %d, sequence_id: %d, size: %d): %w",
			nodeID,
			sequenceID,
			w.cfg.Partition.PartitionSize,
			err,
		)
	}

	return nil
}
