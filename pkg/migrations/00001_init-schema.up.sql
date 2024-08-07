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
	originator_sid BIGINT NOT NULL,
	topic BYTEA NOT NULL,
	originator_envelope BYTEA NOT NULL
);
-- Client queries
CREATE INDEX idx_gateway_envelopes_topic ON gateway_envelopes(topic);
-- Node queries
CREATE UNIQUE INDEX idx_gateway_envelopes_originator_sid ON gateway_envelopes(originator_sid);


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
