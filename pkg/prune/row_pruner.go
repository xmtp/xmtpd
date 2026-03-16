package prune

import (
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

func (e *Executor) PruneRows() error {
	var (
		querier = queries.New(e.writerDB)
		start   = time.Now()
	)

	if e.config.CountDeletable {
		envelopesCount, err := querier.CountExpiredEnvelopes(e.ctx)
		if err != nil {
			e.logger.Error("could not count expired envelopes", zap.Error(err))
			return err
		}
		e.logger.Info("count of envelopes eligible for pruning", utils.CountField(envelopesCount))

		if envelopesCount == 0 {
			e.logger.Info("no envelopes found for pruning")
			return nil
		}
	}

	if e.config.DryRun {
		e.logger.Info("dry run mode enabled, nothing to do")
		return nil
	}

	var (
		cyclesCompleted    = 0
		totalDeletionCount = int64(0)
	)

	latestEnvelopes, err := querier.SelectVectorClock(e.ctx)
	if err != nil {
		e.logger.Error("error selecting vector clock", zap.Error(err))
		return err
	}

	ceilings, err := querier.GetPrunableCeiling(e.ctx)
	if err != nil {
		e.logger.Error("error getting ceilings", zap.Error(err))
		return err
	}
	deletableCeilings := make(map[int32]int64)
	for _, ceiling := range ceilings {
		deletableCeilings[ceiling.OriginatorNodeID] = ceiling.MaxEndSequenceID
	}

	deletableTables := make(map[string]int64)
	for _, t := range latestEnvelopes {
		if t.OriginatorNodeID == 0 || t.OriginatorNodeID == 1 {
			e.logger.Debug(
				"originator is not prunable in this version of XMTPD. Skipping...",
				utils.OriginatorIDField(uint32(t.OriginatorNodeID)),
			)
			continue
		}

		ceilingForThisOriginator := deletableCeilings[t.OriginatorNodeID]

		if ceilingForThisOriginator == 0 {
			e.logger.Debug(
				"originator is not prunable. No reports exist. Skipping...",
				utils.OriginatorIDField(uint32(t.OriginatorNodeID)),
			)
			continue
		}

		e.logger.Debug(
			"Attempting to prune envelopes for originator",
			utils.OriginatorIDField(uint32(t.OriginatorNodeID)),
		)
		deletableTables[fmt.Sprintf("gateway_envelopes_meta_o%d", t.OriginatorNodeID)] = ceilingForThisOriginator
	}

	for {
		if cyclesCompleted >= e.config.MaxCycles {
			e.logger.Warn(
				"reached maximum pruning cycles",
				zap.Int("max_cycles", e.config.MaxCycles),
			)
			break
		}

		if len(deletableTables) == 0 {
			e.logger.Info("all tables have been processed")
			break
		}

		var deletedThisCycle int64

		for tableName, ceiling := range deletableTables {
			result, err := e.writerDB.Exec(
				constructVariableMetaTableQuery(tableName, e.config.BatchSize, ceiling),
			)
			if err != nil {
				e.logger.Error(
					"error pruning envelopes",
					zap.Error(err),
					zap.String("table", tableName),
				)
				delete(deletableTables, tableName)
				continue
			}
			rows, err := result.RowsAffected()
			if err != nil {
				e.logger.Error(
					"Unexpected DB error: could not count envelopes",
					zap.Error(err),
					zap.String("table", tableName),
				)
				delete(deletableTables, tableName)
				continue
			}

			deletedThisCycle += rows

			if rows < int64(e.config.BatchSize) {
				// this one is fully exhausted
				delete(deletableTables, tableName)
			}

			e.logger.Debug(
				"pruned envelopes",
				zap.Int64("deleted", rows),
				zap.String("table", tableName),
			)
		}

		totalDeletionCount += deletedThisCycle

		e.logger.Info("pruned expired envelopes batch", utils.CountField(deletedThisCycle))

		cyclesCompleted++
	}

	if totalDeletionCount == 0 {
		e.logger.Info("no expired envelopes found")
	}

	e.logger.Info(
		"row pruning done",
		utils.CountField(totalDeletionCount),
		utils.DurationMsField(time.Since(start)),
	)

	return nil
}
