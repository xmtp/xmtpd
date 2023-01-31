package crdt

// NodeBroadcaster manages the overall broadcasting capacity of a Node
type NodeBroadcaster interface {
	// NewTopic creates a broadcaster for a specific topic
	NewTopic(name string, node *Node) TopicBroadcaster
}

// TopicBroadcaster manages broadcasts for a given topic
type TopicBroadcaster interface {
	// Broadcast sends an Event out to the network
	Broadcast(*Event)
}
