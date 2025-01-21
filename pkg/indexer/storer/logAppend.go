package storer

import (
	"context"
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"go.uber.org/zap"
)

func GetVersionForAppend(
	ctx context.Context,
	querier *queries.Queries,
	logger *zap.Logger,
	originatorNodeID int32,
	sequenceID int64,
) (sql.NullInt32, error) {
	var version sql.NullInt32

	return version, nil
}
