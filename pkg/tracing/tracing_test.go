package tracing

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GoPanicWrap_WaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	finished := false
	var finishedLock sync.RWMutex
	GoPanicWrap(ctx, &wg, "test", func(ctx context.Context) {
		<-ctx.Done()
		finishedLock.Lock()
		defer finishedLock.Unlock()
		finished = true
	})
	done := false
	var doneLock sync.RWMutex
	go func() {
		wg.Wait()
		doneLock.Lock()
		defer doneLock.Unlock()
		done = true
	}()
	go func() { time.Sleep(time.Millisecond); cancel() }()

	assert.Eventually(t, func() bool {
		finishedLock.RLock()
		defer finishedLock.RUnlock()
		doneLock.RLock()
		defer doneLock.RUnlock()
		return finished && done
	}, time.Second, 10*time.Millisecond)
}
