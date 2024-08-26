package db

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type PollableDBQuery[ValueType any] func(ctx context.Context, lastSeenID int64, numRows int32) (results []ValueType, lastID int64, err error)

// Poll whenever notified, or at an interval if not notified
type PollingOptions struct {
	Interval time.Duration
	Notifier <-chan bool
	NumRows  int32
}

type DBSubscription[ValueType any] struct {
	ctx        context.Context
	log        *zap.Logger
	lastSeenID int64
	options    PollingOptions
	query      PollableDBQuery[ValueType]
	updates    chan<- []ValueType
}

func NewDBSubscription[ValueType any](
	ctx context.Context,
	log *zap.Logger,
	query PollableDBQuery[ValueType],
	lastSeenID int64,
	options PollingOptions,
) *DBSubscription[ValueType] {
	return &DBSubscription[ValueType]{
		ctx:        ctx,
		log:        log,
		lastSeenID: lastSeenID,
		options:    options,
		query:      query,
		updates:    nil,
	}
}

func (s *DBSubscription[ValueType]) Start() (<-chan []ValueType, error) {
	if s.updates != nil {
		return nil, fmt.Errorf("Already started")
	}
	if s.options.NumRows <= 0 || s.log == nil {
		return nil, fmt.Errorf("Required params not provided")
	}
	updates := make(chan []ValueType)
	s.updates = updates

	go func() {
		s.poll()

		timer := time.NewTimer(s.options.Interval)
		for {
			timer.Reset(s.options.Interval)
			select {
			case <-s.ctx.Done():
				s.log.Info("Context done; stopping subscription")
				close(s.updates)
				return
			case <-s.options.Notifier:
				s.poll()
			case <-timer.C:
				s.poll()
			}
		}
	}()

	return updates, nil
}

func (s *DBSubscription[ValueType]) poll() {
	// Repeatedly query page by page until no more results
	for {
		results, lastID, err := s.query(s.ctx, s.lastSeenID, s.options.NumRows)
		if err != nil {
			s.log.Error(
				"Error querying for DB subscription",
				zap.Error(err),
				zap.Int64("lastSeenID", s.lastSeenID),
				zap.Int32("numRows", s.options.NumRows),
			)
			// Did not update lastSeenID; will retry on next poll
			break
		}
		if len(results) == 0 {
			break
		}
		s.lastSeenID = lastID
		s.updates <- results
		if int32(len(results)) < s.options.NumRows {
			break
		}
	}
}
