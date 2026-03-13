package prune

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

func blobPartitionNameFromMeta(metaTable string) (string, error) {
	const metaPrefix = "gateway_envelopes_meta_"
	const blobPrefix = "gateway_envelope_blobs_"

	if !strings.HasPrefix(metaTable, metaPrefix) {
		return "", fmt.Errorf("unexpected meta partition name: %s", metaTable)
	}

	return blobPrefix + strings.TrimPrefix(metaTable, metaPrefix), nil
}

func (e *Executor) DropPrunablePartitions() error {
	ctx := e.ctx

	q := queries.New(e.writerDB)

	parts, err := q.GetPrunableMetaPartitions(ctx)
	if err != nil {
		e.logger.Error("get prunable meta partitions", zap.Error(err))
		return fmt.Errorf("get prunable meta partitions: %w", err)
	}

	if len(parts) == 0 {
		e.logger.Info("no prunable partitions found")
		return nil
	}

	for _, p := range parts {
		e.logger.Info("partition is empty and droppable", zap.String("partition", p.Tablename))
	}

	if e.config.DryRun {
		return nil
	}

	for _, p := range parts {
		blobName, err := blobPartitionNameFromMeta(p.Tablename)
		if err != nil {
			e.logger.Error("derive blob partition from meta", zap.Error(err))
			return fmt.Errorf("derive blob partition from meta %s: %w", p.Tablename, err)
		}

		if _, err := e.writerDB.ExecContext(
			ctx,
			constructDropQuery(p.Tablename, blobName),
		); err != nil {
			e.logger.Error("could not drop partition pair", zap.Error(err))
			return fmt.Errorf("drop partition pair (%s, %s): %w", p.Tablename, blobName, err)
		}

		e.logger.Info(
			"dropped partition pair",
			zap.String("blob_table", blobName),
			zap.String("meta_table", p.Tablename),
		)
	}

	return nil
}
