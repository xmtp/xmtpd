package bolt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/context"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

/* This was used to benchmark the EVENTS bucket arrangements by CID and ByTimeKey:
   To produce benchmark result run something like
		go test -run XXX -bench=. -count=10 | tee events-by-time.bench
   To compare multiple benchmark results use benchstat
		benchstat events-by-cid.bench events-by-time.bench
*/

func BenchmarkQuery(b *testing.B) {
	ctx := test.NewContext(b)
	store := newTestNodeStore(b, ctx)
	defer store.Close()
	topic := newTestStore(b, "topic", store)
	defer topic.Close()
	err := topic.seed(100000)
	require.NoError(b, err)
	b.Log("seeded")

	runQueryBenchmark(b, ctx, topic, 1000, 10000, 10)
	runQueryBenchmark(b, ctx, topic, 1000, 99000, 10)
	runQueryBenchmark(b, ctx, topic, 1000, 10000, 100)
	runQueryBenchmark(b, ctx, topic, 1000, 99000, 100)
	runQueryBenchmark(b, ctx, topic, 1000, 10000, 1000)
	runQueryBenchmark(b, ctx, topic, 1000, 99000, 1000)
}

func runQueryBenchmark(b *testing.B, ctx context.Context, topic *testStore, start, end uint64, pageSize uint32) {
	b.Run(fmt.Sprintf("%d/%d/%d", start, end, pageSize), func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var cursor *messagev1.Cursor
			var resp *messagev1.QueryResponse
			var err error
			for count := end - start + 1; count > 0; count -= uint64(len(resp.Envelopes)) {
				resp, err = topic.Query(ctx, &messagev1.QueryRequest{
					StartTimeNs: start,
					EndTimeNs:   end,
					PagingInfo: &messagev1.PagingInfo{
						Limit:  pageSize,
						Cursor: cursor,
					},
				})
				require.NoError(b, err)
				cursor = resp.PagingInfo.Cursor
			}
		}
	})
}
