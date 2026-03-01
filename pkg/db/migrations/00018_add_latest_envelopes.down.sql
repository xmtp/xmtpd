-- First reinstate previous version of the trigger, as defined in 00006_add_latest_envelopes.up.sql
CREATE TRIGGER gateway_latest_upd
    AFTER INSERT ON gateway_envelopes_meta
    FOR EACH ROW EXECUTE FUNCTION update_latest_envelope();


DROP TRIGGER IF EXISTS gateway_latest_upd_v2 ON gateway_envelopes_meta;
DROP FUNCTION IF EXISTS update_latest_envelope_v2();

