-- Ensures that if the command-line node configuration mutates,
-- the existing data in the DB is invalid
CREATE TABLE node_info(
	node_id INTEGER NOT NULL,
	public_key BYTEA NOT NULL,
	singleton_id SMALLINT PRIMARY KEY DEFAULT 1,
	CONSTRAINT is_singleton CHECK (singleton_id = 1)
);

-- Newly published envelopes will be queued here first (and assigned an originator
-- sequence ID), before being inserted in-order into the gateway_envelopes table.
CREATE TABLE staged_originator_envelopes(
	-- used to construct originator_sid
	id BIGSERIAL PRIMARY KEY,
	originator_time TIMESTAMP NOT NULL DEFAULT now(),
	topic BYTEA NOT NULL,
	payer_envelope BYTEA NOT NULL
);

CREATE FUNCTION insert_staged_originator_envelope(topic BYTEA, payer_envelope BYTEA)
	RETURNS SETOF staged_originator_envelopes
	AS $$
BEGIN
	-- Ensures that the generated sequence ID matches the insertion order
	-- Only released at the end of the enclosing transaction - beware if called within a long transaction
	PERFORM
		pg_advisory_xact_lock(hashtext('staged_originator_envelopes_sequence'));
	RETURN QUERY INSERT INTO staged_originator_envelopes(topic, payer_envelope)
		VALUES(topic, payer_envelope)
	ON CONFLICT
		DO NOTHING
	RETURNING
		*;
END;
$$
LANGUAGE plpgsql;

-- A cached view for looking up the inbox_id that an address belongs to.
-- Relies on a total ordering of updates across all inbox_ids, from which this
-- view can be deterministically generated.
CREATE TABLE address_log(
	address TEXT NOT NULL,
	inbox_id BYTEA NOT NULL,
	association_sequence_id BIGINT,
	revocation_sequence_id BIGINT,
	PRIMARY KEY (address, inbox_id)
);

CREATE TABLE payers(
	id SERIAL PRIMARY KEY,
	address TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS gateway_envelopes_meta
(
    gateway_time           timestamp NOT NULL DEFAULT now(),
    originator_node_id     int       NOT NULL,
    originator_sequence_id bigint    NOT NULL,
    topic                  bytea     NOT NULL,
    -- Leave column nullable since blockchain originated messages won't have a payer_id
    payer_id               int REFERENCES payers(id),
    expiry                 bigint NOT NULL,
    PRIMARY KEY (originator_node_id, originator_sequence_id)
) PARTITION BY LIST (originator_node_id);

-- BLOBS (cold path)
CREATE TABLE IF NOT EXISTS gateway_envelope_blobs
(
    originator_node_id     int    NOT NULL,
    originator_sequence_id bigint NOT NULL,
    originator_envelope    bytea  NOT NULL,
    PRIMARY KEY (originator_node_id, originator_sequence_id),
    FOREIGN KEY (originator_node_id, originator_sequence_id) REFERENCES gateway_envelopes_meta(originator_node_id, originator_sequence_id) ON DELETE CASCADE
) PARTITION BY LIST (originator_node_id);

CREATE OR REPLACE VIEW gateway_envelopes_view AS
SELECT
    m.originator_node_id,
    m.originator_sequence_id,
    m.gateway_time,
    m.topic,
    b.originator_envelope
FROM gateway_envelopes_meta m
         JOIN gateway_envelope_blobs b
              ON b.originator_node_id     = m.originator_node_id
                  AND b.originator_sequence_id = m.originator_sequence_id;