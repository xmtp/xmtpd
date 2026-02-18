package migrations_test

import (
	"database/sql"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

const currentMigration int64 = 19

var (
	originatorIDs = []int32{100, 200, 300}
	topicA        = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicA")).Bytes()
	topicB        = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicB")).Bytes()
	topicC        = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicC")).Bytes()
)

func populateDatabase(t *testing.T, database *sql.DB) {
	querier := queries.New(database)

	privKey1, err := crypto.GenerateKey()
	require.NoError(t, err)

	_, err = querier.InsertNodeInfo(
		t.Context(),
		queries.InsertNodeInfoParams{
			NodeID:    100,
			PublicKey: crypto.FromECDSAPub(&privKey1.PublicKey),
		},
	)
	require.NoError(t, err)

	privKey2, err := crypto.GenerateKey()
	require.NoError(t, err)

	_, err = querier.InsertNodeInfo(
		t.Context(),
		queries.InsertNodeInfoParams{
			NodeID:    200,
			PublicKey: crypto.FromECDSAPub(&privKey2.PublicKey),
		},
	)
	require.NoError(t, err)

	privKey3, err := crypto.GenerateKey()
	require.NoError(t, err)

	_, err = querier.InsertNodeInfo(
		t.Context(),
		queries.InsertNodeInfoParams{
			NodeID:    300,
			PublicKey: crypto.FromECDSAPub(&privKey3.PublicKey),
		},
	)
	require.NoError(t, err)

	payerID1 := testutils.CreatePayer(t, database, testutils.RandomAddress().Hex())
	require.NotZero(t, payerID1, "payerID1 is zero")

	payerID2 := testutils.CreatePayer(t, database, testutils.RandomAddress().Hex())
	require.NotZero(t, payerID2, "payerID2 is zero")

	payerID3 := testutils.CreatePayer(t, database, testutils.RandomAddress().Hex())
	require.NotZero(t, payerID3, "payerID3 is zero")

	// Insert envelopes for each originator across 5 sequence ID bands.
	// Make sure we generate 5 partitions per originator.
	// Each partition is populated with 3 envelopes, one for each topic.
	var (
		topics   = [][]byte{topicA, topicB, topicC}
		payerIDs = []int32{payerID1, payerID2, payerID3}
	)

	for _, originatorID := range originatorIDs {
		for band := range 5 {
			baseSeqID := int64(band) * db.GatewayEnvelopeBandWidth
			var rows []queries.InsertGatewayEnvelopeParams
			for k, topic := range topics {
				seqID := baseSeqID + int64(k)
				rows = append(rows, queries.InsertGatewayEnvelopeParams{
					OriginatorNodeID:     originatorID,
					OriginatorSequenceID: seqID,
					Topic:                topic,
					PayerID:              db.NullInt32(payerIDs[k]),
					OriginatorEnvelope: testutils.Marshal(
						t,
						envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(
							t,
							uint32(originatorID),
							uint64(seqID),
							topic,
						),
					),
				})
			}
			testutils.InsertGatewayEnvelopes(t, database, rows)
		}
	}
}

/* Test entry point */

func TestMigrations(t *testing.T) {
	var (
		ctx         = t.Context()
		database, _ = testutils.NewRawDB(t, ctx)
	)

	populateDatabase(t, database)

	t.Run("schema_migrations", func(t *testing.T) {
		checkSchemaMigrations(t, database)
	})

	t.Run("00001_init-schema", func(t *testing.T) {
		checkInitSchemaMigration(t, database)
	})

	t.Run("00002_partition_management", func(t *testing.T) {
		checkPartitionManagement(t, database)
	})

	t.Run("00003_add-latest-block", func(t *testing.T) {
		checkLatestBlock(t, database)
	})

	t.Run("00004_add_blockchain_columns", func(t *testing.T) {
		checkBlockchainColumns(t, database)
	})

	t.Run("00005_gateway_indexes", func(t *testing.T) {
		checkGatewayIndexes(t, database)
	})

	t.Run("00006_add_latest_envelopes", func(t *testing.T) {
		checkLatestEnvelopes(t, database)
	})

	t.Run("00007_unsettled-usage", func(t *testing.T) {
		checkUnsettledUsage(t, database)
	})

	t.Run("00008_payer-nonces", func(t *testing.T) {
		checkPayerNonces(t, database)
	})

	t.Run("00009_originator-congestion", func(t *testing.T) {
		checkOriginatorCongestion(t, database)
	})

	t.Run("00010_store-payer-reports", func(t *testing.T) {
		checkPayerReports(t, database)
	})

	t.Run("00011_add-migration-tracker", func(t *testing.T) {
		checkMigrationTracker(t, database)
	})

	t.Run("00012_payer_ledger_events", func(t *testing.T) {
		checkPayerLedgerEvents(t, database)
	})

	t.Run("00013_add-commit-messages-migration", func(t *testing.T) {
		checkCommitMessagesMigration(t, database)
	})

	t.Run("00014_add-dead-letter-box", func(t *testing.T) {
		checkDeadLetterBox(t, database)
	})

	t.Run("00015_partition_management_v2", func(t *testing.T) {
		checkPartitionManagementV2(t, database)
	})

	t.Run("00016_insert-gateway-envelopes-batch", func(t *testing.T) {
		checkInsertBatch(t, database)
	})

	t.Run("00017_payer_id-foreign-key", func(t *testing.T) {
		checkPayerForeignKeys(t, database)
	})

	t.Run("00018_add_latest_envelopes_v2", func(t *testing.T) {
		checkLatestEnvelopesV2(t, database)
	})

	t.Run("00019_v3b_indexes", func(t *testing.T) {
		checkV3bIndexes(t, database)
	})

	t.Run("data_verification", func(t *testing.T) {
		checkDataVerification(t, database)
	})
}

/* Per migration checks */

func checkSchemaMigrations(t *testing.T, database *sql.DB) {
	row := database.QueryRowContext(
		t.Context(),
		"SELECT * FROM schema_migrations",
	)

	var (
		version int64
		dirty   bool
	)

	err := row.Scan(&version, &dirty)
	require.NoError(t, err)
	require.Equal(t, currentMigration, version)
	require.False(t, dirty)
}

func checkInitSchemaMigration(t *testing.T, database *sql.DB) {
	tables := []string{
		"node_info",
		"staged_originator_envelopes",
		"address_log",
		"payers",
		"gateway_envelopes_meta",
		"gateway_envelope_blobs",
	}
	for _, tbl := range tables {
		tableExists(t, database, tbl)
	}

	viewExists(t, database, "gateway_envelopes_view")
	functionExists(t, database, "insert_staged_originator_envelope")
}

func checkPartitionManagement(t *testing.T, database *sql.DB) {
	functions := []string{
		"make_meta_originator_part",
		"make_meta_seq_subpart",
		"make_blob_originator_part",
		"make_blob_seq_subpart",
		"ensure_gateway_parts",
	}
	for _, fn := range functions {
		functionExists(t, database, fn)
	}
}

func checkLatestBlock(t *testing.T, database *sql.DB) {
	tableExists(t, database, "latest_block")
}

func checkBlockchainColumns(t *testing.T, database *sql.DB) {
	tableExists(t, database, "blockchain_messages")
	indexExists(t, database, "idx_blockchain_messages_block_canonical")
}

func checkGatewayIndexes(t *testing.T, database *sql.DB) {
	// gem_time_node_seq_idx, gem_topic_time_idx, and gem_originator_node_id
	// are dropped by migration 19 (v3b_indexes).
	indexes := []string{
		"gem_topic_time_desc_idx",
		"gem_expiry_idx",
	}
	for _, idx := range indexes {
		indexExists(t, database, idx)
	}
}

func checkLatestEnvelopes(t *testing.T, database *sql.DB) {
	tableExists(t, database, "gateway_envelopes_latest")
}

func checkUnsettledUsage(t *testing.T, database *sql.DB) {
	tableExists(t, database, "unsettled_usage")
	indexExists(t, database, "idx_unsettled_usage_originator_id_minutes_since_epoch")
}

func checkPayerNonces(t *testing.T, database *sql.DB) {
	tableExists(t, database, "nonce_table")
	functionExists(t, database, "fill_nonce_gap")
}

func checkOriginatorCongestion(t *testing.T, database *sql.DB) {
	tableExists(t, database, "originator_congestion")
}

func checkPayerReports(t *testing.T, database *sql.DB) {
	tableExists(t, database, "payer_reports")
	tableExists(t, database, "payer_report_attestations")

	indexes := []string{
		"payer_reports_submission_status_created_idx",
		"payer_reports_attestation_status_created_idx",
		"payer_report_attestations_payer_report_id_idx",
	}
	for _, idx := range indexes {
		indexExists(t, database, idx)
	}
}

func checkMigrationTracker(t *testing.T, database *sql.DB) {
	tableExists(t, database, "migration_tracker")

	// Verify initial rows inserted by migration 11.
	expectedTables := []string{
		"group_messages",
		"inbox_log",
		"key_packages",
		"welcome_messages",
	}
	for _, tbl := range expectedTables {
		var exists bool
		err := database.QueryRowContext(
			t.Context(),
			`SELECT EXISTS (
				SELECT 1 FROM migration_tracker WHERE source_table = $1
			)`,
			tbl,
		).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "migration_tracker should have row for %s", tbl)
	}
}

func checkPayerLedgerEvents(t *testing.T, database *sql.DB) {
	tableExists(t, database, "payer_ledger_events")
	indexExists(t, database, "idx_payer_ledger_events_payer_id")
}

func checkCommitMessagesMigration(t *testing.T, database *sql.DB) {
	var exists bool
	err := database.QueryRowContext(
		t.Context(),
		`SELECT EXISTS (
			SELECT 1 FROM migration_tracker WHERE source_table = 'commit_messages'
		)`,
	).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "migration_tracker should have row for commit_messages")
}

func checkDeadLetterBox(t *testing.T, database *sql.DB) {
	tableExists(t, database, "migration_dead_letter_box")
	indexExists(t, database, "migration_dead_letter_box_source_table_added_at_idx")
	indexExists(t, database, "migration_dead_letter_box_retryable_retried_at_idx")
	functionExists(t, database, "insert_migration_dead_letter_box")
	functionExists(t, database, "delete_migration_dead_letter_box")
}

func checkPartitionManagementV2(t *testing.T, database *sql.DB) {
	functions := []string{
		"make_meta_originator_part_v2",
		"make_meta_seq_subpart_v2",
		"make_blob_originator_part_v2",
		"make_blob_seq_subpart_v2",
		"ensure_gateway_parts_v2",
	}
	for _, fn := range functions {
		functionExists(t, database, fn)
	}
}

func checkInsertBatch(t *testing.T, database *sql.DB) {
	functionExists(t, database, "insert_gateway_envelope_batch")
}

func checkPayerForeignKeys(t *testing.T, database *sql.DB) {
	constraintExists(t, database, "fk_unsettled_usage_payer_id")
	constraintExists(t, database, "fk_payer_ledger_events_payer_id")
}

func checkLatestEnvelopesV2(t *testing.T, database *sql.DB) {
	functionExists(t, database, "update_latest_envelope_v2")
	triggerExists(t, database, "gateway_latest_upd_v2")
}

func checkV3bIndexes(t *testing.T, database *sql.DB) {
	indexExists(t, database, "gem_topic_orig_seq_idx")
}

// --- Data verification after populateDatabase ---

func checkDataVerification(t *testing.T, database *sql.DB) {
	ctx := t.Context()

	t.Run("node_info_rows", func(t *testing.T) {
		var count int
		err := database.QueryRowContext(ctx, "SELECT COUNT(*) FROM node_info").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "node_info should have 1 row (singleton constraint)")
	})

	t.Run("payers_rows", func(t *testing.T) {
		var count int
		err := database.QueryRowContext(ctx, "SELECT COUNT(*) FROM payers").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 3, count, "payers should have 3 rows")
	})

	t.Run("gateway_envelopes_meta_rows", func(t *testing.T) {
		var count int
		err := database.QueryRowContext(ctx, "SELECT COUNT(*) FROM gateway_envelopes_meta").
			Scan(&count)
		require.NoError(t, err)
		assert.Positive(t, count, "gateway_envelopes_meta should have rows")
	})

	t.Run("gateway_envelopes_view_returns_data", func(t *testing.T) {
		var count int
		err := database.QueryRowContext(ctx, "SELECT COUNT(*) FROM gateway_envelopes_view").
			Scan(&count)
		require.NoError(t, err)
		assert.Positive(t, count, "gateway_envelopes_view should return joined data")
	})
}

/* Helpers for catalog queries */

func tableExists(t *testing.T, database *sql.DB, tableName string) {
	t.Helper()
	var exists bool
	err := database.QueryRowContext(
		t.Context(),
		`SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = $1
		)`,
		tableName,
	).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "table %s should exist", tableName)
}

func indexExists(t *testing.T, database *sql.DB, indexName string) {
	t.Helper()
	var exists bool
	err := database.QueryRowContext(
		t.Context(),
		`SELECT EXISTS (
			SELECT 1 FROM pg_indexes
			WHERE schemaname = 'public' AND indexname = $1
		)`,
		indexName,
	).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "index %s should exist", indexName)
}

func functionExists(t *testing.T, database *sql.DB, funcName string) {
	t.Helper()
	var exists bool
	err := database.QueryRowContext(
		t.Context(),
		`SELECT EXISTS (
			SELECT 1 FROM pg_proc p
			JOIN pg_namespace n ON p.pronamespace = n.oid
			WHERE n.nspname = 'public' AND p.proname = $1
		)`,
		funcName,
	).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "function %s should exist", funcName)
}

func triggerExists(t *testing.T, database *sql.DB, triggerName string) {
	t.Helper()
	var exists bool
	err := database.QueryRowContext(
		t.Context(),
		`SELECT EXISTS (
			SELECT 1 FROM information_schema.triggers
			WHERE trigger_schema = 'public' AND trigger_name = $1
		)`,
		triggerName,
	).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "trigger %s should exist", triggerName)
}

func viewExists(t *testing.T, database *sql.DB, viewName string) {
	t.Helper()
	var exists bool
	err := database.QueryRowContext(
		t.Context(),
		`SELECT EXISTS (
			SELECT 1 FROM information_schema.views
			WHERE table_schema = 'public' AND table_name = $1
		)`,
		viewName,
	).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "view %s should exist", viewName)
}

func constraintExists(t *testing.T, database *sql.DB, constraintName string) {
	t.Helper()
	var exists bool
	err := database.QueryRowContext(
		t.Context(),
		`SELECT EXISTS (
			SELECT 1 FROM information_schema.table_constraints
			WHERE constraint_schema = 'public' AND constraint_name = $1
		)`,
		constraintName,
	).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "constraint %s should exist", constraintName)
}
