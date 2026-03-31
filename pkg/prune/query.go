package prune

import (
	"fmt"

	"github.com/lib/pq"
)

func constructVariableMetaTableQuery(
	tableName string,
	batchSize int32,
	maxSequenceId int64,
) string {
	return fmt.Sprintf(`
WITH to_delete AS (
  SELECT ctid
  FROM %s
  WHERE expiry < EXTRACT(EPOCH FROM now())::bigint
    AND originator_sequence_id < %d
  ORDER BY expiry
  LIMIT %d
  FOR UPDATE SKIP LOCKED
)
DELETE FROM %s
WHERE ctid IN (SELECT ctid FROM to_delete);
`, pq.QuoteIdentifier(tableName), maxSequenceId, batchSize, pq.QuoteIdentifier(tableName))
}

func constructDropQuery(metaTable string, blobTable string) string {
	return fmt.Sprintf(
		"DROP TABLE IF EXISTS %s,%s CASCADE",
		pq.QuoteIdentifier(metaTable),
		pq.QuoteIdentifier(blobTable),
	)
}
