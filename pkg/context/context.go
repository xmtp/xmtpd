package context

import (
	"context"
	"time"

	"github.com/xmtp/xmtpd/pkg/zap"
)

var _ Context = &runtimeContext{}

var (
	// re-export commonly used bits
	Background       = context.Background
	Canceled         = context.Canceled
	DeadlineExceeded = context.DeadlineExceeded
)

// Context is a cancellable context that provides synchronization
// primitives with goroutines it controls and common runtime facilities, e.g. Logger.
//
// Child Context's goroutines are controlled by both the child and the parent.
// Parent's goroutines are out of scope of the child context.
type Context interface {
	context.Context
	// Go runs f in a goroutine that can be cancelled and waited for its exit.
	// f MUST respond to ctx.Done()
	Go(f func(ctx Context))
	// Close cancels All goroutines and waits for them to exit.
	Close()
	// Logger returns the logger carried by this context.
	Logger() *zap.Logger

	// This is a private helper that we need to implement the standard Context API.
	// It will also prevent others from implementing this interface.
	clone() *runtimeContext
}

type runtimeContext struct {
	context.Context
	cancel context.CancelFunc
	wg     *waitGroup
	logger *zap.Logger
}

// New converts a standard Context into our Context,
// it also attaches the logger to it.
func New(ctx context.Context, logger *zap.Logger) Context {
	ctx, cancel := context.WithCancel(ctx)
	return &runtimeContext{
		Context: ctx,
		cancel:  cancel,
		logger:  logger,
		wg:      newWaitGroup(nil),
	}
}

func WithCancel(ctx Context) Context {
	rc := ctx.clone()
	rc.wg = newWaitGroup(rc.wg)
	rc.Context, rc.cancel = context.WithCancel(ctx)
	return rc
}

func WithTimeout(ctx Context, t time.Duration) Context {
	rc := ctx.clone()
	rc.wg = newWaitGroup(rc.wg)
	rc.Context, rc.cancel = context.WithTimeout(ctx, t)
	return rc
}

func WithDeadline(ctx Context, t time.Time) Context {
	rc := ctx.clone()
	rc.wg = newWaitGroup(rc.wg)
	rc.Context, rc.cancel = context.WithDeadline(ctx, t)
	return rc
}

func WithValue(ctx Context, key, val any) Context {
	rc := ctx.clone()
	rc.Context = context.WithValue(ctx, key, val)
	return rc
}

func WithLogger(ctx Context, logger *zap.Logger) Context {
	rc := ctx.clone()
	rc.logger = logger
	return rc
}

func (ctx *runtimeContext) Go(f func(context Context)) {
	ctx.wg.Add(1)
	go func() {
		defer ctx.wg.Done()
		f(ctx)
	}()
}

func (ctx *runtimeContext) Close() {
	ctx.cancel()
	ctx.wg.Wait()
}

func (ctx *runtimeContext) Logger() *zap.Logger {
	return ctx.logger
}

func (ctx *runtimeContext) clone() *runtimeContext {
	return &runtimeContext{
		Context: ctx.Context,
		cancel:  ctx.cancel,
		logger:  ctx.logger,
		wg:      ctx.wg,
	}
}
