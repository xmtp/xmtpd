ALTER TABLE unsettled_usage DROP CONSTRAINT IF EXISTS fk_unsettled_usage_payer_id;
ALTER TABLE payer_ledger_events DROP CONSTRAINT IF EXISTS fk_payer_ledger_events_payer_id;
