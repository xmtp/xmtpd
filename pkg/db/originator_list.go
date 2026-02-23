package db

import (
	"context"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"go.uber.org/zap"
)

// OriginatorLister returns a list of all known originator node IDs.
type OriginatorLister interface {
	GetOriginatorNodeIDs(ctx context.Context) ([]uint32, error)
}

// CachedOriginatorList queries gateway_envelopes_latest for originator IDs
// and caches the result for a configurable TTL.
type CachedOriginatorList struct {
	querier *queries.Queries
	ttl     time.Duration
	logger  *zap.Logger

	mu        sync.RWMutex
	cached    []uint32
	fetchedAt time.Time
}

func NewCachedOriginatorList(
	querier *queries.Queries,
	ttl time.Duration,
	logger *zap.Logger,
) *CachedOriginatorList {
	return &CachedOriginatorList{
		querier: querier,
		ttl:     ttl,
		logger:  logger,
	}
}

func (c *CachedOriginatorList) GetOriginatorNodeIDs(
	ctx context.Context,
) ([]uint32, error) {
	c.mu.RLock()
	if !c.fetchedAt.IsZero() && time.Since(c.fetchedAt) < c.ttl {
		result := cloneUint32s(c.cached)
		c.mu.RUnlock()
		return result, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock.
	if !c.fetchedAt.IsZero() && time.Since(c.fetchedAt) < c.ttl {
		return cloneUint32s(c.cached), nil
	}

	ids, err := c.querier.SelectOriginatorNodeIDs(ctx)
	if err != nil {
		return nil, err
	}

	c.cached = convertInt32sToUint32s(ids, c.logger)
	c.fetchedAt = time.Now()
	return cloneUint32s(c.cached), nil
}

func convertInt32sToUint32s(ids []int32, logger *zap.Logger) []uint32 {
	out := make([]uint32, 0, len(ids))
	for _, id := range ids {
		if id < 0 {
			logger.Warn(
				"skipping negative originator node ID",
				zap.Int32("originator_node_id", id),
			)
			continue
		}
		out = append(out, uint32(id))
	}
	return out
}

func cloneUint32s(s []uint32) []uint32 {
	if s == nil {
		return nil
	}
	out := make([]uint32, len(s))
	copy(out, s)
	return out
}
