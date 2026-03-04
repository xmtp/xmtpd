ALTER TABLE unsettled_usage 
    ADD CONSTRAINT fk_unsettled_usage_payer_id 
    FOREIGN KEY (payer_id) REFERENCES payers(id);

ALTER TABLE payer_ledger_events 
    ADD CONSTRAINT fk_payer_ledger_events_payer_id 
    FOREIGN KEY (payer_id) REFERENCES payers(id);
