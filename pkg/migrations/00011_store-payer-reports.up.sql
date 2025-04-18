CREATE TABLE payer_reports (
    id BYTEA PRIMARY KEY,
    originator_node_id INT NOT NULL,
    start_sequence_id BIGINT NOT NULL,
    end_sequence_id BIGINT NOT NULL,
    payers_merkle_root BYTEA NOT NULL,
    payers_leaf_count BIGINT NOT NULL,
    nodes_hash BYTEA NOT NULL,
    nodes_count INT NOT NULL,
    -- 0 = pending, 1 = submitted, 2 = settled
    submission_status SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX ON payer_reports (submission_status);

CREATE TABLE payer_report_attestations (
    -- Do not reference the payer reports table since attestations may arrive before the report is stored
    payer_report_id BYTEA NOT NULL,
    node_id BIGINT NOT NULL,
    signature BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (payer_report_id, node_id)
);

CREATE INDEX ON payer_report_attestations (payer_report_id);