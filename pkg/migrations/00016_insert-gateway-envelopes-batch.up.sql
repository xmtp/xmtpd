CREATE TYPE gateway_envelope_row AS (
  originator_node_id     int,
  originator_sequence_id bigint,
  topic                  bytea,
  payer_id               int,
  gateway_time           timestamp,
  expiry                 bigint,
  originator_envelope    bytea
);
