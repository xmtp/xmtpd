package migrations_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/migrations"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

// TestMigration00024_EnsureGatewayPartsV4_AttachesL1AndL2 verifies that the
// hardened partition-ensure function attaches both the L1 (per-originator LIST)
// and L2 (per-(originator, seq-band) RANGE) partitions to their parents and
// leaves no seed CHECK constraints behind.
//
// The v2 `make_meta_seq_subpart` helper had a `format()` arity bug that
// produced a bogus `seq_id_check` (see migration 00024). The v4 helpers write
// a correct CHECK for the duration of ATTACH and drop it afterwards — the test
// proves no stray `oid_check` or `seq_id_check` constraints survive.
func TestMigration00024_EnsureGatewayPartsV4_AttachesL1AndL2(t *testing.T) {
	ctx := t.Context()
	database, _ := testutils.NewRawDB(t, ctx)

	const (
		oid    = 100
		seqID  = 1
		l1Meta = "gateway_envelopes_meta_o100"
		l1Blob = "gateway_envelopes_blob_o100"
		l2Meta = "gateway_envelopes_meta_o100_s0_1000000"
		l2Blob = "gateway_envelopes_blob_o100_s0_1000000"
		bandBw = int64(1_000_000)
	)

	_, err := database.ExecContext(
		ctx,
		`SELECT ensure_gateway_parts_v4($1, $2, $3)`,
		oid, seqID, bandBw,
	)
	require.NoError(t, err, "ensure_gateway_parts_v4 should succeed on a fresh schema")

	// L1 (meta) attached to gateway_envelopes_meta.
	assertInherits(t, database, l1Meta, "gateway_envelopes_meta")
	// L1 (blob) attached to gateway_envelopes_blob.
	assertInherits(t, database, l1Blob, "gateway_envelopes_blob")
	// L2 (meta) attached to the L1 meta.
	assertInherits(t, database, l2Meta, l1Meta)
	// L2 (blob) attached to the L1 blob.
	assertInherits(t, database, l2Blob, l1Blob)

	// No residual CHECK constraints named oid_check / seq_id_check on any
	// partition child — the v4 helpers drop them after ATTACH.
	for _, child := range []string{l1Meta, l1Blob, l2Meta, l2Blob} {
		assertNoCheckConstraint(t, database, child, "oid_check")
		assertNoCheckConstraint(t, database, child, "seq_id_check")
	}
}

// TestMigration00024_EnsureGatewayPartsV4_Idempotent verifies that calling
// ensure_gateway_parts_v4 twice for the same (oid, seq) is a no-op on the
// second call — no duplicate L1/L2 children, no errors, and the pg_inherits
// row count is unchanged.
func TestMigration00024_EnsureGatewayPartsV4_Idempotent(t *testing.T) {
	ctx := t.Context()
	database, _ := testutils.NewRawDB(t, ctx)

	const (
		oid    = 101
		seqID  = 1
		bandBw = int64(1_000_000)
	)

	_, err := database.ExecContext(
		ctx,
		`SELECT ensure_gateway_parts_v4($1, $2, $3)`,
		oid, seqID, bandBw,
	)
	require.NoError(t, err)

	before := countInheritsForOriginator(t, database, oid)

	// Call again — must be a schema no-op.
	_, err = database.ExecContext(
		ctx,
		`SELECT ensure_gateway_parts_v4($1, $2, $3)`,
		oid, seqID, bandBw,
	)
	require.NoError(t, err, "second call must not error")

	after := countInheritsForOriginator(t, database, oid)
	assert.Equal(t, before, after, "second call should not create new partitions")
}

// TestMigration00024_EnsureGatewayPartsV4_CoexistsWithV3 verifies the
// "append-only migration" invariant: the v3 helpers and v4 helpers can both
// be called against the same schema, and calling v4 for (oid, seq) that were
// already seeded by v3 is a no-op (no errors, no duplicate rows).
//
// The test seeds partitions at schema version 23 via ensure_gateway_parts_v3,
// migrates to HEAD (applying 00024), then calls ensure_gateway_parts_v4 for
// the same (oid, seq) and asserts it reports success without re-attaching.
func TestMigration00024_EnsureGatewayPartsV4_CoexistsWithV3(t *testing.T) {
	ctx := t.Context()
	database, _ := testutils.NewRawDBAtVersion(t, ctx, 23)

	const (
		oid    = 102
		seqID  = 1
		bandBw = int64(1_000_000)
	)

	// Seed via v3 at schema version 23.
	_, err := database.ExecContext(
		ctx,
		`SELECT ensure_gateway_parts_v3($1, $2, $3)`,
		oid, seqID, bandBw,
	)
	require.NoError(t, err, "ensure_gateway_parts_v3 should succeed at v23")

	// Apply migration 00024.
	require.NoError(t, migrations.Migrate(ctx, database))

	before := countInheritsForOriginator(t, database, oid)

	_, err = database.ExecContext(
		ctx,
		`SELECT ensure_gateway_parts_v4($1, $2, $3)`,
		oid, seqID, bandBw,
	)
	require.NoError(t, err, "ensure_gateway_parts_v4 should be a no-op on v3-seeded partitions")

	after := countInheritsForOriginator(t, database, oid)
	assert.Equal(t, before, after, "v4 should not alter partitions already attached by v3")
}

// assertInherits requires child to be attached as a partition of parent
// (via pg_inherits). Both names are unqualified public-schema relation names.
func assertInherits(t *testing.T, database *sql.DB, child, parent string) {
	t.Helper()
	var exists bool
	err := database.QueryRowContext(
		t.Context(),
		`SELECT EXISTS (
			SELECT 1
			FROM pg_inherits i
			JOIN pg_class c ON c.oid = i.inhrelid
			JOIN pg_class p ON p.oid = i.inhparent
			WHERE c.relname = $1
			  AND p.relname = $2
		)`,
		child, parent,
	).Scan(&exists)
	require.NoError(t, err)
	assert.Truef(t, exists, "%s should be attached as a partition of %s", child, parent)
}

// assertNoCheckConstraint fails if a CHECK constraint by the given name exists
// on the given relation. The v4 helpers install a seed CHECK for the duration
// of ATTACH, then drop it — so a surviving constraint is evidence that ATTACH
// failed silently (the exact regression this migration prevents).
func assertNoCheckConstraint(t *testing.T, database *sql.DB, relName, constraintName string) {
	t.Helper()
	var exists bool
	err := database.QueryRowContext(
		t.Context(),
		`SELECT EXISTS (
			SELECT 1
			FROM pg_constraint con
			JOIN pg_class c ON c.oid = con.conrelid
			WHERE c.relname = $1
			  AND con.conname = $2
		)`,
		relName, constraintName,
	).Scan(&exists)
	require.NoError(t, err)
	assert.Falsef(
		t,
		exists,
		"%s should not carry residual CHECK constraint %s after ATTACH",
		relName, constraintName,
	)
}

// countInheritsForOriginator returns the number of pg_inherits rows whose
// child relation is named gateway_envelopes_{meta,blob}_o<oid>[...]. Used to
// assert schema idempotence.
func countInheritsForOriginator(t *testing.T, database *sql.DB, oid int) int {
	t.Helper()
	metaPrefix := fmt.Sprintf("gateway_envelopes_meta_o%d%%", oid)
	blobPrefix := fmt.Sprintf("gateway_envelopes_blob_o%d%%", oid)
	var n int
	err := database.QueryRowContext(
		t.Context(),
		`SELECT COUNT(*)
		 FROM pg_inherits i
		 JOIN pg_class c ON c.oid = i.inhrelid
		 WHERE c.relname LIKE $1
		    OR c.relname LIKE $2`,
		metaPrefix, blobPrefix,
	).Scan(&n)
	require.NoError(t, err)
	return n
}
