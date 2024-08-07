-- name: InsertStagedOriginatorEnvelope :one
INSERT INTO staged_originator_envelopes(payer_envelope)
	VALUES (@payer_envelope)
RETURNING
	*;
