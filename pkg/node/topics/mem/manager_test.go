package memtopics_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/crdt"
	memtopics "github.com/xmtp/xmtpd/pkg/node/topics/mem"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestMemoryTopics(t *testing.T) {
	ctx := context.Background()
	log := test.NewLogger(t)

	expectedReplica1 := &crdt.Replica{}
	expectedReplica2 := &crdt.Replica{}
	require.NotSame(t, expectedReplica1, expectedReplica2)

	topics, err := memtopics.New(log, func(topicId string) (*crdt.Replica, error) {
		switch topicId {
		case "topic1":
			return expectedReplica1, nil
		case "topic2":
			return expectedReplica2, nil
		}
		return nil, nil
	})
	require.NoError(t, err)
	defer topics.Close()

	replica, err := topics.GetOrCreateTopic(ctx, "topic1")
	require.NoError(t, err)
	require.Same(t, expectedReplica1, replica)

	replica, err = topics.GetOrCreateTopic(ctx, "topic2")
	require.NoError(t, err)
	require.Same(t, expectedReplica2, replica)
}
