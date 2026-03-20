-- Orphan blob rows that were deleted cannot be restored.
-- Drop the FK constraint if it was added by this migration.
ALTER TABLE gateway_envelope_blobs
    DROP CONSTRAINT IF EXISTS gateway_envelope_blobs_meta_fk;
