package context_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/context"
	tests "github.com/xmtp/xmtpd/pkg/testing"
)

func Test_GoClose(t *testing.T) {
	ctx := context.New(context.Background(), tests.NewLogger(t))
	ctx.Go(func(ctx context.Context) {
		<-ctx.Done()
	})
	ctx.Close()
}

func Test_GoCloseWithCancel(t *testing.T) {
	var finished sync.Map
	ctx := context.New(context.Background(), tests.NewLogger(t))
	ctx.Go(func(ctx context.Context) {
		<-ctx.Done()
		finished.Store("one", true)
	})
	ctx2 := context.WithCancel(ctx)
	ctx2.Go(func(ctx context.Context) {
		<-ctx.Done()
		finished.Store("two", true)
	})
	// Close child
	ctx2.Close()
	_, ok := finished.Load("one")
	require.False(t, ok)
	_, ok = finished.Load("two")
	require.True(t, ok)
	// Close parent
	ctx.Close()
	_, ok = finished.Load("one")
	require.True(t, ok)
}

func Test_GoCloseWithTimeout(t *testing.T) {
	var finished sync.Map
	ctx := context.New(context.Background(), tests.NewLogger(t))
	ctx.Go(func(ctx context.Context) {
		<-ctx.Done()
		finished.Store("one", true)
	})
	ctx2 := context.WithTimeout(ctx, 10*time.Millisecond)
	ctx2.Go(func(ctx context.Context) {
		<-ctx.Done()
		finished.Store("two", true)
	})
	// Wait for child to expire
	time.Sleep(20 * time.Millisecond)
	_, ok := finished.Load("two")
	require.True(t, ok)
	// Close the child
	ctx2.Close()
	_, ok = finished.Load("one")
	require.False(t, ok)
	// Close parent
	ctx.Close()
	_, ok = finished.Load("one")
	require.True(t, ok)
}

func Test_GoCloseWithDeadline(t *testing.T) {
	var finished sync.Map
	ctx := context.New(context.Background(), tests.NewLogger(t))
	ctx.Go(func(ctx context.Context) {
		<-ctx.Done()
		finished.Store("one", true)
	})
	ctx2 := context.WithDeadline(ctx, time.Now().Add(10*time.Millisecond))
	ctx2.Go(func(ctx context.Context) {
		<-ctx.Done()
		finished.Store("two", true)
	})
	// Wait for child to expire
	time.Sleep(20 * time.Millisecond)
	_, ok := finished.Load("two")
	require.True(t, ok)
	// Close the child
	ctx2.Close()
	_, ok = finished.Load("one")
	require.False(t, ok)
	// Close parent
	ctx.Close()
	_, ok = finished.Load("one")
	require.True(t, ok)
}

func Test_GoCloseWithWithValue(t *testing.T) {
	var finished sync.Map
	ctx := context.New(context.Background(), tests.NewLogger(t))
	ctx.Go(func(ctx context.Context) {
		<-ctx.Done()
		finished.Store("one", true)
	})
	key := "key"
	ctx2 := context.WithValue(ctx, key, true)
	require.Equal(t, true, ctx2.Value(key))

	// Closing child closes the parent
	ctx2.Close()
	_, ok := finished.Load("one")
	require.True(t, ok)
}
