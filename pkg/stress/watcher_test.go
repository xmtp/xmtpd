package stress

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

func TestProcessLogsStopsWhenChannelCloses(t *testing.T) {
	watcher := &Watcher{logger: zap.NewNop()}
	logs := make(chan types.Log)
	close(logs)

	select {
	case <-watcher.processLogs(context.Background(), logs):
	case <-time.After(time.Second):
		t.Fatal("processLogs did not stop after the log channel closed")
	}
}
