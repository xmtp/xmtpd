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
	getCurrentCount := CountChannel(channel)

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
	getCurrentCount1 := CountChannel(channel1)
	channel2 := registry.register()
	getCurrentCount2 := CountChannel(channel2)

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
	getCurrentCount := CountChannel(channel)

	for i := 0; i < 100; i++ {
		go registry.trigger(1)
	}

	require.Eventually(t, func() bool {
		return getCurrentCount() == 100
	}, time.Second, 10*time.Millisecond)
}

func CountChannel[Kind any](ch <-chan Kind, validators ...func(Kind)) func() int {
	var count int
	var mutex sync.RWMutex
	go func() {
		for v := range ch {
			for _, validate := range validators {
				if validate != nil {
					validate(v)
				}
			}
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
