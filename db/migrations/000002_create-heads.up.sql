BEGIN;

CREATE TABLE IF NOT EXISTS heads(
   cid VARCHAR(100) NOT NULL PRIMARY KEY REFERENCES events (cid),
   topic VARCHAR(300) NOT NULL
);

CREATE UNIQUE INDEX heads_topic_cid_idx ON heads (topic, cid);

COMMIT;
