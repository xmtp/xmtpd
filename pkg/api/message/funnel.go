package message

import (
	"sync"
)

type funnel[T any] struct {
	*sync.Mutex

	in     []<-chan T
	out    chan T
	wg     sync.WaitGroup
	closer sync.Once
	closed bool
}

// Funnel many input channels into a single one. Output channel is closed once all input channels are.
// Adding channels after the output channel has been closed is invalid.
func newFunnel[T any](ch ...<-chan T) *funnel[T] {
	f := &funnel[T]{
		Mutex: &sync.Mutex{},
		in:    make([]<-chan T, 0, len(ch)),
		out:   make(chan T),
	}

	// Function will forward entries from its channel to the common channel.
	for _, c := range ch {
		f.addChannel(c)
	}

	return f
}

func (f *funnel[T]) addChannel(ch <-chan T) {
	f.Lock()
	defer f.Unlock()

	if f.closed {
		return
	}

	f.in = append(f.in, ch)

	f.wg.Add(1)

	f.startCloser()

	go func() {
		defer f.wg.Done()
		for e := range ch {
			f.out <- e
		}
	}()
}

func (f *funnel[T]) startCloser() {
	// Start closer goroutine (once).
	f.closer.Do(func() {
		go func() {
			// When all input channels are closed, close our channel too.
			f.wg.Wait()

			f.Lock()
			defer f.Unlock()

			close(f.out)
			f.closed = true
		}()
	})
}

func (f *funnel[T]) output() <-chan T {
	return f.out
}
