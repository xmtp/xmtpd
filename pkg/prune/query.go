package prune

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// metaTableToBlobTable derives the blob partition table name from a meta partition
// table name by replacing the "gateway_envelopes_meta" prefix with
// "gateway_envelope_blobs".
//
// Example: "gateway_envelopes_meta_o100" → "gateway_envelope_blobs_o100"
func metaTableToBlobTable(metaTable string) string {
	return strings.Replace(metaTable, "gateway_envelopes_meta", "gateway_envelope_blobs", 1)
}

// constructVariableMetaTableQuery builds a CTE query that:
//  1. Selects up to batchSize expired meta rows (SKIP LOCKED).
//  2. Deletes those meta rows.
//  3. Deletes the corresponding blob rows in the same transaction.
//
// It returns the number of deleted meta rows so the caller can detect when a
// partition has been fully exhausted.
func constructVariableMetaTableQuery(
	tableName string,
	batchSize int32,
	maxSequenceId int64,
) string {
	blobTable := metaTableToBlobTable(tableName)
	return fmt.Sprintf(`
WITH to_delete AS (
  SELECT ctid, originator_sequence_id
  FROM %s
  WHERE expiry < EXTRACT(EPOCH FROM now())::bigint
    AND originator_sequence_id <= %d
  ORDER BY expiry
  LIMIT %d
  FOR UPDATE SKIP LOCKED
),
deleted_meta AS (
  DELETE FROM %s
  WHERE ctid IN (SELECT ctid FROM to_delete)
  RETURNING originator_sequence_id
),
deleted_blobs AS (
  DELETE FROM %s
  WHERE originator_sequence_id IN (SELECT originator_sequence_id FROM deleted_meta)
)
SELECT count(*) FROM deleted_meta;
`, pq.QuoteIdentifier(tableName), maxSequenceId, batchSize,
		pq.QuoteIdentifier(tableName), pq.QuoteIdentifier(blobTable))
}

func constructDropQuery(metaTable string, blobTable string) string {
	return fmt.Sprintf(
		"DROP TABLE IF EXISTS %s,%s CASCADE",
		pq.QuoteIdentifier(metaTable),
		pq.QuoteIdentifier(blobTable),
	)
}
