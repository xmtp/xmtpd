package registry

import (
	"sync"
)

type SingleNotificationNotifier[ValueType any] struct {
	channels map[chan<- ValueType]bool
	mutex    sync.RWMutex
}

func newNotifier[ValueType any]() *SingleNotificationNotifier[ValueType] {
	return &SingleNotificationNotifier[ValueType]{
		channels: make(map[chan<- ValueType]bool),
	}
}

func (c *SingleNotificationNotifier[ValueType]) register() <-chan ValueType {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	newChannel := make(chan ValueType, 1)
	c.channels[newChannel] = true

	return newChannel
}

func (c *SingleNotificationNotifier[ValueType]) trigger(value ValueType) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for channel := range c.channels {
		channel <- value
	}
}
