// Package noncemanager provides interfaces and implementations for managing transaction nonces
// in a blockchain environment. It supports both SQL and Redis-backed storage implementations
// with concurrent request limiting and atomic nonce allocation.
package noncemanager

import (
	"context"
	"math/big"
)

// NonceContext represents a reserved nonce with associated lifecycle management functions
type NonceContext struct {
	// Nonce is the allocated nonce value
	Nonce big.Int
	// Cancel releases the nonce reservation, making it available for reuse
	Cancel func()
	// Consume commits the nonce, marking it as permanently used
	Consume func() error
}

// NonceManager defines the interface for managing transaction nonces.
// Implementations must provide thread-safe operations for nonce allocation,
// replenishment, and fast-forwarding.
type NonceManager interface {
	// GetNonce atomically reserves and returns the next available nonce.
	// The returned NonceContext allows the caller to either consume or cancel the nonce.
	// This operation may block if the maximum concurrent request limit is reached.
	GetNonce(ctx context.Context) (*NonceContext, error)

	// FastForwardNonce sets the nonce sequence to start from the given value,
	// abandoning all nonces below it. This is typically used when recovering
	// from blockchain state inconsistencies.
	FastForwardNonce(ctx context.Context, nonce big.Int) error

	// Replenish ensures a sufficient number of nonces are available starting
	// from the given nonce value. This is used to pre-populate the nonce pool.
	Replenish(ctx context.Context, nonce big.Int) error
}
