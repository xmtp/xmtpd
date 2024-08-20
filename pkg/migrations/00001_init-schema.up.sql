-- Ensures that if the command-line node configuration mutates,
-- the existing data in the DB is invalid
CREATE TABLE node_info(
	node_id INTEGER NOT NULL,
	public_key BYTEA NOT NULL,
	singleton_id SMALLINT PRIMARY KEY DEFAULT 1,
	CONSTRAINT is_singleton CHECK (singleton_id = 1)
);

-- Includes all envelopes, whether they were originated locally or not
CREATE TABLE gateway_envelopes(
	-- used to construct gateway_sid
	id BIGSERIAL PRIMARY KEY,
	originator_node_id INT NOT NULL,
	originator_sequence_id BIGINT NOT NULL,
	topic BYTEA NOT NULL,
	originator_envelope BYTEA NOT NULL
);

-- Client queries
CREATE INDEX idx_gateway_envelopes_topic ON gateway_envelopes(topic);

-- Node queries
CREATE UNIQUE INDEX idx_gateway_envelopes_originator_sid ON gateway_envelopes(originator_node_id, originator_sequence_id);

CREATE FUNCTION insert_gateway_envelope(originator_node_id INT, originator_sequence_id BIGINT, topic BYTEA, originator_envelope BYTEA)
	RETURNS SETOF gateway_envelopes
	AS $$
BEGIN
	-- Ensures that the generated sequence ID matches the insertion order
	-- Only released at the end of the enclosing transaction - beware if called within a long transaction
	PERFORM
		pg_advisory_xact_lock(hashtext('gateway_envelopes_sequence'));
	RETURN QUERY INSERT INTO gateway_envelopes(originator_node_id, originator_sequence_id, topic, originator_envelope)
		VALUES(originator_node_id, originator_sequence_id, topic, originator_envelope)
	ON CONFLICT
		DO NOTHING
	RETURNING
		*;
END;
$$
LANGUAGE plpgsql;

-- Process for originating envelopes:
-- 1. Perform any necessary validation
-- 2. Insert into originated_envelopes
-- 3. Singleton background task will continuously query (or subscribe to)
--    staged_originated_envelopes, and for each envelope in order of ID:
--     2.1. Construct and sign OriginatorEnvelope proto
--     2.2. Atomically insert into all_envelopes and delete from originated_envelopes,
--	        ignoring unique index violations on originator_sid
-- This preserves total ordering, while avoiding gaps in sequence ID's.
CREATE TABLE staged_originator_envelopes(
	-- used to construct originator_sid
	id BIGSERIAL PRIMARY KEY,
	originator_time TIMESTAMP NOT NULL DEFAULT now(),
	payer_envelope BYTEA NOT NULL
);

CREATE FUNCTION insert_staged_originator_envelope(payer_envelope BYTEA)
	RETURNS SETOF staged_originator_envelopes
	AS $$
BEGIN
	PERFORM
		pg_advisory_xact_lock(hashtext('staged_originator_envelopes_sequence'));
	RETURN QUERY INSERT INTO staged_originator_envelopes(payer_envelope)
		VALUES(payer_envelope)
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

