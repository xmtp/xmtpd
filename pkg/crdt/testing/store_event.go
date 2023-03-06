package crdttest

import (
	"testing"

	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
)

func RunStoreEventTests(t *testing.T, topic string, storeMaker TestStoreMaker) {
	t.Run("insert events", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		defer s.Close()

		ev1 := s.addRandomEvent(t, topic)
		s.requireEventsEqual(t, ev1)

		ev2 := s.addRandomEvent(t, topic)
		s.requireEventsEqual(t, ev1, ev2)
	})

	t.Run("insert existing event", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		defer s.Close()

		ev1 := s.addRandomEvent(t, topic)
		s.requireEventsEqual(t, ev1)

		s.addExistingEvent(t, ev1)
		s.requireEventsEqual(t, ev1)
	})

	t.Run("append events", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		defer s.Close()

		ev1 := s.appendRandomEvent(t, topic)
		require.NotNil(t, ev1.Cid)
		require.Empty(t, ev1.Links)
		s.requireEventsEqual(t, ev1)

		ev2 := s.appendRandomEvent(t, topic)
		require.NotNil(t, ev2.Cid)
		require.Equal(t, []multihash.Multihash{ev1.Cid}, ev2.Links)
		s.requireEventsEqual(t, ev1, ev2)
	})

	t.Run("insert remove heads", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		defer s.Close()

		head := s.addRandomHead(t, topic)
		s.requireEventsEqual(t, head)

		ev1 := s.newRandomEventWithHeads(t, topic, []multihash.Multihash{head.Cid})
		require.Equal(t, []multihash.Multihash{head.Cid}, ev1.Links)
		s.addEvent(t, ev1)
		s.requireEventsEqual(t, head, ev1)

		ev2 := s.newRandomEvent(t, topic)
		require.Nil(t, ev2.Links)
		s.addEvent(t, ev2)
		s.requireEventsEqual(t, head, ev1, ev2)
	})

	t.Run("insert existing head", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		defer s.Close()

		head := s.addRandomHead(t, topic)
		s.requireEventsEqual(t, head)

		s.addExistingEvent(t, head)
		s.requireEventsEqual(t, head)
	})

	t.Run("find missing links", func(t *testing.T) {
		t.Parallel()
		s := storeMaker(t)
		defer s.Close()

		missingEv1 := s.newRandomEvent(t, topic)
		s.requireNoEvents(t)
		require.NotNil(t, missingEv1)

		missingEv2 := s.newRandomEvent(t, topic)
		s.requireNoEvents(t)
		require.NotNil(t, missingEv2)

		missingEv3 := s.newRandomEvent(t, topic)
		s.requireNoEvents(t)
		require.NotNil(t, missingEv3)

		ev1 := s.newRandomEvent(t, topic)
		s.requireNoEvents(t)
		ev1.Links = []multihash.Multihash{missingEv1.Cid, missingEv2.Cid}

		s.addEvent(t, ev1)
		s.requireEventsEqual(t, ev1)

		ev2 := s.newRandomEvent(t, topic)
		s.requireEventsEqual(t, ev1)
		ev2.Links = []multihash.Multihash{missingEv3.Cid}

		s.addEvent(t, ev2)
		s.requireEventsEqual(t, ev1, ev2)

		cids, err := s.FindMissingLinks(s.ctx)
		require.NoError(t, err)
		require.ElementsMatch(t, []multihash.Multihash{missingEv1.Cid, missingEv2.Cid, missingEv3.Cid}, cids)
	})
}
