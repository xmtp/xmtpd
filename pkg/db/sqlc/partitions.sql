-- name: EnsureGatewayParts :exec
SELECT ensure_gateway_parts_v2(
               @originator_node_id,
               @originator_sequence_id,
               @band_width
       );

-- name: MakeMetaOriginatorPart :exec
SELECT make_meta_originator_part_v2(@originator_node_id);

-- name: MakeBlobOriginatorPart :exec
SELECT make_blob_originator_part_v2(@originator_node_id);

-- name: MakeMetaSeqBand :exec
SELECT make_meta_seq_subpart_v2(@originator_node_id, @band_start, @band_end);

-- name: MakeBlobSeqBand :exec
SELECT make_blob_seq_subpart_v2(@originator_node_id, @band_start, @band_end);

-- name: InsertSavePoint :exec
SAVEPOINT sp_part;

-- name: InsertSavePointRelease :exec
RELEASE SAVEPOINT sp_part;

-- name: InsertSavePointRollback :exec
ROLLBACK TO SAVEPOINT sp_part;