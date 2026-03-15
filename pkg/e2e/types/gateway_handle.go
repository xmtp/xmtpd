package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/e2e/gateway"
	"github.com/xmtp/xmtpd/pkg/utils"
)

// GatewayHandle provides a fluent API for interacting with an xmtpd gateway.
// It wraps the underlying gateway container and composes lifecycle management
// and chaos injection into a single addressable handle.
//
// Gateways are addressed by their creation index (0, 1, 2, ...):
//
//	env.Gateway(0).AddLatency(ctx, 500)
//	env.Gateway(1).Stop(ctx)
type GatewayHandle struct {
	gateway *gateway.Gateway
	env     *Environment
	index   int
}

// newGatewayHandle creates a GatewayHandle wrapping the given gateway with a
// reference to the environment's chaos controller.
func newGatewayHandle(gw *gateway.Gateway, index int, env *Environment) *GatewayHandle {
	return &GatewayHandle{
		gateway: gw,
		index:   index,
		env:     env,
	}
}

// --- Identity ---

// Index returns the gateway's creation index (0, 1, 2, ...).
func (h *GatewayHandle) Index() int {
	return h.index
}

// Alias returns the container's network alias (e.g. "gateway-0").
func (h *GatewayHandle) Alias() string {
	return h.gateway.Alias()
}

// Endpoint returns the internal network address for this gateway (e.g. "http://gateway-0:5050").
func (h *GatewayHandle) Endpoint() string {
	return h.gateway.InternalAddr()
}

// Address returns the Ethereum address derived from this gateway's signer key.
// This is the gateway's payer address on-chain — the address that gets charged
// in the PayerRegistry when the gateway wraps client messages into payer envelopes.
func (h *GatewayHandle) Address() common.Address {
	privKey, err := utils.ParseEcdsaPrivateKey(h.gateway.SignerKey())
	if err != nil {
		panic("failed to parse gateway signer key: " + err.Error())
	}
	return crypto.PubkeyToAddress(privKey.PublicKey)
}

// --- Lifecycle ---

// Stop terminates the gateway's container.
func (h *GatewayHandle) Stop(ctx context.Context) error {
	return h.gateway.Stop(ctx)
}

// --- Chaos injection ---
// All chaos methods delegate to the toxiproxy controller using this gateway's alias.

// AddLatency injects a network latency toxic on this gateway.
// All connections to/from the gateway will experience the specified delay in milliseconds.
func (h *GatewayHandle) AddLatency(ctx context.Context, ms int) error {
	return h.env.Chaos.AddLatency(ctx, h.gateway.Alias(), ms)
}

// AddBandwidthLimit restricts the gateway's network throughput to the specified rate in KB/s.
func (h *GatewayHandle) AddBandwidthLimit(ctx context.Context, kbps int) error {
	return h.env.Chaos.AddBandwidthLimit(ctx, h.gateway.Alias(), kbps)
}

// AddConnectionReset simulates TCP connection resets (RST) on the gateway's connections.
// Connections are reset after the specified timeout in milliseconds.
func (h *GatewayHandle) AddConnectionReset(ctx context.Context, timeoutMs int) error {
	return h.env.Chaos.AddConnectionReset(ctx, h.gateway.Alias(), timeoutMs)
}

// AddTimeout blocks all data and closes connections after the specified timeout.
// If timeoutMs is 0, data is dropped indefinitely without closing the connection,
// effectively simulating a network partition (black hole).
func (h *GatewayHandle) AddTimeout(ctx context.Context, timeoutMs int) error {
	return h.env.Chaos.AddTimeout(ctx, h.gateway.Alias(), timeoutMs)
}

// RemoveAllToxics removes all active toxics from this gateway, restoring normal
// network conditions.
func (h *GatewayHandle) RemoveAllToxics(ctx context.Context) error {
	return h.env.Chaos.RemoveAllToxics(ctx, h.gateway.Alias())
}

// --- Balance queries ---

// GetPayerBalance returns this gateway's balance in the PayerRegistry.
func (h *GatewayHandle) GetPayerBalance(ctx context.Context) (*big.Int, error) {
	return h.env.GetPayerBalance(ctx, h.Address())
}

// GetFeeTokenBalance returns the fee token (xUSD) balance for this gateway's address.
func (h *GatewayHandle) GetFeeTokenBalance(ctx context.Context) (*big.Int, error) {
	return h.env.GetFeeTokenBalance(ctx, h.Address())
}

// GetGasBalance returns the native ETH balance for this gateway's address.
func (h *GatewayHandle) GetGasBalance(ctx context.Context) (*big.Int, error) {
	return h.env.GetGasBalance(ctx, h.Address())
}

// --- Payer operations ---

// Deposit mints fee tokens and deposits them into the PayerRegistry for this
// gateway's payer address. Handles the full flow: mint → wrap → approve → deposit.
func (h *GatewayHandle) Deposit(ctx context.Context, amount *big.Int) error {
	return h.env.FundPayer(ctx, h.Address(), amount)
}

// RequestWithdrawal requests a withdrawal from the PayerRegistry for this gateway.
func (h *GatewayHandle) RequestWithdrawal(ctx context.Context, amount *big.Int) error {
	return h.env.RequestPayerWithdrawal(ctx, h.gateway.SignerKey(), amount)
}

// CancelWithdrawal cancels a pending withdrawal from the PayerRegistry for this gateway.
func (h *GatewayHandle) CancelWithdrawal(ctx context.Context) error {
	return h.env.CancelPayerWithdrawal(ctx, h.gateway.SignerKey())
}

// FinalizeWithdrawal finalizes a pending withdrawal, transferring funds to the recipient.
func (h *GatewayHandle) FinalizeWithdrawal(ctx context.Context, recipient common.Address) error {
	return h.env.FinalizePayerWithdrawal(ctx, h.gateway.SignerKey(), recipient)
}

// DisableProxy completely disables this gateway's proxy, refusing all connections.
func (h *GatewayHandle) DisableProxy(ctx context.Context) error {
	return h.env.Chaos.DisableProxy(ctx, h.gateway.Alias())
}

// EnableProxy re-enables this gateway's proxy after it was disabled.
func (h *GatewayHandle) EnableProxy(ctx context.Context) error {
	return h.env.Chaos.EnableProxy(ctx, h.gateway.Alias())
}
