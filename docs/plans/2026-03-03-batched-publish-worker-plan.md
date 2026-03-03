# Batched Publish Worker Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Restructure the publish worker to batch DB operations, reducing round trips from ~400 to ~4 per batch of 100 envelopes.

**Architecture:** Two-phase processing — CPU prep (per-envelope: parse, fees, sign) then batch DB ops (bulk payer upsert, batch gateway insert with usage+congestion, bulk delete staged). In-memory congestion tracking within a batch ensures fee accuracy.

**Tech Stack:** Go, PostgreSQL, sqlc, golang-migrate

**Design doc:** `docs/plans/2026-03-03-batched-publish-worker-design.md`

---

### Task 1: Migration — `insert_gateway_envelope_batch_v2` SQL Function

Creates a new SQL function that extends V1 with `p_is_reserved` boolean array and congestion increment.

**Files:**
- Create: `pkg/db/migrations/00020_insert-gateway-envelopes-batch-v2.up.sql`
- Create: `pkg/db/migrations/00020_insert-gateway-envelopes-batch-v2.down.sql`

**Step 1: Create the up migration**

Use skill `@writing-migrations` for conventions.

```sql
CREATE OR REPLACE FUNCTION insert_gateway_envelope_batch_v2(
    p_originator_node_ids     int[],
    p_originator_sequence_ids bigint[],
    p_topics                  bytea[],
    p_payer_ids               int[],
    p_gateway_times           timestamp[],
    p_expiries                bigint[],
    p_originator_envelopes    bytea[],
    p_spend_picodollars       bigint[],
    p_is_reserved             boolean[]
)
RETURNS TABLE (
    inserted_meta_rows    bigint,
    inserted_blob_rows    bigint,
    affected_usage_rows   bigint,
    affected_congestion_rows bigint
)
LANGUAGE SQL
AS $$
WITH input AS (
    SELECT
        originator_node_id,
        originator_sequence_id,
        topic,
        NULLIF(payer_id, 0) AS payer_id,
        gateway_time,
        expiry,
        originator_envelope,
        spend_picodollars,
        is_reserved
    FROM unnest(
        p_originator_node_ids,
        p_originator_sequence_ids,
        p_topics,
        p_payer_ids,
        p_gateway_times,
        p_expiries,
        p_originator_envelopes,
        p_spend_picodollars,
        p_is_reserved
    ) AS t(
        originator_node_id,
        originator_sequence_id,
        topic,
        payer_id,
        gateway_time,
        expiry,
        originator_envelope,
        spend_picodollars,
        is_reserved
    )
),

m AS (
    INSERT INTO gateway_envelopes_meta (
        originator_node_id,
        originator_sequence_id,
        topic,
        payer_id,
        gateway_time,
        expiry
    )
    SELECT originator_node_id, originator_sequence_id, topic, payer_id, gateway_time, expiry
    FROM input
    ON CONFLICT DO NOTHING
    RETURNING originator_node_id, originator_sequence_id, payer_id, gateway_time
),

b AS (
    INSERT INTO gateway_envelope_blobs (
        originator_node_id,
        originator_sequence_id,
        originator_envelope
    )
    SELECT originator_node_id, originator_sequence_id, originator_envelope
    FROM input
    ON CONFLICT DO NOTHING
    RETURNING originator_node_id, originator_sequence_id
),

m_with_spend AS (
    SELECT
        m.originator_node_id,
        m.originator_sequence_id,
        m.payer_id,
        m.gateway_time,
        i.spend_picodollars,
        i.is_reserved
    FROM m
    JOIN b USING (originator_node_id, originator_sequence_id)
    JOIN input i USING (originator_node_id, originator_sequence_id)
),

u_prep AS (
    SELECT
        payer_id,
        originator_node_id AS originator_id,
        floor(extract(epoch from gateway_time) / 60)::int AS minutes_since_epoch,
        sum(spend_picodollars)::bigint AS spend_picodollars,
        max(originator_sequence_id)::bigint AS last_sequence_id,
        count(*)::int AS message_count
    FROM m_with_spend
    WHERE payer_id IS NOT NULL AND NOT is_reserved
    GROUP BY 1, 2, 3
),

u AS (
    INSERT INTO unsettled_usage (
        payer_id,
        originator_id,
        minutes_since_epoch,
        spend_picodollars,
        last_sequence_id,
        message_count
    )
    SELECT payer_id, originator_id, minutes_since_epoch, spend_picodollars, last_sequence_id, message_count
    FROM u_prep
    ORDER BY payer_id, originator_id, minutes_since_epoch
    ON CONFLICT (payer_id, originator_id, minutes_since_epoch) DO UPDATE
    SET
        spend_picodollars = unsettled_usage.spend_picodollars + EXCLUDED.spend_picodollars,
        message_count     = unsettled_usage.message_count + EXCLUDED.message_count,
        last_sequence_id  = GREATEST(unsettled_usage.last_sequence_id, EXCLUDED.last_sequence_id)
    RETURNING 1
),

c_prep AS (
    SELECT
        originator_node_id AS originator_id,
        floor(extract(epoch from gateway_time) / 60)::int AS minutes_since_epoch,
        count(*)::int AS num_messages
    FROM m_with_spend
    WHERE NOT is_reserved
    GROUP BY 1, 2
),

c AS (
    INSERT INTO originator_congestion (originator_id, minutes_since_epoch, num_messages)
    SELECT originator_id, minutes_since_epoch, num_messages
    FROM c_prep
    ON CONFLICT (originator_id, minutes_since_epoch) DO UPDATE
    SET num_messages = originator_congestion.num_messages + EXCLUDED.num_messages
    RETURNING 1
)

SELECT
    (SELECT COUNT(*) FROM m) AS inserted_meta_rows,
    (SELECT COUNT(*) FROM b) AS inserted_blob_rows,
    (SELECT COUNT(*) FROM u) AS affected_usage_rows,
    (SELECT COUNT(*) FROM c) AS affected_congestion_rows;
$$;
```

**Step 2: Create the down migration**

```sql
DROP FUNCTION IF EXISTS insert_gateway_envelope_batch_v2(int[], bigint[], bytea[], int[], timestamp[], bigint[], bytea[], bigint[], boolean[]);
```

**Step 3: Verify migration applies**

Run: `dev/up` then connect with `dev/psql` and run `\df insert_gateway_envelope_batch_v2` to confirm the function exists.

**Step 4: Commit**

```
gt create -m "Add insert_gateway_envelope_batch_v2 migration

Extends V1 with p_is_reserved boolean array parameter and congestion
CTE. Reserved topics skip both unsettled_usage and originator_congestion
increments."
```

---

### Task 2: New sqlc Queries

Adds `BulkFindOrCreatePayers`, `BulkDeleteStagedOriginatorEnvelopes`, and `InsertGatewayEnvelopeBatchV2` sqlc queries. Removes old `DeleteStagedOriginatorEnvelope`.

**Files:**
- Modify: `pkg/db/sqlc/payer_reports.sql` (add `BulkFindOrCreatePayers`)
- Modify: `pkg/db/sqlc/envelopes.sql` (replace `DeleteStagedOriginatorEnvelope` with `BulkDeleteStagedOriginatorEnvelopes`)
- Modify: `pkg/db/sqlc/envelopes_v2.sql` (add `InsertGatewayEnvelopeBatchV2`)

Use skill `@writing-queries` for conventions.

**Step 1: Add BulkFindOrCreatePayers to payer_reports.sql**

Append after line 14 (after `FindOrCreatePayer` query):

```sql
-- name: BulkFindOrCreatePayers :many
WITH input AS (
    SELECT address FROM unnest(@addresses::TEXT[]) AS t(address)
),
ins AS (
    INSERT INTO payers(address)
    SELECT address FROM input
    ON CONFLICT (address) DO NOTHING
    RETURNING id, address
)
SELECT address, id
FROM ins
UNION ALL
SELECT i.address, p.id
FROM input i
JOIN payers p ON p.address = i.address
WHERE i.address NOT IN (SELECT address FROM ins);
```

**Step 2: Replace DeleteStagedOriginatorEnvelope in envelopes.sql**

Replace lines 12-14 of `pkg/db/sqlc/envelopes.sql`:

```sql
-- name: DeleteStagedOriginatorEnvelope :execrows
DELETE FROM staged_originator_envelopes
WHERE id = @id;
```

With:

```sql
-- name: BulkDeleteStagedOriginatorEnvelopes :execrows
DELETE FROM staged_originator_envelopes WHERE id = ANY(@ids::BIGINT[]);
```

**Step 3: Add InsertGatewayEnvelopeBatchV2 to envelopes_v2.sql**

Append at end of `pkg/db/sqlc/envelopes_v2.sql`:

```sql
-- name: InsertGatewayEnvelopeBatchV2 :one
SELECT
    inserted_meta_rows::bigint,
    inserted_blob_rows::bigint,
    affected_usage_rows::bigint,
    affected_congestion_rows::bigint
FROM insert_gateway_envelope_batch_v2(
    @originator_node_ids::int[],
    @originator_sequence_ids::bigint[],
    @topics::bytea[],
    @payer_ids::int[],
    @gateway_times::timestamp[],
    @expiries::bigint[],
    @originator_envelopes::bytea[],
    @spend_picodollars::bigint[],
    @is_reserved::boolean[]
);
```

**Step 4: Run code generation**

Run: `dev/gen/sqlc`

**Step 5: Fix all Go compilation errors from removed `DeleteStagedOriginatorEnvelope`**

Two callers need updating:
- `pkg/api/message/publish_worker.go:261` — will be rewritten in Task 6
- `pkg/db/bench/hot_path_bench_test.go:186,258,354` — will be rewritten in Task 7

For now, create a temporary compatibility wrapper or update these call sites to use the new bulk delete with a single-element slice. The simplest fix:

In `pkg/api/message/publish_worker.go:261`, change:
```go
deleted, err := p.store.WriteQuery().DeleteStagedOriginatorEnvelope(p.ctx, stagedEnv.ID)
```
to:
```go
deleted, err := p.store.WriteQuery().BulkDeleteStagedOriginatorEnvelopes(p.ctx, []int64{stagedEnv.ID})
```

In `pkg/db/bench/hot_path_bench_test.go`, change all `DeleteStagedOriginatorEnvelope(ctx, id)` calls to `BulkDeleteStagedOriginatorEnvelopes(ctx, []int64{id})`.

**Step 6: Verify compilation**

Run: `go build ./...`

**Step 7: Commit**

```
gt create -m "Add batch sqlc queries and replace single-row delete

New queries: BulkFindOrCreatePayers, BulkDeleteStagedOriginatorEnvelopes,
InsertGatewayEnvelopeBatchV2. Removes DeleteStagedOriginatorEnvelope."
```

---

### Task 3: GatewayEnvelopeBatch V2 Types and Wrapper

Adds `IsReserved` field to batch types and creates the V2 Go wrapper.

**Files:**
- Modify: `pkg/db/types/gateway_envelope_batch.go` (add `IsReserved`, add `ToParamsV2`)
- Modify: `pkg/db/gateway_envelope_batch.go` (add V2 wrapper functions)
- Test: `pkg/db/gateway_envelope_batch_test.go` (add V2 tests)

**Step 1: Write failing test for V2 batch insert**

Add a test to `pkg/db/gateway_envelope_batch_test.go` that calls `InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage` with a mix of reserved and non-reserved envelopes. Assert:
- All envelopes are inserted (meta + blob rows)
- Only non-reserved envelopes create `unsettled_usage` rows
- Only non-reserved envelopes create `originator_congestion` rows
- Reserved envelopes have zero usage/congestion impact

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/db/ -run TestBatchInsertV2 -v`
Expected: compilation error (function doesn't exist yet)

**Step 3: Add `IsReserved` field to `GatewayEnvelopeRow`**

In `pkg/db/types/gateway_envelope_batch.go`, add to the struct:

```go
type GatewayEnvelopeRow struct {
	OriginatorNodeID     int32
	OriginatorSequenceID int64
	Topic                []byte
	PayerID              int32
	GatewayTime          time.Time
	Expiry               int64
	OriginatorEnvelope   []byte
	SpendPicodollars     int64
	IsReserved           bool  // new field
}
```

**Step 4: Add `ToParamsV2` method**

Add a new method that returns `queries.InsertGatewayEnvelopeBatchV2Params` (generated by sqlc in Task 2). Same structure as `ToParams()` but includes the `IsReserved` boolean array.

```go
func (b *GatewayEnvelopeBatch) ToParamsV2() queries.InsertGatewayEnvelopeBatchV2Params {
	n := b.Len()
	b.ensureOrdered()

	params := queries.InsertGatewayEnvelopeBatchV2Params{
		OriginatorNodeIds:     make([]int32, n),
		OriginatorSequenceIds: make([]int64, n),
		Topics:                make([][]byte, n),
		PayerIds:              make([]int32, n),
		GatewayTimes:          make([]time.Time, n),
		Expiries:              make([]int64, n),
		OriginatorEnvelopes:   make([][]byte, n),
		SpendPicodollars:      make([]int64, n),
		IsReserved:            make([]bool, n),
	}

	for i, row := range b.Envelopes {
		params.OriginatorNodeIds[i] = row.OriginatorNodeID
		params.OriginatorSequenceIds[i] = row.OriginatorSequenceID
		params.Topics[i] = row.Topic
		params.PayerIds[i] = row.PayerID
		params.GatewayTimes[i] = row.GatewayTime
		params.Expiries[i] = row.Expiry
		params.OriginatorEnvelopes[i] = row.OriginatorEnvelope
		params.SpendPicodollars[i] = row.SpendPicodollars
		params.IsReserved[i] = row.IsReserved
	}

	return params
}
```

**Step 5: Add V2 wrapper functions in `gateway_envelope_batch.go`**

```go
func InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(
	ctx context.Context,
	db *sql.DB,
	input *types.GatewayEnvelopeBatch,
) (int64, error) {
	return RunInTxWithResult(ctx, db, &sql.TxOptions{},
		func(ctx context.Context, q *queries.Queries) (int64, error) {
			return InsertGatewayEnvelopeBatchV2Transactional(ctx, q, input)
		})
}

func InsertGatewayEnvelopeBatchV2Transactional(
	ctx context.Context,
	q *queries.Queries,
	input *types.GatewayEnvelopeBatch,
) (int64, error) {
	// Same structure as InsertGatewayEnvelopeBatchTransactional but calls
	// q.InsertGatewayEnvelopeBatchV2(ctx, params) and uses ToParamsV2().
	// Same savepoint + EnsureGatewayParts retry pattern.
}
```

**Step 6: Run test to verify it passes**

Run: `go test ./pkg/db/ -run TestBatchInsertV2 -v`
Expected: PASS

**Step 7: Commit**

```
gt create -m "Add V2 batch insert types and wrapper with IsReserved support

GatewayEnvelopeRow gains IsReserved field. ToParamsV2 maps to the new
SQL function. V2 wrapper uses same savepoint retry pattern as V1."
```

---

### Task 4: Fee Calculator — In-Memory Congestion Adjustment

Modifies `CalculateCongestionFee` to accept an `additionalMessages` parameter for in-memory congestion tracking.

**Files:**
- Modify: `pkg/fees/interface.go:27-39` (update `IFeeCalculator` interface)
- Modify: `pkg/fees/calculator.go:57-96` (update `CalculateCongestionFee`)
- Modify: `pkg/fees/calculator_test.go` (add test for additionalMessages)

**Step 1: Write failing test**

Add to `pkg/fees/calculator_test.go`:

```go
func TestCalculateCongestionFeeWithAdditionalMessages(t *testing.T) {
	calculator := setupCalculator()
	db, _ := testutils.NewRawDB(t, context.Background())

	ctx := context.Background()
	querier := queries.New(db)
	originatorID := uint32(testutils.RandomInt32())
	messageTime := time.Now()
	minutesSinceEpoch := utils.MinutesSinceEpoch(messageTime)

	// Add congestion just below threshold (targetRatePerMinute=4)
	addCongestion(t, querier, originatorID, minutesSinceEpoch, 4)

	// With 0 additional messages, fee should be 0 (at target, not above)
	fee0, err := calculator.CalculateCongestionFee(ctx, querier, messageTime, originatorID, 0)
	require.NoError(t, err)
	require.Equal(t, currency.PicoDollar(0), fee0)

	// With additional messages pushing above target, fee should be > 0
	fee5, err := calculator.CalculateCongestionFee(ctx, querier, messageTime, originatorID, 5)
	require.NoError(t, err)
	require.Greater(t, fee5, currency.PicoDollar(0))

	// More additional messages should produce higher fee
	fee10, err := calculator.CalculateCongestionFee(ctx, querier, messageTime, originatorID, 10)
	require.NoError(t, err)
	require.Greater(t, fee10, fee5)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/fees/ -run TestCalculateCongestionFeeWithAdditionalMessages -v`
Expected: compilation error (wrong number of arguments)

**Step 3: Update the interface**

In `pkg/fees/interface.go`, change `CalculateCongestionFee` signature:

```go
CalculateCongestionFee(
    ctx context.Context,
    querier *queries.Queries,
    messageTime time.Time,
    originatorID uint32,
    additionalMessages int32,
) (currency.PicoDollar, error)
```

**Step 4: Update the implementation**

In `pkg/fees/calculator.go`, update `CalculateCongestionFee`:

```go
func (c *FeeCalculator) CalculateCongestionFee(
	ctx context.Context,
	querier *queries.Queries,
	messageTime time.Time,
	originatorID uint32,
	additionalMessages int32,
) (currency.PicoDollar, error) {
	last5MinutesCongestion, err := db.Get5MinutesOfCongestion(
		ctx,
		querier,
		int32(originatorID),
		int32(utils.MinutesSinceEpoch(messageTime)),
	)
	if err != nil {
		return 0, err
	}

	// Add in-memory batch adjustment to the current minute (index 0)
	last5MinutesCongestion[0] += additionalMessages

	rates, err := c.ratesFetcher.GetRates(messageTime)
	if err != nil {
		return 0, err
	}

	// ... rest unchanged
```

**Step 5: Fix all callers**

Search for all `CalculateCongestionFee` calls and add `, 0` as the last argument. Known callers:
- `pkg/api/message/publish_worker.go:291-296` — add `, 0` (will be rewritten in Task 6)
- `pkg/sync/envelope_sink.go` — add `, 0`
- `pkg/migrator/transformer.go` — add `, 0`
- Any generated mocks — regenerate with `dev/gen/mocks`

Run: `grep -rn "CalculateCongestionFee" --include="*.go"` to find all callers.

**Step 6: Run test to verify it passes**

Run: `go test ./pkg/fees/ -run TestCalculateCongestionFee -v`
Expected: all congestion fee tests PASS

**Step 7: Run full test suite to verify no regressions**

Run: `dev/test`
Expected: PASS

**Step 8: Commit**

```
gt create -m "Add additionalMessages parameter to CalculateCongestionFee

Allows callers to pass an in-memory running count of batch-processed
messages for accurate per-envelope congestion fee calculation without
requiring a DB round trip per envelope."
```

---

### Task 5: Refactor Publish Worker to Batch Mode

The core change — replaces the per-envelope `publishStagedEnvelope` loop with a two-phase `publishBatch` method.

**Files:**
- Modify: `pkg/api/message/publish_worker.go`

**Step 1: Add the `publishBatch` method**

This is the new method that replaces the inner loop in `start()`. Structure:

```go
// publishBatch processes a batch of staged envelopes in two phases:
// Phase 1: CPU prep (per-envelope) - parse, fees, sign
// Phase 2: Batch DB ops - bulk payer upsert, batch insert, bulk delete
func (p *publishWorker) publishBatch(batch []queries.StagedOriginatorEnvelope) bool {
	originatorID := int32(p.registrant.NodeID())

	// Phase 1: CPU prep
	type preparedEnvelope struct {
		staged         queries.StagedOriginatorEnvelope
		originatorBytes []byte
		payerAddress   string
		topic          []byte
		isReserved     bool
		baseFee        currency.PicoDollar
		congestionFee  currency.PicoDollar
		retentionDays  uint32
		expiry         int64
		gatewayTime    time.Time
	}

	// Query congestion once for the entire batch
	// Use the first envelope's time as representative (all in same batch window)
	congestionData, err := db.Get5MinutesOfCongestion(
		p.ctx,
		p.store.Query(),
		originatorID,
		int32(utils.MinutesSinceEpoch(batch[0].OriginatorTime)),
	)
	if err != nil {
		p.logger.Error("failed to get congestion data", zap.Error(err))
		return false
	}

	prepared := make([]preparedEnvelope, 0, len(batch))
	var additionalMessages int32

	for _, stagedEnv := range batch {
		// Parse payer envelope
		env, err := envelopes.NewPayerEnvelopeFromBytes(stagedEnv.PayerEnvelope)
		if err != nil {
			p.logger.Warn("failed to unmarshall originator envelope", zap.Error(err))
			return false
		}

		parsedTopic, err := topic.ParseTopic(stagedEnv.Topic)
		if err != nil {
			return false
		}

		isReserved := parsedTopic.IsReserved()
		retentionDays := env.RetentionDays()

		var baseFee, congestionFee currency.PicoDollar

		if !isReserved {
			baseFee, err = p.feeCalculator.CalculateBaseFee(
				stagedEnv.OriginatorTime,
				int64(len(stagedEnv.PayerEnvelope)),
				retentionDays,
			)
			if err != nil {
				p.logger.Error("failed to calculate base fee", zap.Error(err))
				return false
			}

			// Calculate congestion fee using DB data + in-memory running count
			// Build the adjusted congestion array for this envelope
			adjusted := congestionData
			adjusted[0] += additionalMessages

			rates, err := p.feeCalculator.GetRates(stagedEnv.OriginatorTime)
			// ... use CalculateCongestion directly with adjusted data
			// OR use the new additionalMessages parameter:
			congestionFee, err = p.feeCalculator.CalculateCongestionFee(
				p.ctx,
				p.store.Query(),
				stagedEnv.OriginatorTime,
				uint32(originatorID),
				additionalMessages,
			)
			if err != nil {
				p.logger.Error("failed to calculate congestion fee", zap.Error(err))
				return false
			}

			additionalMessages++
		}

		// Sign the envelope
		originatorEnv, err := p.registrant.SignStagedEnvelope(
			stagedEnv, baseFee, congestionFee, retentionDays,
		)
		if err != nil {
			p.logger.Error("failed to sign staged envelope", zap.Error(err))
			return false
		}

		validatedEnvelope, err := envelopes.NewOriginatorEnvelope(originatorEnv)
		if err != nil {
			p.logger.Error("failed to validate originator envelope", zap.Error(err))
			return false
		}

		originatorBytes, err := validatedEnvelope.Bytes()
		if err != nil {
			p.logger.Error("failed to marshal originator envelope", zap.Error(err))
			return false
		}

		payerAddress, err := validatedEnvelope.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
		if err != nil {
			p.logger.Error("failed to recover payer address", zap.Error(err))
			return false
		}

		prepared = append(prepared, preparedEnvelope{
			staged:          stagedEnv,
			originatorBytes: originatorBytes,
			payerAddress:    payerAddress.Hex(),
			topic:           stagedEnv.Topic,
			isReserved:      isReserved,
			baseFee:         baseFee,
			congestionFee:   congestionFee,
			retentionDays:   retentionDays,
			expiry:          int64(validatedEnvelope.UnsignedOriginatorEnvelope.Proto().GetExpiryUnixtime()),
			gatewayTime:     stagedEnv.OriginatorTime,
		})
	}

	// Phase 2: Batch DB ops

	// 2a. Bulk find/create payers (deduplicated)
	uniqueAddresses := deduplicateAddresses(prepared)
	payerRows, err := p.store.WriteQuery().BulkFindOrCreatePayers(p.ctx, uniqueAddresses)
	if err != nil {
		p.logger.Error("failed to bulk find/create payers", zap.Error(err))
		return false
	}
	payerMap := make(map[string]int32, len(payerRows))
	for _, row := range payerRows {
		payerMap[row.Address] = row.ID
	}

	// 2b. Build batch and insert
	batchInput := types.NewGatewayEnvelopeBatch()
	for _, prep := range prepared {
		batchInput.Add(types.GatewayEnvelopeRow{
			OriginatorNodeID:     originatorID,
			OriginatorSequenceID: prep.staged.ID,
			Topic:                prep.topic,
			PayerID:              payerMap[prep.payerAddress],
			GatewayTime:          prep.gatewayTime,
			Expiry:               prep.expiry,
			OriginatorEnvelope:   prep.originatorBytes,
			SpendPicodollars:     int64(prep.baseFee) + int64(prep.congestionFee),
			IsReserved:           prep.isReserved,
		})
	}

	inserted, err := db.InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(
		p.ctx, p.store.DB(), batchInput,
	)
	if p.ctx.Err() != nil {
		return false
	}
	if err != nil {
		p.logger.Error("failed to batch insert gateway envelopes", zap.Error(err))
		return false
	}

	// 2c. Bulk delete staged envelopes
	stagedIDs := make([]int64, len(prepared))
	for i, prep := range prepared {
		stagedIDs[i] = prep.staged.ID
	}
	_, err = p.store.WriteQuery().BulkDeleteStagedOriginatorEnvelopes(p.ctx, stagedIDs)
	if p.ctx.Err() != nil {
		return true
	}
	if err != nil {
		p.logger.Error("failed to bulk delete staged envelopes", zap.Error(err))
		// Envelopes already inserted, safe to continue
		return true
	}

	// Emit metrics for each envelope
	for _, prep := range prepared {
		metrics.EmitAPIStagedEnvelopeProcessingDelay(time.Since(prep.staged.OriginatorTime))
	}

	if inserted > 0 {
		p.logger.Info("batch published",
			zap.Int("batch_size", len(prepared)),
			zap.Int64("inserted", inserted),
		)
	}

	return true
}
```

**Step 2: Update the `start()` method**

Replace the per-envelope loop in `start()` (lines 101-123):

```go
func (p *publishWorker) start() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case batch, ok := <-p.listener:
			if !ok {
				p.logger.Error("listener is closed")
				return
			}

			for !p.publishBatch(batch) {
				time.Sleep(p.sleepOnFailureTime)
			}

			if len(batch) > 0 {
				p.lastProcessed.Store(batch[len(batch)-1].ID)
			}
		}
	}
}
```

**Step 3: Remove old `publishStagedEnvelope` and `calculateFees` methods**

Delete the `publishStagedEnvelope` method (lines 126-274) and `calculateFees` method (lines 276-302) since all logic is now in `publishBatch`.

**Step 4: Add helper function for address deduplication**

```go
func deduplicateAddresses(prepared []preparedEnvelope) []string {
	seen := make(map[string]struct{}, len(prepared))
	result := make([]string, 0, len(prepared))
	for _, p := range prepared {
		if _, ok := seen[p.payerAddress]; !ok {
			seen[p.payerAddress] = struct{}{}
			result = append(result, p.payerAddress)
		}
	}
	return result
}
```

**Step 5: Verify compilation**

Run: `go build ./...`

**Step 6: Address the congestion fee DB call issue**

Note: The current `CalculateCongestionFee` still queries the DB internally. For the batch optimization to work (1 DB call per batch, not per envelope), we need to either:

Option A: Call `Get5MinutesOfCongestion` once before the loop and use `CalculateCongestion` directly with the adjusted array. This bypasses the `IFeeCalculator` interface for congestion.

Option B: Accept the per-envelope DB call for now (still batches the writes). This is simpler and still provides significant improvement.

**Recommendation:** Option A for maximum benefit. Query congestion once, then in the loop use `CalculateCongestion()` + rates directly:

```go
// Before the loop:
congestionData, err := db.Get5MinutesOfCongestion(...)
rates, err := p.feeCalculator.GetRates(batch[0].OriginatorTime)  // need to expose GetRates or fetch rates

// Inside the loop for non-reserved envelopes:
adjusted := congestionData
adjusted[0] += additionalMessages
congestionUnits := fees.CalculateCongestion(adjusted, int32(rates.TargetRatePerMinute))
congestionFee = rates.CongestionFee * currency.PicoDollar(congestionUnits)
additionalMessages++
```

This requires exposing `GetRates` from the fee calculator (or adding it to the interface). Check if `IRatesFetcher` is accessible.

**Step 7: Run existing tests**

Run: `dev/test`
Expected: PASS (existing publish worker tests should still pass)

**Step 8: Commit**

```
gt create -m "Refactor publish worker to batch mode

Replaces per-envelope publishStagedEnvelope with two-phase publishBatch:
Phase 1: CPU prep (parse, fees with in-memory congestion tracking, sign)
Phase 2: Batch DB ops (bulk payer upsert, batch insert, bulk delete)

Reduces DB round trips from ~4N to ~4 per batch of N envelopes."
```

---

### Task 6: Lint, Format, and Final Verification

**Step 1: Run linter**

Run: `dev/lint-fix`

**Step 2: Run full test suite**

Run: `dev/test`
Expected: PASS

**Step 3: Fix any issues found**

**Step 4: Commit fixes if any**

```
gt create -m "Fix lint and formatting issues"
```

---

### Task 7: Update Benchmarks

Updates the hot path benchmarks to use the new batch APIs for comparison.

**Files:**
- Modify: `pkg/db/bench/hot_path_bench_test.go`

**Step 1: Add a batched version of BenchmarkHotPathBatchCycle**

Add `BenchmarkHotPathBatchCycleV2` that uses the V2 batch insert + bulk delete instead of per-envelope operations. Same batch sizes (1, 10, 100, 500) for direct comparison.

```go
func BenchmarkHotPathBatchCycleV2(b *testing.B) {
	batchSizes := []int{1, 10, 100, 500}

	for _, batchSize := range batchSizes {
		b.Run(fmt.Sprintf("batch=%d", batchSize), func(b *testing.B) {
			// Same setup as BenchmarkHotPathBatchCycle
			// But timed section uses:
			// 1. BulkFindOrCreatePayers (single call)
			// 2. InsertGatewayEnvelopeBatchV2 (single call)
			// 3. BulkDeleteStagedOriginatorEnvelopes (single call)
		})
	}
}
```

**Step 2: Run benchmarks to compare**

Run: `go test -tags bench -bench=BenchmarkHotPathBatchCycle -benchmem -count=3 -timeout=30m -run='^$' ./pkg/db/bench/`

Compare V1 vs V2 output using `benchstat`.

**Step 3: Commit**

```
gt create -m "Add V2 batch cycle benchmark for comparison

BenchmarkHotPathBatchCycleV2 uses bulk payer lookup, batch insert V2,
and bulk delete for direct comparison with the per-envelope V1 path."
```

---

### Task 8: Critical Risk Tests

Tests targeting the riskiest assumptions in this change. These are the tests most likely to catch regressions.

**Files:**
- Modify: `pkg/db/gateway_envelope_batch_test.go` (V2 risk tests)
- Modify: `pkg/fees/calculator_test.go` (congestion parity test)

**Risk 1: Reserved topics must not affect usage or congestion**

Add `TestBatchInsertV2_ReservedTopicsNoUsageNoCongestion`:

```go
func TestBatchInsertV2_ReservedTopicsNoUsageNoCongestion(t *testing.T) {
	// Setup: create a batch with 3 reserved and 2 non-reserved envelopes
	// All have valid payer IDs and non-zero spend_picodollars

	// Insert the batch via V2

	// Assert: unsettled_usage rows exist ONLY for non-reserved envelopes
	//   - Query unsettled_usage for the payer/originator/minute
	//   - Verify spend_picodollars matches sum of non-reserved envelopes only
	//   - Verify message_count == 2 (not 5)

	// Assert: originator_congestion incremented ONLY by non-reserved count
	//   - Query originator_congestion for the originator/minute
	//   - Verify num_messages == 2 (not 5)
}
```

**Risk 2: All-reserved batch creates no usage/congestion rows at all**

Add `TestBatchInsertV2_AllReservedNoSideEffects`:

```go
func TestBatchInsertV2_AllReservedNoSideEffects(t *testing.T) {
	// Setup: batch of 5 envelopes, ALL marked is_reserved=true
	// Insert via V2

	// Assert: zero unsettled_usage rows for this originator/minute
	// Assert: zero originator_congestion rows for this originator/minute
	// Assert: all 5 meta + blob rows inserted correctly
}
```

**Risk 3: In-memory congestion tracking produces same fees as sequential**

Add `TestCongestionFeeParity_BatchVsSequential`:

```go
func TestCongestionFeeParity_BatchVsSequential(t *testing.T) {
	// Setup: pre-seed congestion data in DB (e.g., 3 messages in current minute)
	calculator := setupCalculator()
	db, _ := testutils.NewRawDB(t, context.Background())
	querier := queries.New(db)

	originatorID := uint32(testutils.RandomInt32())
	messageTime := time.Now()
	minutesSinceEpoch := utils.MinutesSinceEpoch(messageTime)
	addCongestion(t, querier, originatorID, minutesSinceEpoch, 3)

	// Sequential: calculate 10 fees, incrementing congestion in DB each time
	sequentialFees := make([]currency.PicoDollar, 10)
	for i := range 10 {
		fee, err := calculator.CalculateCongestionFee(
			context.Background(), querier, messageTime, originatorID, 0,
		)
		require.NoError(t, err)
		sequentialFees[i] = fee

		// Simulate the real sequential path: increment congestion in DB
		err = querier.IncrementOriginatorCongestion(context.Background(),
			queries.IncrementOriginatorCongestionParams{
				OriginatorID:      int32(originatorID),
				MinutesSinceEpoch: minutesSinceEpoch,
			},
		)
		require.NoError(t, err)
	}

	// Batched: same starting state, calculate 10 fees with additionalMessages
	// Reset: create fresh DB with same seed
	db2, _ := testutils.NewRawDB(t, context.Background())
	querier2 := queries.New(db2)
	addCongestion(t, querier2, originatorID, minutesSinceEpoch, 3)

	batchedFees := make([]currency.PicoDollar, 10)
	for i := range 10 {
		fee, err := calculator.CalculateCongestionFee(
			context.Background(), querier2, messageTime, originatorID, int32(i),
		)
		require.NoError(t, err)
		batchedFees[i] = fee
	}

	// Assert: every fee is identical
	for i := range 10 {
		require.Equal(t, sequentialFees[i], batchedFees[i],
			"fee mismatch at index %d: sequential=%d, batched=%d",
			i, sequentialFees[i], batchedFees[i],
		)
	}
}
```

**Risk 4: BulkFindOrCreatePayers handles duplicate addresses in same batch**

Add `TestBulkFindOrCreatePayers_DuplicateAddresses`:

```go
func TestBulkFindOrCreatePayers_DuplicateAddresses(t *testing.T) {
	db, _ := testutils.NewRawDB(t, context.Background())
	querier := queries.New(db)

	addr1 := utils.HexEncode(testutils.RandomBytes(20))
	addr2 := utils.HexEncode(testutils.RandomBytes(20))

	// Pass duplicates: [addr1, addr2, addr1, addr1]
	rows, err := querier.BulkFindOrCreatePayers(context.Background(),
		[]string{addr1, addr2, addr1, addr1},
	)
	require.NoError(t, err)

	// Should return exactly as many rows as unique addresses (2), not 4
	// Build map to verify
	idMap := make(map[string]int32)
	for _, row := range rows {
		idMap[row.Address] = row.ID
	}
	require.Len(t, idMap, 2)
	require.NotZero(t, idMap[addr1])
	require.NotZero(t, idMap[addr2])
}
```

**NOTE:** If this test reveals the SQL returns duplicate rows (one per input occurrence), we need to either deduplicate input addresses before calling the query, or add `DISTINCT` to the SQL. The publish worker already deduplicates via `deduplicateAddresses()`, but the SQL should also be robust.

**Risk 5: Batch idempotency — concurrent V2 inserts don't double-count**

Add `TestBatchInsertV2_ConcurrentIdempotency`:

```go
func TestBatchInsertV2_ConcurrentIdempotency(t *testing.T) {
	// Setup: create identical batch input
	// Run two goroutines inserting the same batch simultaneously

	var wg sync.WaitGroup
	var inserted1, inserted2 int64

	wg.Add(2)
	go func() {
		defer wg.Done()
		inserted1, _ = db.InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(ctx, rawDB, input1)
	}()
	go func() {
		defer wg.Done()
		inserted2, _ = db.InsertGatewayEnvelopeBatchV2AndIncrementUnsettledUsage(ctx, rawDB, input2)
	}()
	wg.Wait()

	// Assert: exactly one winner (N rows), one loser (0 rows)
	require.Equal(t, int64(N), inserted1+inserted2)

	// Assert: unsettled_usage was incremented exactly once (not twice)
	usage, err := querier.GetPayerUnsettledUsage(...)
	require.Equal(t, expectedSpend, usage.TotalSpendPicodollars)

	// Assert: originator_congestion was incremented exactly once
	congestion, err := querier.SumOriginatorCongestion(...)
	require.Equal(t, int64(N), congestion.NumMessages)
}
```

**Step 1: Write all risk tests**
**Step 2: Run tests — expect failures for unimplemented V2 (if running before Task 3)**
**Step 3: After Tasks 3-5 complete, all risk tests should pass**
**Step 4: Commit**

```
gt create -m "Add critical risk tests for batched publish worker

Tests: reserved topic isolation, congestion fee parity (batch vs
sequential), bulk payer duplicate handling, concurrent insert
idempotency."
```

---

### Task 9: Submit Stack

**Step 1: Submit the full stack**

Run: `gt submit --no-interactive`

**Step 2: Verify CI passes**

Check the PR links and monitor CI.
