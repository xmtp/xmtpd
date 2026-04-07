-- name: EnsureGatewayParts :exec
-- Pre-rename partition ensure. Calls ensure_gateway_parts_v2, which still
-- references the legacy gateway_envelope_blobs table. Surviving only for
-- migration-behavior tests; production code uses EnsureGatewayPartsV3.
SELECT ensure_gateway_parts_v2(
               @originator_node_id,
               @originator_sequence_id,
               @band_width
       );

-- name: EnsureGatewayPartsV3 :exec
-- Production partition ensure. Calls ensure_gateway_parts_v3, which targets
-- the renamed gateway_envelopes_blobs table.
SELECT ensure_gateway_parts_v3(
               @originator_node_id,
               @originator_sequence_id,
               @band_width
       );

-- name: MakeMetaOriginatorPart :exec
SELECT make_meta_originator_part_v2(@originator_node_id);

-- name: MakeBlobOriginatorPart :exec
-- Pre-rename L1 blob partition maker. Kept for migration-behavior tests.
SELECT make_blob_originator_part_v2(@originator_node_id);

-- name: MakeBlobOriginatorPartV3 :exec
SELECT make_blob_originator_part_v3(@originator_node_id);

-- name: MakeMetaSeqBand :exec
SELECT make_meta_seq_subpart_v2(@originator_node_id, @band_start, @band_end);

-- name: MakeBlobSeqBand :exec
-- Pre-rename L2 blob subpartition maker. Kept for migration-behavior tests.
SELECT make_blob_seq_subpart_v2(@originator_node_id, @band_start, @band_end);

-- name: MakeBlobSeqBandV3 :exec
SELECT make_blob_seq_subpart_v3(@originator_node_id, @band_start, @band_end);

-- name: InsertSavePoint :exec
SAVEPOINT sp_part;

-- name: InsertSavePointRelease :exec
RELEASE SAVEPOINT sp_part;

-- name: InsertSavePointRollback :exec
ROLLBACK TO SAVEPOINT sp_part;
