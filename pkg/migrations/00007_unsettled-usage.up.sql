CREATE TABLE unsettled_usage(
	payer_id INTEGER NOT NULL,
	originator_id INTEGER NOT NULL,
	minutes_since_epoch INTEGER NOT NULL,
	spend_picodollars BIGINT NOT NULL,
	PRIMARY KEY (payer_id, originator_id, minutes_since_epoch)
);

CREATE INDEX idx_unsettled_usage_payer_id ON unsettled_usage(payer_id);