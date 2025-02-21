ALTER TABLE gateway_envelopes
-- Leave column nullable since blockchain originated messages won't have a payer_id
	ADD COLUMN payer_id INT REFERENCES payers(id);

