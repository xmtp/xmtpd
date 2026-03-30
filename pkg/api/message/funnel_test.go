package message

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFunnelBasic(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	f := newFunnel(ch1, ch2)

	var received []int
	var mu sync.Mutex
	done := make(chan struct{})

	go func() {
		for v := range f.output() {
			mu.Lock()
			received = append(received, v)
			mu.Unlock()
		}
		close(done)
	}()

	ch1 <- 1
	ch2 <- 2
	ch1 <- 3

	close(ch1)
	close(ch2)

	<-done

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, received, 3)
	require.ElementsMatch(t, []int{1, 2, 3}, received)
}

func TestFunnelAddChannelAfterAllClosed(t *testing.T) {
	ch1 := make(chan int)
	f := newFunnel(ch1)

	var received []int
	var mu sync.Mutex
	done := make(chan struct{})

	go func() {
		for v := range f.output() {
			mu.Lock()
			received = append(received, v)
			mu.Unlock()
		}
		close(done)
	}()

	ch1 <- 1
	close(ch1)

	// Give the closer time to detect ch1 closure and close f.out.
	time.Sleep(50 * time.Millisecond)

	// Add a new channel after the first one closed.
	// This must not panic (send on closed channel).
	ch2 := make(chan int, 1)
	f.addChannel(ch2)

	ch2 <- 2
	close(ch2)

	<-done

	mu.Lock()
	defer mu.Unlock()
	// ch2's value may or may not be delivered (output might already be closed).
	// The critical assertion: no panic occurred.
	require.GreaterOrEqual(t, len(received), 1, "should have received at least the value from ch1")
}
