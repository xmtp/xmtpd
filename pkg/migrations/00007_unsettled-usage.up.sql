CREATE TABLE unsettled_usage(
                                payer_id INTEGER NOT NULL,
                                originator_id INTEGER NOT NULL,
                                minutes_since_epoch INTEGER NOT NULL,
                                spend_picodollars BIGINT NOT NULL,
                                last_sequence_id BIGINT NOT NULL,
                                message_count INTEGER NOT NULL DEFAULT 0,
                                PRIMARY KEY (payer_id, originator_id, minutes_since_epoch)
);

CREATE INDEX idx_unsettled_usage_originator_id_minutes_since_epoch
    ON unsettled_usage(originator_id, minutes_since_epoch DESC);