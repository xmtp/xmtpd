package tests

import (
	"testing"
	"time"
)

func Test_BasicBroadcast(t *testing.T) {
	net := NewNetwork(t, 5, 1)
	defer net.Close()
	net.Publish(t, 0, t0, "hi")
	net.AssertEventuallyConsistent(t, time.Second)
}
