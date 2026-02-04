// Package delegation provides delegation verification with caching support.
// It enables user-funded messages by verifying that a gateway (delegate) is authorized
// to sign payer envelopes on behalf of a user (payer).
package delegation

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// DelegationInfo represents the delegation status from the on-chain contract.
type DelegationInfo struct {
	IsActive  bool
	Expiry    uint64
	CreatedAt uint64
}

// ChainVerifier is an interface for verifying delegations on-chain.
type ChainVerifier interface {
	// IsAuthorized checks if the delegate is authorized to sign on behalf of the payer.
	IsAuthorized(ctx context.Context, payer, delegate common.Address) (bool, error)

	// GetDelegation returns the full delegation info for a payer/delegate pair.
	GetDelegation(ctx context.Context, payer, delegate common.Address) (*DelegationInfo, error)
}

// CacheEntry represents a cached delegation status.
type CacheEntry struct {
	IsAuthorized bool
	Expiry       uint64    // Delegation expiry from contract (0 = no expiry)
	CachedAt     time.Time // When this entry was cached
}

// CachingVerifier wraps a ChainVerifier with a time-based cache.
// Caching Strategy:
//   - Cache delegation status for a configurable TTL (default 5 minutes)
//   - On cache hit: return cached value if not expired
//   - On cache miss or expiry: fetch from chain and cache
//   - Delegation revocations may have latency up to the TTL
//   - For delegations with on-chain expiry, we also check if expiry has passed
type CachingVerifier struct {
	chainVerifier ChainVerifier
	cache         map[delegationKey]*CacheEntry
	mu            sync.RWMutex
	cacheTTL      time.Duration
}

type delegationKey struct {
	payer    common.Address
	delegate common.Address
}

// DefaultCacheTTL is the default time-to-live for cached delegation entries.
// This represents the maximum latency for revocation propagation.
const DefaultCacheTTL = 5 * time.Minute

// NewCachingVerifier creates a new CachingVerifier with the given chain verifier.
func NewCachingVerifier(chainVerifier ChainVerifier, cacheTTL time.Duration) *CachingVerifier {
	if cacheTTL == 0 {
		cacheTTL = DefaultCacheTTL
	}
	return &CachingVerifier{
		chainVerifier: chainVerifier,
		cache:         make(map[delegationKey]*CacheEntry),
		cacheTTL:      cacheTTL,
	}
}

// IsAuthorized checks if a delegate is authorized to sign on behalf of a payer.
// Uses caching to reduce on-chain lookups.
func (v *CachingVerifier) IsAuthorized(
	ctx context.Context,
	payer, delegate common.Address,
) (bool, error) {
	key := delegationKey{payer: payer, delegate: delegate}

	// Try cache first
	v.mu.RLock()
	entry, exists := v.cache[key]
	v.mu.RUnlock()

	now := time.Now()

	if exists {
		// Check if cache entry is still valid
		if now.Sub(entry.CachedAt) < v.cacheTTL {
			// Also check if the on-chain delegation has expired
			if entry.Expiry != 0 && uint64(now.Unix()) >= entry.Expiry {
				return false, nil
			}
			return entry.IsAuthorized, nil
		}
	}

	// Cache miss or expired - fetch from chain
	authorized, err := v.chainVerifier.IsAuthorized(ctx, payer, delegate)
	if err != nil {
		return false, err
	}

	// Get delegation info for expiry
	delegationInfo, delegationErr := v.chainVerifier.GetDelegation(ctx, payer, delegate)
	if delegationErr != nil {
		// If we can't get delegation info, still cache the authorization status
		v.mu.Lock()
		v.cache[key] = &CacheEntry{
			IsAuthorized: authorized,
			Expiry:       0,
			CachedAt:     now,
		}
		v.mu.Unlock()
		return authorized, nil //nolint:nilerr // delegationErr is separate from err; err was handled above
	}

	// Cache the result
	v.mu.Lock()
	v.cache[key] = &CacheEntry{
		IsAuthorized: authorized,
		Expiry:       delegationInfo.Expiry,
		CachedAt:     now,
	}
	v.mu.Unlock()

	return authorized, nil
}

// InvalidateCache removes a specific delegation entry from the cache.
// Useful when we know a revocation has occurred.
func (v *CachingVerifier) InvalidateCache(payer, delegate common.Address) {
	key := delegationKey{payer: payer, delegate: delegate}
	v.mu.Lock()
	delete(v.cache, key)
	v.mu.Unlock()
}

// ClearCache removes all entries from the cache.
func (v *CachingVerifier) ClearCache() {
	v.mu.Lock()
	v.cache = make(map[delegationKey]*CacheEntry)
	v.mu.Unlock()
}

// CacheSize returns the current number of entries in the cache.
func (v *CachingVerifier) CacheSize() int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return len(v.cache)
}
