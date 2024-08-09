-- name: InsertStagedOriginatorEnvelope :one
INSERT INTO staged_originator_envelopes(payer_envelope)
	VALUES (@payer_envelope)
RETURNING
	*;

-- name: InsertNodeInfo :one
INSERT INTO node_info(node_id, public_key)
	VALUES (@node_id, @public_key)
	RETURNING *;

-- name: SelectNodeInfo :one
SELECT * FROM node_info WHERE singleton_id = 1;
