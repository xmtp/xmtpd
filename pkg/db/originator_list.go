package db

import (
	"context"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/db/queries"
)

// OriginatorLister returns a list of all known originator node IDs.
type OriginatorLister interface {
	GetOriginatorNodeIDs(ctx context.Context) ([]int32, error)
}

// CachedOriginatorList queries gateway_envelopes_latest for originator IDs
// and caches the result for a configurable TTL.
type CachedOriginatorList struct {
	querier *queries.Queries
	ttl     time.Duration

	mu        sync.RWMutex
	cached    []int32
	fetchedAt time.Time
}

func NewCachedOriginatorList(
	querier *queries.Queries,
	ttl time.Duration,
) *CachedOriginatorList {
	return &CachedOriginatorList{
		querier: querier,
		ttl:     ttl,
	}
}

func (c *CachedOriginatorList) GetOriginatorNodeIDs(
	ctx context.Context,
) ([]int32, error) {
	c.mu.RLock()
	if !c.fetchedAt.IsZero() && time.Since(c.fetchedAt) < c.ttl {
		result := cloneInt32s(c.cached)
		c.mu.RUnlock()
		return result, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock.
	if !c.fetchedAt.IsZero() && time.Since(c.fetchedAt) < c.ttl {
		return cloneInt32s(c.cached), nil
	}

	ids, err := c.querier.SelectOriginatorNodeIDs(ctx)
	if err != nil {
		return nil, err
	}

	c.cached = ids
	c.fetchedAt = time.Now()
	return cloneInt32s(ids), nil
}

func cloneInt32s(s []int32) []int32 {
	if s == nil {
		return nil
	}
	out := make([]int32, len(s))
	copy(out, s)
	return out
}
