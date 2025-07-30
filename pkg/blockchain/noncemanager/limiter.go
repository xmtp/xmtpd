package noncemanager

import "sync"

// MaxConcurrentRequests defines the maximum number of concurrent nonce requests allowed.
// The blockchain mempool can usually only hold 64 transactions from the same fromAddress.
const MaxConcurrentRequests = 64

// BestGuessConcurrency is the recommended default concurrency limit for most use cases
const BestGuessConcurrency = 32

// OpenConnectionsLimiter controls the number of concurrent nonce allocation requests
// to prevent overwhelming the blockchain mempool with too many transactions from
// the same address.
type OpenConnectionsLimiter struct {
	Semaphore chan struct{}
	WG        sync.WaitGroup
}

// NewOpenConnectionsLimiter creates a new limiter with the specified maximum concurrent requests.
// If maxConcurrent exceeds MaxConcurrentRequests, it will be capped at that value.
func NewOpenConnectionsLimiter(maxConcurrent int) *OpenConnectionsLimiter {
	if maxConcurrent > MaxConcurrentRequests {
		maxConcurrent = MaxConcurrentRequests
	}
	return &OpenConnectionsLimiter{
		Semaphore: make(chan struct{}, maxConcurrent),
	}
}
