CREATE TABLE originator_congestion(
	originator_id INTEGER NOT NULL,
	num_messages INTEGER NOT NULL DEFAULT 0,
	minutes_since_epoch INTEGER NOT NULL,
	PRIMARY KEY (originator_id, minutes_since_epoch)
);

