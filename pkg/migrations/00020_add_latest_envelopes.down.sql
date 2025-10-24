DROP TRIGGER IF EXISTS gateway_latest_upd ON gateway_envelopes;
DROP FUNCTION IF EXISTS update_latest_envelope();
DROP TABLE IF EXISTS gateway_envelopes_latest;