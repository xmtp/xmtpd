package indexerpoc

import (
	"context"
	"fmt"
	"sync"
)

// MemoryStorage provides an in-memory implementation of the Storage interface
// Useful for testing or simple applications
type MemoryStorage struct {
	mu     sync.RWMutex
	states map[string]*taskState // Key is contractName:network
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		states: make(map[string]*taskState),
	}
}

// makeKey creates a composite key for storing state
func makeKey(contractName, network string) string {
	return contractName + ":" + network
}

// SaveState saves the indexing state for a contract
func (s *MemoryStorage) SaveState(ctx context.Context, state *taskState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Make a copy to prevent external modifications
	stateCopy := *state
	key := makeKey(state.ContractName, state.NetworkName)
	s.states[key] = &stateCopy

	return nil
}

// GetState retrieves the indexing state for a contract
func (s *MemoryStorage) GetState(
	ctx context.Context,
	contractName string,
	network string,
) (*taskState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := makeKey(contractName, network)
	state, ok := s.states[key]
	if !ok {
		return nil, fmt.Errorf(
			"no state found for contract: %s on network: %s",
			contractName,
			network,
		)
	}

	// Return a copy to prevent external modifications
	stateCopy := *state
	return &stateCopy, nil
}

// DeleteFromBlock deletes all data for a contract from the specified block number onwards
func (s *MemoryStorage) DeleteFromBlock(
	ctx context.Context,
	contractName string,
	network string,
	blockNumber uint64,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := makeKey(contractName, network)
	state, ok := s.states[key]
	if !ok {
		return nil // Nothing to delete
	}

	// If the stored state is for a higher block, revert it
	if state.BlockNumber >= blockNumber {
		state.BlockNumber = blockNumber - 1
		// BlockHash would be updated separately
	}

	return nil
}

// MemoryTransaction implements a simple in-memory transaction
type MemoryTransaction struct {
	storage   *MemoryStorage
	committed bool
}

// Begin starts a new transaction
func (s *MemoryStorage) Begin(ctx context.Context) (Transaction, error) {
	return &MemoryTransaction{
		storage:   s,
		committed: false,
	}, nil
}

// Commit commits the transaction
func (t *MemoryTransaction) Commit() error {
	t.committed = true
	return nil
}

// Rollback rolls back the transaction
func (t *MemoryTransaction) Rollback() error {
	// Nothing to do for this simple implementation
	return nil
}
