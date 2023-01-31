package crdt

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

var visTopicM int // number of messages for VisualiseTopic test
var visTopicN int // number of nodes for VisualiseTopic test

func init() {
	flag.IntVar(&visTopicM, "visTopic", 0, "run VisualiseTopic test with specified number of messages")
	flag.IntVar(&visTopicN, "visTopicN", 3, "if running VisualiseTopic test, run with this many nodes")
}

type fixture struct {
	nodes    int
	topics   int
	messages int
}

func Test_RandomMessages(t *testing.T) {
	fixtures := []fixture{
		{5, 1, 100},
		{3, 3, 100},
		{10, 10, 1000},
		{10, 5, 10000},
	}
	if !testing.Short() {
		fixtures = append(fixtures, fixture{50, 1000, 100000}) // should take about 6s locally
	}
	for i, fix := range fixtures {
		t.Run(fmt.Sprintf("%d/%dn/%dt/%dm", i, fix.nodes, fix.topics, fix.messages),
			func(t *testing.T) { randomMsgTest(t, fix.nodes, fix.topics, fix.messages).Close() },
		)
	}
}

func Test_NewNodeJoin(t *testing.T) {
	// create a network with some pre-existing traffic
	net := randomMsgTest(t, 3, 1, 10)
	defer net.Close()
	// add a new node and observe that it catches up
	net.AddNode(newMapStore())
	// need to trigger a sync for the node to catch up
	net.Publish(0, t0, "ahoy new node")
	net.AssertEventuallyConsistent(time.Second)
}

func Test_NodeRestart(t *testing.T) {
	// create a network with some pre-existing traffic
	net := randomMsgTest(t, 3, 1, 10)
	defer net.Close()
	// replace node 2 reusing its store
	n := net.RemoveNode(2)
	store := n.NodeStore.(*mapStore)
	// delete some early events from the node store
	// to see if they get re-fetched during bootstrap
	for _, ev := range net.events[:5] {
		delete(store.topics[t0].events, ev.cid.String())
	}
	net.AddNode(store)
	net.AssertEventuallyConsistent(time.Second)
}

// Run a single topic test with given number of nodes and messages and visualise the resulting topic DAG after.
// Usage:
//
//		go test -visTopic=<messageCount> [visTopicN=<nodeCount>] >t0.dot
//	 	dot -Tjpg t0.dot >t0.jpg
//	 	open t0.jpg
func Test_VisualiseTopic(t *testing.T) {
	if visTopicM == 0 {
		return
	}
	net := randomMsgTest(t, visTopicN, 1, visTopicM)
	defer net.Close()
	net.visualiseTopic(os.Stdout, t0)
}

func randomMsgTest(t *testing.T, nodes, topics, messages int) *network {
	// to emulate significant concurrent activity we want nodes to be adding
	// events concurrently, but we also want to allow propagation at the same time.
	// So we need to introduce short delays to allow the network
	// to make some propagation progress. Given the random spraying approach
	// injecting a delay at every (nodes*topics)th event should allow most nodes
	// to inject an event to most topics, and then the random length of the delay
	// should allow some amount of propagation to happen before the next burst.
	delayEvery := nodes * topics
	net := newNetwork(t, nodes, topics)
	for i := 0; i < messages; i++ {
		topic := fmt.Sprintf("t%d", rand.Intn(topics))
		msg := fmt.Sprintf("gm %d", i)
		net.Publish(rand.Intn(nodes), topic, msg)
		if i%delayEvery == 0 {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
		}
	}
	net.AssertEventuallyConsistent(time.Duration(messages*nodes) * time.Millisecond)
	return net
}
