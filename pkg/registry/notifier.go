package registry

import (
	"sync"
)

type notifier[ValueType any] struct {
	channels map[chan<- ValueType]bool
	mutex    sync.RWMutex
}

func newNotifier[ValueType any]() *notifier[ValueType] {
	return &notifier[ValueType]{
		channels: make(map[chan<- ValueType]bool),
	}
}

func (c *notifier[Node]) register() (<-chan Node, CancelSubscription) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	newChannel := make(chan Node)
	c.channels[newChannel] = true

	return newChannel, func() {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		close(newChannel)
		delete(c.channels, newChannel)
	}
}

func (c *notifier[any]) trigger(value any) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for channel := range c.channels {
		// Write to the channel in a goroutine to avoid blocking the caller
		go func(channel chan<- any) {
			channel <- value
		}(channel)
	}
}
