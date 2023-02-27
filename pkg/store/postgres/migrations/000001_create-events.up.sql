BEGIN;

CREATE TABLE IF NOT EXISTS events(
   cid VARCHAR(100) PRIMARY KEY,
   links JSON NOT NULL,
   topic VARCHAR (300) NOT NULL,
   timestamp_ns BIGINT NOT NULL,
   message BYTEA NOT NULL 
);

CREATE UNIQUE INDEX events_topic_ts_cid_idx ON events (topic, timestamp_ns, cid);

COMMIT;
