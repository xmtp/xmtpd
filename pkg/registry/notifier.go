package registry

import (
	"sync"
)

type SingleNotificationNotifier[T any] struct {
	channels map[chan<- T]bool
	mutex    sync.RWMutex
}

func newNotifier[T any]() *SingleNotificationNotifier[T] {
	return &SingleNotificationNotifier[T]{
		channels: make(map[chan<- T]bool),
	}
}

func (c *SingleNotificationNotifier[T]) register() <-chan T {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ch := make(chan T, 1)
	c.channels[ch] = true

	return ch
}

func (c *SingleNotificationNotifier[T]) trigger(value T) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for channel := range c.channels {
		channel <- value
	}
}
