package networkwatcher

import (
	"context"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api/metadata_apiconnect"
)

// fakeMetadataServer implements SubscribeSyncCursor for tests.
type fakeMetadataServer struct {
	metadata_apiconnect.UnimplementedMetadataApiHandler
	sends []map[uint32]uint64
	// hold blocks the handler until the channel is closed, so the test can
	// keep the stream open.
	hold chan struct{}
}

func (f *fakeMetadataServer) SubscribeSyncCursor(
	ctx context.Context,
	_ *connect.Request[metadata_api.GetSyncCursorRequest],
	stream *connect.ServerStream[metadata_api.GetSyncCursorResponse],
) error {
	for _, snap := range f.sends {
		cursor := &envelopes.Cursor{NodeIdToSequenceId: make(map[uint32]uint64, len(snap))}
		maps.Copy(cursor.NodeIdToSequenceId, snap)
		if err := stream.Send(&metadata_api.GetSyncCursorResponse{LatestSync: cursor}); err != nil {
			return err
		}
	}
	if f.hold != nil {
		select {
		case <-f.hold:
		case <-ctx.Done():
		}
	}
	return nil
}

func newTestServer(t *testing.T, handler metadata_apiconnect.MetadataApiHandler) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	path, h := metadata_apiconnect.NewMetadataApiHandler(handler)
	mux.Handle(path, h)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestSubscriber_StreamsCursorsIntoAggregator(t *testing.T) {
	resetAggregatorMetrics()

	hold := make(chan struct{})
	defer close(hold)

	srv := newTestServer(t, &fakeMetadataServer{
		sends: []map[uint32]uint64{
			{100: 10, 200: 20},
			{100: 11, 200: 21},
		},
		hold: hold,
	})

	a := NewAggregator()
	sub := NewSubscriber(SubscriberConfig{
		NodeID:     42,
		BaseURL:    srv.URL,
		Aggregator: a,
		Logger:     zap.NewNop(),
		MinBackoff: 10 * time.Millisecond,
		MaxBackoff: 50 * time.Millisecond,
		HTTPClient: srv.Client(),
	})

	ctx := t.Context()

	go sub.Run(ctx)

	require.Eventually(t, func() bool {
		return metricValue(t, cursorGauge.WithLabelValues("42", "100")) == 11 &&
			metricValue(t, cursorGauge.WithLabelValues("42", "200")) == 21 &&
			metricValue(t, nodeUp.WithLabelValues("42")) == 1
	}, 2*time.Second, 10*time.Millisecond, "expected aggregator to receive both cursor snapshots")
}
