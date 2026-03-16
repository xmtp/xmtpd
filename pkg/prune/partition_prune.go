package prune

import (
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/constants"

	"github.com/xmtp/xmtpd/pkg/utils"

	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

func constructBlobName(row queries.GetPrunableMetaPartitionsRow) string {
	return fmt.Sprintf(
		"gateway_envelope_blobs_o%d_s%d_%d",
		row.OriginatorNodeID,
		row.BandStart,
		row.BandEnd,
	)
}

func (e *Executor) DropPrunablePartitions() error {
	ctx := e.ctx

	start := time.Now()
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

	for _, droppableMetaRow := range parts {
		if droppableMetaRow.OriginatorNodeID == constants.GroupMessageOriginatorID ||
			droppableMetaRow.OriginatorNodeID == constants.IdentityUpdateOriginatorID {
			e.logger.Info(
				"refusing to drop this partition in this version of XMTPD",
				zap.String("partition", droppableMetaRow.Tablename),
				utils.OriginatorIDField(uint32(droppableMetaRow.OriginatorNodeID)),
			)
			continue
		}

		blobName := constructBlobName(droppableMetaRow)

		if _, err := e.writerDB.ExecContext(
			ctx,
			constructDropQuery(droppableMetaRow.Tablename, blobName),
		); err != nil {
			e.logger.Error("could not drop partition pair", zap.Error(err))
			return fmt.Errorf(
				"drop partition pair (%s, %s): %w",
				droppableMetaRow.Tablename,
				blobName,
				err,
			)
		}

		e.logger.Info(
			"dropped partition pair",
			zap.String("blob_table", blobName),
			zap.String("meta_table", droppableMetaRow.Tablename),
		)
	}

	e.logger.Info(
		"partition deletion done",
		utils.DurationMsField(time.Since(start)),
	)

	return nil
}
