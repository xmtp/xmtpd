package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func RandomMsgTest(t *testing.T, nodes, topics, messages int) *network {
	// to emulate significant concurrent activity we want nodes to be adding
	// events concurrently, but we also want to allow propagation at the same time.
	// So we need to introduce short delays to allow the network
	// to make some propagation progress. Given the random spraying approach
	// injecting a delay at every (nodes*topics)th event should allow most nodes
	// to inject an event to most topics, and then the random length of the delay
	// should allow some amount of propagation to happen before the next burst.
	delayEvery := nodes * topics
	net := NewNetwork(t, nodes, topics)
	for i := 0; i < messages; i++ {
		topic := fmt.Sprintf("t%d", rand.Intn(topics))
		msg := fmt.Sprintf("gm %d", i)
		net.Publish(t, rand.Intn(nodes), topic, msg)
		if i%delayEvery == 0 {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
		}
	}
	net.AssertEventuallyConsistent(t, time.Duration(messages*nodes)*time.Millisecond)
	return net
}
