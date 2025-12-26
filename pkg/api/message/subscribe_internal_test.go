package message

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInputChannelMerge(t *testing.T) {
	var (
		a = []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
		b = []int{20, 21, 22, 23, 24}

		chA = make(chan int)
		chB = make(chan int)
	)

	ch := merge(chA, chB)

	sendFn := func(s []int, ch chan int) {
		for _, n := range s {
			ch <- n
		}
		close(ch)
	}

	// Start goroutines to send the individual elements and then close respective channels.
	go sendFn(a, chA)
	go sendFn(b, chB)

	received := make(map[int]struct{})
	for n := range ch {
		t.Logf("received %v", n)
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
}
