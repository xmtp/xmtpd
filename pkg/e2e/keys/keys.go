// Package keys manages Anvil private key allocation for E2E tests.
//
// The 10 default Anvil accounts are partitioned by role to avoid collisions:
//
//	Account 0:   admin (on-chain registration, canonical network ops)
//	Account 1:   client (payer envelope signing)
//	Accounts 2-4: gateways (signer + payer keys)
//	Accounts 5-9: nodes (signer keys)
//
// When a pool is exhausted, new ECDSA keys are generated and funded
// from the admin account via the Anvil RPC.
package keys

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"
)

// anvilKeys are the 10 default Anvil accounts, pre-funded with 10000 ETH each.
var anvilKeys = [10]string{
	"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", // account 0
	"0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d", // account 1
	"0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a", // account 2
	"0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6", // account 3
	"0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a", // account 4
	"0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba", // account 5
	"0x92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e", // account 6
	"0x4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356", // account 7
	"0xdbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97", // account 8
	"0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6", // account 9
}

// AdminKey returns the admin private key (account 0).
// Used exclusively for on-chain registration and canonical network operations.
func AdminKey() string {
	return anvilKeys[0]
}

// ClientKey returns the client payer private key (account 1).
// Used for signing payer envelopes in traffic generation.
func ClientKey() string {
	return anvilKeys[1]
}

// Pool indices for each role.
const (
	gatewayPoolStart = 2 // accounts 2-4
	gatewayPoolEnd   = 5
	nodePoolStart    = 5 // accounts 5-9
	nodePoolEnd      = 10
)

// Manager allocates private keys from role-specific pools.
// When a pool is exhausted, it generates new ECDSA keys and funds them
// from the admin account via the Anvil RPC endpoint.
type Manager struct {
	logger *zap.Logger
	rpcURL string

	mu           sync.Mutex
	nextNode     int
	nextGateway  int
	generatedIdx int
}

// NewManager creates a key manager that uses the given RPC URL to fund
// generated keys when pre-funded pools are exhausted.
func NewManager(logger *zap.Logger, rpcURL string) *Manager {
	return &Manager{
		logger: logger,
		rpcURL: rpcURL,
	}
}

// Reset resets the key allocation counters to zero.
// Call this between test runs to ensure each test gets the same
// set of keys, avoiding "public key mismatch" errors when databases
// persist across tests.
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextNode = 0
	m.nextGateway = 0
}

// NextNodeKey returns the next available node signer key.
// Returns pre-funded keys from accounts 5-9 first, then generates
// and funds new keys as needed.
func (m *Manager) NextNodeKey(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	poolIdx := nodePoolStart + m.nextNode
	m.nextNode++

	if poolIdx < nodePoolEnd {
		return anvilKeys[poolIdx], nil
	}

	return m.generateAndFundKey(ctx, "node")
}

// NextGatewayKey returns the next available gateway key.
// Returns pre-funded keys from accounts 2-4 first, then generates
// and funds new keys as needed.
func (m *Manager) NextGatewayKey(ctx context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	poolIdx := gatewayPoolStart + m.nextGateway
	m.nextGateway++

	if poolIdx < gatewayPoolEnd {
		return anvilKeys[poolIdx], nil
	}

	return m.generateAndFundKey(ctx, "gateway")
}

// generateAndFundKey creates a new ECDSA private key and funds it with ETH
// from the admin account. Must be called with m.mu held.
func (m *Manager) generateAndFundKey(ctx context.Context, role string) (string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate %s key: %w", role, err)
	}

	keyHex := "0x" + fmt.Sprintf("%064x", privateKey.D)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	m.logger.Info("generated new key (pool exhausted)",
		zap.String("role", role),
		zap.String("address", address.Hex()),
	)

	if err := m.fundAddress(ctx, address.Hex()); err != nil {
		return "", fmt.Errorf("failed to fund generated %s key: %w", role, err)
	}

	return keyHex, nil
}

// fundAddress sends 1000 ETH from the admin account to the given address
// using Anvil's eth_sendTransaction (no signing needed for pre-funded accounts).
func (m *Manager) fundAddress(ctx context.Context, toAddress string) error {
	adminKey, err := crypto.HexToECDSA(AdminKey()[2:])
	if err != nil {
		return fmt.Errorf("failed to parse admin key: %w", err)
	}
	fromAddress := crypto.PubkeyToAddress(adminKey.PublicKey).Hex()

	// 1000 ETH in wei = 0x3635C9ADC5DEA00000
	payload := map[string]any{
		"jsonrpc": "2.0",
		"method":  "eth_sendTransaction",
		"params": []any{
			map[string]string{
				"from":  fromAddress,
				"to":    toAddress,
				"value": "0x3635C9ADC5DEA00000",
			},
		},
		"id": 1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal rpc request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.rpcURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create rpc request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send funding tx: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("funding tx failed with status %d: %s", resp.StatusCode, respBody)
	}

	var rpcResp struct {
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respBody, &rpcResp); err == nil && rpcResp.Error != nil {
		return fmt.Errorf("funding tx rpc error: %s", rpcResp.Error.Message)
	}

	m.logger.Info("funded generated key",
		zap.String("to", toAddress),
		zap.String("amount", "1000 ETH"),
	)

	return nil
}

// PrivateKeyToHex converts an ECDSA private key to a 0x-prefixed hex string.
func PrivateKeyToHex(key *ecdsa.PrivateKey) string {
	return "0x" + fmt.Sprintf("%064x", key.D)
}
