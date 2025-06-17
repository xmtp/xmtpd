-- This table includes Payer Reports sent from any node.
CREATE TABLE payer_reports (
    id BYTEA PRIMARY KEY,
    originator_node_id INT NOT NULL,
    start_sequence_id BIGINT NOT NULL,
    end_sequence_id BIGINT NOT NULL,
    end_minute_since_epoch INT NOT NULL,
    payers_merkle_root BYTEA NOT NULL,
    active_node_ids INT [] NOT NULL,
    -- 0 = pending, 1 = submitted, 2 = settled
    submission_status SMALLINT NOT NULL DEFAULT 0,
    -- 0 = pending, 1 = approved, 2 = rejected
    attestation_status SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX payer_reports_submission_status_created_idx ON payer_reports (submission_status, created_at);

CREATE INDEX payer_reports_attestation_status_created_idx ON payer_reports (attestation_status, created_at);

CREATE TABLE payer_report_attestations (
    -- Do not reference the payer reports table since attestations may arrive before the report is stored
    payer_report_id BYTEA NOT NULL,
    node_id BIGINT NOT NULL,
    signature BYTEA NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (payer_report_id, node_id)
);

CREATE INDEX payer_report_attestations_payer_report_id_idx ON payer_report_attestations (payer_report_id);