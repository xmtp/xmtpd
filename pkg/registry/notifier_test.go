package registry

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNotifier(t *testing.T) {
	registry := newNotifier[int]()
	channel, cancel := registry.register()
	getCurrentCount := CountChannel(channel)

	// Make sure the value is getting writter to the cannel
	registry.trigger(1)
	// Sleep for 10ms since we read from the channel in a goroutinee
	time.Sleep(10 * time.Millisecond)
	require.Equal(t, 1, getCurrentCount())

	// Trigger again and make sure it still works
	registry.trigger(1)
	time.Sleep(10 * time.Millisecond)
	require.Equal(t, 2, getCurrentCount())

	// Unregister the subscription
	cancel()
	registry.trigger(1)
	time.Sleep(10 * time.Millisecond)
	// Make sure the count hasn't changed
	require.Equal(t, 2, getCurrentCount())
}

func TestNotifierMultiple(t *testing.T) {
	registry := newNotifier[int]()

	channel1, cancel1 := registry.register()
	getCurrentCount1 := CountChannel(channel1)
	channel2, cancel2 := registry.register()
	getCurrentCount2 := CountChannel(channel2)

	registry.trigger(1)
	time.Sleep(10 * time.Millisecond)
	require.Equal(t, 1, getCurrentCount1())
	require.Equal(t, 1, getCurrentCount2())

	cancel1()
	registry.trigger(1)
	time.Sleep(10 * time.Millisecond)
	require.Equal(t, 1, getCurrentCount1())
	require.Equal(t, 2, getCurrentCount2())
	cancel2()
}

func TestNotifierConcurrent(t *testing.T) {
	registry := newNotifier[int]()
	channel, cancel := registry.register()
	getCurrentCount := CountChannel(channel)
	defer cancel()

	for i := 0; i < 100; i++ {
		go registry.trigger(1)
	}
	time.Sleep(30 * time.Millisecond)
	require.Equal(t, 100, getCurrentCount())
}

func CountChannel[Kind any](ch <-chan Kind) func() int {
	var count int
	var mutex sync.RWMutex
	go func() {
		for range ch {
			mutex.Lock()
			count++
			mutex.Unlock()
		}
	}()

	return func() int {
		mutex.RLock()
		defer mutex.RUnlock()
		return count
	}
}
