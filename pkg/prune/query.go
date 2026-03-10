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

WITH to_delete AS (SELECT ge.originator_node_id,
                          ge.originator_sequence_id
                   FROM %s ge
                   WHERE ge.expiry < EXTRACT(EPOCH FROM now())::bigint
                     AND ge.originator_sequence_id <= %d
                   ORDER BY ge.expiry, ge.originator_node_id, ge.originator_sequence_id
                   LIMIT %d FOR UPDATE SKIP LOCKED)
DELETE
FROM %s ge
    USING to_delete td
WHERE ge.originator_node_id = td.originator_node_id
  AND ge.originator_sequence_id = td.originator_sequence_id
RETURNING ge.originator_node_id, ge.originator_sequence_id, ge.expiry;
`, pq.QuoteIdentifier(tableName), maxSequenceId, batchSize, pq.QuoteIdentifier(tableName))
}
