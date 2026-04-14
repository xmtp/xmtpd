-- Revert migration 00024 by dropping the hardened partition helpers. The
-- pre-existing v2/v3 helpers are untouched by the up migration, so no rename
-- or restore is required here.

DROP FUNCTION IF EXISTS ensure_gateway_parts_v4(int, bigint, bigint);
DROP FUNCTION IF EXISTS make_blob_seq_subpart_v4(int, bigint, bigint);
DROP FUNCTION IF EXISTS make_blob_originator_part_v4(int);
DROP FUNCTION IF EXISTS make_meta_seq_subpart_v3(int, bigint, bigint);
DROP FUNCTION IF EXISTS make_meta_originator_part_v3(int);
