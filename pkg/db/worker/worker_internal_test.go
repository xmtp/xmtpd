package worker

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

func TestWorker_CheckPartition(t *testing.T) {
	t.Run("empty partition", func(t *testing.T) {
		ts := newTestScaffold(t)

		err := ts.worker.runPartitionCheck(t.Context(), 100, 0)
		require.NoError(t, err)
		require.Empty(t, readPartitionSet(t, ts.db))
	})
	t.Run("partitions with enough room", func(t *testing.T) {
		ts := newTestScaffold(t)

		envelopes := map[uint32]uint64{
			100: 1,
			200: 100_000,
			300: 500_000,
		}

		testutils.InsertGatewayEnvelopes(t, ts.db.DB(), generateEnvelopes(t, envelopes))
		require.Len(t, readPartitionSet(t, ts.db), len(envelopes))

		err := ts.worker.runDBCheck(t.Context())
		require.NoError(t, err)

		// No new partitions should be created.
		require.Len(t, readPartitionSet(t, ts.db), len(envelopes))
	})
	t.Run("checks create new partitions", func(t *testing.T) {
		ts := newTestScaffold(t)

		require.Empty(t, readPartitionSet(t, ts.db))

		envelopes := map[uint32]uint64{
			100: 1,
			200: 900_000,
			300: 750_000,
		}

		testutils.InsertGatewayEnvelopes(t, ts.db.DB(), generateEnvelopes(t, envelopes))
		require.Len(t, readPartitionSet(t, ts.db), len(envelopes))

		err := ts.worker.runDBCheck(t.Context())
		require.NoError(t, err)

		// Two new partitions should be created.
		partitions := readPartitionSet(t, ts.db)
		require.Len(t, partitions, 5)

		_, ok := partitions["gateway_envelopes_meta_o100_s1000000_2000000"]
		require.False(t, ok)

		_, ok = partitions["gateway_envelopes_meta_o200_s1000000_2000000"]
		require.True(t, ok)

		_, ok = partitions["gateway_envelopes_meta_o300_s1000000_2000000"]
		require.True(t, ok)
	})
	t.Run("checks create new higher partitions", func(t *testing.T) {
		ts := newTestScaffold(t)

		envelopes := map[uint32]uint64{
			100: 50,
			200: 3_750_000,
		}

		testutils.InsertGatewayEnvelopes(t, ts.db.DB(), generateEnvelopes(t, envelopes))
		require.Len(t, readPartitionSet(t, ts.db), len(envelopes))

		err := ts.worker.runDBCheck(t.Context())
		require.NoError(t, err)

		// One new partition should be created.
		partitions := readPartitionSet(t, ts.db)
		require.Len(t, partitions, 3)

		_, ok := partitions["gateway_envelopes_meta_o100_s1000000_2000000"]
		require.False(t, ok)

		_, ok = partitions["gateway_envelopes_meta_o200_s3000000_4000000"]
		require.True(t, ok)
	})
	t.Run("trying to create the same partition multiple times does nothing", func(t *testing.T) {
		ts := newTestScaffold(t)

		envelopes := map[uint32]uint64{
			100: 4_900_000,
		}

		testutils.InsertGatewayEnvelopes(t, ts.db.DB(), generateEnvelopes(t, envelopes))
		require.Len(t, readPartitionSet(t, ts.db), 1)

		err := ts.worker.runDBCheck(t.Context())
		require.NoError(t, err)

		// One new partition should be created.
		partitions := readPartitionSet(t, ts.db)
		require.Len(t, partitions, 2)

		_, ok := partitions["gateway_envelopes_meta_o100_s4000000_5000000"]
		require.True(t, ok)

		_, ok = partitions["gateway_envelopes_meta_o100_s5000000_6000000"]
		require.True(t, ok)

		err = ts.worker.runDBCheck(t.Context())
		require.NoError(t, err)

		// No new partitions should be created.
		partitions = readPartitionSet(t, ts.db)
		require.Len(t, partitions, 2)
	})
}

func readPartitionSet(t *testing.T, db *db.Handler) map[string]struct{} {
	t.Helper()

	query := `SELECT table_name
			  	FROM information_schema.tables
				WHERE table_name ~ '^gateway_envelopes_meta_o\d+_s\d+_\d+$'`

	rows, err := db.DB().QueryContext(t.Context(), query)
	require.NoError(t, err)

	defer rows.Close()
	out := make(map[string]struct{})
	for rows.Next() {

		var tableName string
		err := rows.Scan(&tableName)
		require.NoError(t, err)

		out[tableName] = struct{}{}
	}

	err = rows.Close()
	require.NoError(t, err)

	err = rows.Err()
	require.NoError(t, err)

	return out
}

func generateEnvelopes(
	t *testing.T,
	seqIDs map[uint32]uint64,
) []queries.InsertGatewayEnvelopeParams {
	t.Helper()

	topic := topic.NewTopic(
		topic.TopicKindGroupMessagesV1,
		[]byte(fmt.Sprintf("generic-topic-%v", rand.Int())),
	)

	out := make([]queries.InsertGatewayEnvelopeParams, 0, len(seqIDs))

	for nodeID, seqID := range seqIDs {

		oe := testutils.Marshal(
			t,
			envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(
				t,
				uint32(nodeID),
				uint64(seqID),
				topic.Bytes(),
			),
		)

		out = append(out, queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     int32(nodeID),
			OriginatorSequenceID: int64(seqID),
			Topic:                topic.Bytes(),
			OriginatorEnvelope:   oe,
		})

	}

	return out
}

type testScaffold struct {
	db     *db.Handler
	worker *Worker
}

func newTestScaffold(t *testing.T) testScaffold {
	t.Helper()

	var (
		db, _ = testutils.NewDB(t, t.Context())
		log   = testutils.NewLog(t)
	)

	worker := NewWorker(log, db)

	ts := testScaffold{
		db:     db,
		worker: worker,
	}

	return ts
}
