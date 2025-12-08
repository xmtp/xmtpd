package db

import (
	"context"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type PollableDBQuery[ValueType any, CursorType any] func(
	ctx context.Context,
	lastSeen CursorType,
	numRows int32,
) (results []ValueType, nextCursor CursorType, err error)

// PollingOptions specifies the polling options for a DB subscription.
// It can poll whenever notified, or at an interval if not notified.
type PollingOptions struct {
	Interval time.Duration
	Notifier <-chan bool
	NumRows  int32
}

// DBSubscription is a subscription that polls a DB for updates
// Assumes there is only one listener (updates block on a single unbuffered channel)
type DBSubscription[ValueType any, CursorType any] struct {
	ctx      context.Context
	logger   *zap.Logger
	lastSeen CursorType
	options  PollingOptions
	query    PollableDBQuery[ValueType, CursorType]
	updates  chan<- []ValueType
}

func NewDBSubscription[ValueType any, CursorType any](
	ctx context.Context,
	logger *zap.Logger,
	query PollableDBQuery[ValueType, CursorType],
	lastSeen CursorType,
	options PollingOptions,
) *DBSubscription[ValueType, CursorType] {
	logger = logger.Named(utils.DatabaseSubscriptionLoggerName)
	return &DBSubscription[ValueType, CursorType]{
		ctx:      ctx,
		logger:   logger,
		lastSeen: lastSeen,
		options:  options,
		query:    query,
		updates:  nil,
	}
}

func (s *DBSubscription[ValueType, CursorType]) Start() (<-chan []ValueType, error) {
	if s.updates != nil {
		return nil, fmt.Errorf("already started")
	}
	if s.options.NumRows <= 0 || s.logger == nil {
		return nil, fmt.Errorf("required params not provided")
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
				s.logger.Debug("context done; stopping")
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

func (s *DBSubscription[ValueType, CursorType]) poll() {
	// Repeatedly query page by page until no more results
	for {
		results, lastID, err := s.query(s.ctx, s.lastSeen, s.options.NumRows)
		if s.ctx.Err() != nil {
			break
		} else if err != nil {
			// Log is extremely noisy during test teardown
			s.logger.Error(
				"",
				zap.Error(err),
				utils.NumRowsField(s.options.NumRows),
			)
			// Did not update lastSeen; will retry on next poll
			break
		} else if len(results) == 0 {
			break
		}
		s.lastSeen = lastID
		s.updates <- results
		if int32(len(results)) < s.options.NumRows {
			break
		}
	}
}
