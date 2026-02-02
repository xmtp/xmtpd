package db

import (
	"context"
	"errors"
	"time"

	"github.com/xmtp/xmtpd/pkg/tracing"
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
		return nil, errors.New("already started")
	}
	if s.options.NumRows <= 0 || s.logger == nil {
		return nil, errors.New("required params not provided")
	}
	updates := make(chan []ValueType)
	s.updates = updates

	go func() {
		s.poll("startup")

		timer := time.NewTimer(s.options.Interval)
		for {
			timer.Reset(s.options.Interval)
			select {
			case <-s.ctx.Done():
				s.logger.Debug("context done; stopping")
				close(s.updates)
				return
			case <-s.options.Notifier:
				s.poll("notification")
			case <-timer.C:
				s.poll("timer_fallback")
			}
		}
	}()

	return updates, nil
}

func (s *DBSubscription[ValueType, CursorType]) poll(trigger string) {
	// Create APM span for polling - this helps identify notification vs timer_fallback
	span, ctx := tracing.StartSpanFromContext(s.ctx, tracing.SpanDBSubscriptionPoll)
	defer span.Finish()

	// Tag with trigger type - this is KEY for debugging the read-replica issue!
	// If you see lots of "timer_fallback" with num_results > 0, the notification
	// poll is missing data (likely due to read-replica lag)
	tracing.SpanTag(span, "trigger", trigger)

	// Repeatedly query page by page until no more results
	totalResults := 0
	for {
		results, lastID, err := s.query(ctx, s.lastSeen, s.options.NumRows)
		if ctx.Err() != nil {
			return
		}

		if err != nil {
			span.Finish(tracing.WithError(err))
			// Log is extremely noisy during test teardown
			s.logger.Error(
				"error querying for database subscription",
				zap.Error(err),
				utils.NumRowsField(s.options.NumRows),
			)

			// Did not update lastSeen; will retry on next poll
			return
		}

		if len(results) == 0 {
			tracing.SpanTag(span, "num_results", totalResults)
			if totalResults == 0 && trigger == "notification" {
				// This indicates the notification poll missed data - likely read-replica lag!
				tracing.SpanTag(span, "notification_miss", true)
			}
			return
		}

		totalResults += len(results)
		s.lastSeen = lastID
		s.updates <- results

		// If we have less results than allowed, it means there's currently no more items to retrieve.
		// Else repeat query and return more batches.
		if int32(len(results)) < s.options.NumRows {
			tracing.SpanTag(span, "num_results", totalResults)
			return
		}
	}
}
