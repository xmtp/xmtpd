-- Create the sequence that starts at 0
CREATE SEQUENCE sequence_id_seq MINVALUE 0 START WITH 0 INCREMENT BY 1 NO CYCLE;

-- Create the table using the sequence
CREATE TABLE payer_sequences (
    id BIGINT PRIMARY KEY DEFAULT nextval('sequence_id_seq'),
    available BOOLEAN DEFAULT TRUE
);
