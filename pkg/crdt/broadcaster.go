package crdt

// NodeBroadcaster manages the overall broadcasting capacity of a Node
type NodeBroadcaster interface {
	// AddNode registers the node with the broadcaster.
	AddNode(*Node)

	// NewTopic creates a TopicBroadcaster for given topic and node.
	NewTopic(name string, node *Node) TopicBroadcaster
}

// TopicBroadcaster manages broadcasts for a given topic
type TopicBroadcaster interface {
	// Broadcast sends an Event out to the network
	Broadcast(*Event)
}
