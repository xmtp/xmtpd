CREATE TABLE payer_ledger_events(
    event_id BYTEA PRIMARY KEY,
    payer_id INTEGER NOT NULL,
    amount_picodollars BIGINT NOT NULL,
    -- 0: deposit, 1: withdrawal, 2: settlement, 3: canceled withdrawal, 4: reorg reversal
    event_type SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payer_ledger_events_payer_id ON payer_ledger_events(payer_id);