package message

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunnel(t *testing.T) {
	sendFn := func(s []int, ch chan int) {
		for _, n := range s {
			ch <- n
		}
		close(ch)
	}

	t.Run("channel merge", func(t *testing.T) {
		var (
			a = []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
			b = []int{20, 21, 22, 23, 24}

			chA = make(chan int)
			chB = make(chan int)
		)

		ch := newFunnel(chA, chB).output()

		// Start goroutines to send the individual elements and then close respective channels.
		go sendFn(a, chA)
		go sendFn(b, chB)

		received := make(map[int]struct{})
		for n := range ch {
			received[n] = struct{}{}
		}

		// Make sure all elements were received.
		require.Len(t, received, len(a)+len(b))

		// Make sure all individual elements from input channels were seen.
		full := slices.Concat(a, b)
		for _, x := range full {
			_, seen := received[x]
			require.True(t, seen)
		}
	})
	t.Run("adding a channel", func(t *testing.T) {
		var (
			a = []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
			b = []int{20, 21, 22, 23, 24}
			c = []int{30, 31, 32, 33, 34, 35, 36}

			chA = make(chan int)
			chB = make(chan int)
			chC = make(chan int)

			funnel = newFunnel(chA, chB)
		)

		// Start goroutines to send the individual elements and then close respective channels.
		go sendFn(a, chA)
		go sendFn(b, chB)

		// Existing channel should receive all data, even if handle was retrieved prior to adding more channels.
		ch := funnel.output()

		go sendFn(c, chC)
		funnel.addChannel(chC)

		received := make(map[int]struct{})
		for n := range ch {
			received[n] = struct{}{}
		}

		// Make sure all elements were received.
		require.Len(t, received, len(a)+len(b)+len(c))

		// Make sure all individual elements from input channels were seen.
		full := slices.Concat(a, b, c)
		for _, x := range full {
			_, seen := received[x]
			require.True(t, seen)
		}
	})
}
