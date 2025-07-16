CREATE TABLE key_packages(
	sequence_id BIGSERIAL PRIMARY KEY,
	installation_id BYTEA NOT NULL,
	key_package BYTEA NOT NULL
);

CREATE TABLE group_messages(
	id BIGSERIAL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	group_id BYTEA NOT NULL,
	data BYTEA NOT NULL,
	group_id_data_hash BYTEA NOT NULL
);

CREATE TABLE welcome_messages(
	id BIGSERIAL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	installation_key BYTEA NOT NULL,
	data BYTEA NOT NULL,
	hpke_public_key BYTEA NOT NULL,
	installation_key_data_hash BYTEA NOT NULL,
	wrapper_algorithm SMALLINT NOT NULL DEFAULT 0
);

CREATE TABLE inbox_log(
	sequence_id BIGSERIAL PRIMARY KEY,
	inbox_id BYTEA NOT NULL,
	server_timestamp_ns BIGINT NOT NULL,
	identity_update_proto BYTEA NOT NULL
);
