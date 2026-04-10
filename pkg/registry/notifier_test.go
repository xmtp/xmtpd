package registry

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNotifier(t *testing.T) {
	registry := newNotifier[int]()
	channel := registry.register()
	getCurrentCount := CountChannel(t, channel)

	// Make sure the value is getting written to the channel
	registry.trigger(1)
	require.Eventually(t, func() bool {
		return getCurrentCount() == 1
	}, time.Second, 10*time.Millisecond)

	// Trigger again and make sure it still works
	registry.trigger(1)
	require.Eventually(t, func() bool {
		return getCurrentCount() == 2
	}, time.Second, 10*time.Millisecond)
}

func TestNotifierMultiple(t *testing.T) {
	registry := newNotifier[int]()

	channel1 := registry.register()
	getCurrentCount1 := CountChannel(t, channel1)
	channel2 := registry.register()
	getCurrentCount2 := CountChannel(t, channel2)

	registry.trigger(1)
	require.Eventually(t, func() bool {
		return getCurrentCount1() == 1 && getCurrentCount2() == 1
	}, time.Second, 10*time.Millisecond)

	registry.trigger(1)
	require.Eventually(t, func() bool {
		return getCurrentCount1() == 2 && getCurrentCount2() == 2
	}, time.Second, 10*time.Millisecond)
}

func TestNotifierConcurrent(t *testing.T) {
	registry := newNotifier[int]()
	channel := registry.register()
	getCurrentCount := CountChannel(t, channel)

	for range 100 {
		go registry.trigger(1)
	}

	require.Eventually(t, func() bool {
		return getCurrentCount() == 100
	}, 5*time.Second, 10*time.Millisecond)
}

// CountChannel spawns a reader goroutine that counts values received on ch.
// The goroutine exits when ch closes OR when the test completes (via t.Cleanup),
// preventing goroutine leaks when the notifier never closes its channels.
func CountChannel[Kind any](
	t *testing.T,
	ch <-chan Kind,
	validators ...func(Kind),
) func() int {
	t.Helper()

	var (
		count int
		mutex sync.RWMutex
		stop  = make(chan struct{})
		done  = make(chan struct{})
	)

	t.Cleanup(func() {
		close(stop)
		<-done
	})

	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			case v, ok := <-ch:
				if !ok {
					return
				}
				for _, validate := range validators {
					if validate != nil {
						validate(v)
					}
				}
				mutex.Lock()
				count++
				mutex.Unlock()
			}
		}
	}()

	return func() int {
		mutex.RLock()
		defer mutex.RUnlock()
		return count
	}
}
