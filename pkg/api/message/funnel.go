package message

import (
	"sync"
)

type funnel[T any] struct {
	mu     sync.Mutex
	in     []<-chan T
	out    chan T
	active int
	closed bool
}

// Funnel many input channels into a single one. Output channel is closed once all input channels are.
// Adding channels after the output channel has been closed is a no-op.
func newFunnel[T any](ch ...<-chan T) *funnel[T] {
	f := &funnel[T]{
		in:  make([]<-chan T, 0, len(ch)),
		out: make(chan T),
	}

	for _, c := range ch {
		f.addChannel(c)
	}

	return f
}

func (f *funnel[T]) addChannel(ch <-chan T) {
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return
	}

	f.in = append(f.in, ch)
	f.active++
	f.mu.Unlock()

	go func() {
		for e := range ch {
			f.out <- e
		}

		f.mu.Lock()
		f.active--
		if f.active == 0 {
			close(f.out)
			f.closed = true
		}
		f.mu.Unlock()
	}()
}

func (f *funnel[T]) output() <-chan T {
	return f.out
}
