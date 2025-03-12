package utils

import (
	"context"
	"math/rand"
	"time"
)

func RandomSleep(ctx context.Context, maxDuration time.Duration) {
	if maxDuration <= 0 {
		return // No sleep if duration is invalid
	}
	randDuration := time.Duration(rand.Float64() * float64(maxDuration))

	select {
	case <-time.After(randDuration):
	case <-ctx.Done():
		return
	}
}
