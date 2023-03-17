package node_test

import (
	"testing"

	ntest "github.com/xmtp/xmtpd/pkg/node/testing"
)

func TestNetwork(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow under race detector")
	}
	net := ntest.NewNetwork(t, 3)
	defer net.Close()

	sub := net.Subscribe(t, "topic1")
	envs := net.PublishRandom(t, sub.Topic, 1)
	sub.RequireEventuallyCapturedEvents(t, envs)
	net.RequireEventuallyStoredEvents(t, "topic1", envs)
}
