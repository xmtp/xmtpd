-- name: InsertPayerLedgerEvent :exec
INSERT INTO payer_ledger_events(
        event_id,
        payer_id,
        amount_picodollars,
        event_type
    )
VALUES (
        @event_id,
        @payer_id,
        @amount_picodollars,
        @event_type
    ) ON CONFLICT (event_id) DO NOTHING;

-- name: GetPayerBalance :one
SELECT COALESCE(SUM(amount_picodollars), 0)::BIGINT AS balance
FROM payer_ledger_events
WHERE payer_id = @payer_id;

-- name: GetLastEvent :one
SELECT *
FROM payer_ledger_events
WHERE payer_id = @payer_id
    AND event_type = @event_type
ORDER BY created_at DESC
LIMIT 1;