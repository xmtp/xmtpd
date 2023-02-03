package tests

import (
	"flag"
	"fmt"
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
	net.AddNode(t, newMapStore())
	// need to trigger a sync for the node to catch up
	net.Publish(t, 0, t0, "ahoy new node")
	net.AssertEventuallyConsistent(t, time.Second)
}

func Test_NodeRestart(t *testing.T) {
	// create a network with some pre-existing traffic
	net := randomMsgTest(t, 3, 1, 10)
	defer net.Close()
	// replace node 2 reusing its store
	n := net.RemoveNode(t, 2)
	store := n.NodeStore.(*mapStore)
	// delete some early events from the node store
	// to see if they get re-fetched during bootstrap
	for _, ev := range net.events[:5] {
		delete(store.topics[t0].events, ev.Cid.String())
	}
	net.AddNode(t, store)
	net.AssertEventuallyConsistent(t, time.Second)
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
