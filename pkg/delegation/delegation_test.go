package delegation

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockChainVerifier is a mock implementation of ChainVerifier for testing.
type MockChainVerifier struct {
	authorizations map[delegationKey]bool
	delegations    map[delegationKey]*DelegationInfo
	callCount      int
}

func NewMockChainVerifier() *MockChainVerifier {
	return &MockChainVerifier{
		authorizations: make(map[delegationKey]bool),
		delegations:    make(map[delegationKey]*DelegationInfo),
	}
}

func (m *MockChainVerifier) SetAuthorized(payer, delegate common.Address, authorized bool) {
	key := delegationKey{payer: payer, delegate: delegate}
	m.authorizations[key] = authorized
}

func (m *MockChainVerifier) SetDelegation(payer, delegate common.Address, info *DelegationInfo) {
	key := delegationKey{payer: payer, delegate: delegate}
	m.delegations[key] = info
}

func (m *MockChainVerifier) IsAuthorized(
	ctx context.Context,
	payer, delegate common.Address,
) (bool, error) {
	m.callCount++
	key := delegationKey{payer: payer, delegate: delegate}
	return m.authorizations[key], nil
}

func (m *MockChainVerifier) GetDelegation(
	ctx context.Context,
	payer, delegate common.Address,
) (*DelegationInfo, error) {
	key := delegationKey{payer: payer, delegate: delegate}
	if info, exists := m.delegations[key]; exists {
		return info, nil
	}
	return &DelegationInfo{}, nil
}

func TestCachingVerifier_CacheHit(t *testing.T) {
	mock := NewMockChainVerifier()
	payer := common.HexToAddress("0x1234567890123456789012345678901234567890")
	delegate := common.HexToAddress("0x0987654321098765432109876543210987654321")

	mock.SetAuthorized(payer, delegate, true)
	mock.SetDelegation(payer, delegate, &DelegationInfo{
		IsActive:  true,
		Expiry:    0,
		CreatedAt: uint64(time.Now().Unix()),
	})

	verifier := NewCachingVerifier(mock, time.Minute)

	ctx := context.Background()

	// First call - should hit chain
	authorized, err := verifier.IsAuthorized(ctx, payer, delegate)
	require.NoError(t, err)
	assert.True(t, authorized)
	assert.Equal(t, 1, mock.callCount)

	// Second call - should hit cache
	authorized, err = verifier.IsAuthorized(ctx, payer, delegate)
	require.NoError(t, err)
	assert.True(t, authorized)
	assert.Equal(t, 1, mock.callCount) // No additional chain calls

	// Third call - still cached
	authorized, err = verifier.IsAuthorized(ctx, payer, delegate)
	require.NoError(t, err)
	assert.True(t, authorized)
	assert.Equal(t, 1, mock.callCount)
}

func TestCachingVerifier_CacheExpiry(t *testing.T) {
	mock := NewMockChainVerifier()
	payer := common.HexToAddress("0x1234567890123456789012345678901234567890")
	delegate := common.HexToAddress("0x0987654321098765432109876543210987654321")

	mock.SetAuthorized(payer, delegate, true)
	mock.SetDelegation(payer, delegate, &DelegationInfo{
		IsActive:  true,
		Expiry:    0,
		CreatedAt: uint64(time.Now().Unix()),
	})

	// Very short TTL for testing
	verifier := NewCachingVerifier(mock, 10*time.Millisecond)

	ctx := context.Background()

	// First call
	_, err := verifier.IsAuthorized(ctx, payer, delegate)
	require.NoError(t, err)
	assert.Equal(t, 1, mock.callCount)

	// Wait for cache to expire
	time.Sleep(20 * time.Millisecond)

	// Second call - should hit chain again
	_, err = verifier.IsAuthorized(ctx, payer, delegate)
	require.NoError(t, err)
	assert.Equal(t, 2, mock.callCount)
}

func TestCachingVerifier_DelegationExpiry(t *testing.T) {
	mock := NewMockChainVerifier()
	payer := common.HexToAddress("0x1234567890123456789012345678901234567890")
	delegate := common.HexToAddress("0x0987654321098765432109876543210987654321")

	// Set delegation that expires soon
	expiryTime := uint64(time.Now().Add(100 * time.Millisecond).Unix())
	mock.SetAuthorized(payer, delegate, true)
	mock.SetDelegation(payer, delegate, &DelegationInfo{
		IsActive:  true,
		Expiry:    expiryTime,
		CreatedAt: uint64(time.Now().Unix()),
	})

	// Long cache TTL to ensure we're testing delegation expiry, not cache expiry
	verifier := NewCachingVerifier(mock, time.Hour)

	ctx := context.Background()

	// First call - should be authorized
	authorized, err := verifier.IsAuthorized(ctx, payer, delegate)
	require.NoError(t, err)
	assert.True(t, authorized)

	// Wait for delegation to expire
	time.Sleep(150 * time.Millisecond)

	// Second call - should return false due to delegation expiry
	authorized, err = verifier.IsAuthorized(ctx, payer, delegate)
	require.NoError(t, err)
	assert.False(t, authorized)
}

func TestCachingVerifier_NotAuthorized(t *testing.T) {
	mock := NewMockChainVerifier()
	payer := common.HexToAddress("0x1234567890123456789012345678901234567890")
	delegate := common.HexToAddress("0x0987654321098765432109876543210987654321")

	mock.SetAuthorized(payer, delegate, false)

	verifier := NewCachingVerifier(mock, time.Minute)

	ctx := context.Background()

	authorized, err := verifier.IsAuthorized(ctx, payer, delegate)
	require.NoError(t, err)
	assert.False(t, authorized)
}

func TestCachingVerifier_InvalidateCache(t *testing.T) {
	mock := NewMockChainVerifier()
	payer := common.HexToAddress("0x1234567890123456789012345678901234567890")
	delegate := common.HexToAddress("0x0987654321098765432109876543210987654321")

	mock.SetAuthorized(payer, delegate, true)
	mock.SetDelegation(payer, delegate, &DelegationInfo{IsActive: true})

	verifier := NewCachingVerifier(mock, time.Hour)

	ctx := context.Background()

	// First call
	_, err := verifier.IsAuthorized(ctx, payer, delegate)
	require.NoError(t, err)
	assert.Equal(t, 1, mock.callCount)

	// Invalidate cache
	verifier.InvalidateCache(payer, delegate)

	// Next call should hit chain
	_, err = verifier.IsAuthorized(ctx, payer, delegate)
	require.NoError(t, err)
	assert.Equal(t, 2, mock.callCount)
}

func TestCachingVerifier_ClearCache(t *testing.T) {
	mock := NewMockChainVerifier()
	payer := common.HexToAddress("0x1234567890123456789012345678901234567890")
	delegate := common.HexToAddress("0x0987654321098765432109876543210987654321")

	mock.SetAuthorized(payer, delegate, true)
	mock.SetDelegation(payer, delegate, &DelegationInfo{IsActive: true})

	verifier := NewCachingVerifier(mock, time.Hour)

	ctx := context.Background()

	_, err := verifier.IsAuthorized(ctx, payer, delegate)
	require.NoError(t, err)
	assert.Equal(t, 1, verifier.CacheSize())

	verifier.ClearCache()
	assert.Equal(t, 0, verifier.CacheSize())
}
