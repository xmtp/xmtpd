package types

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/e2e/client"
	"github.com/xmtp/xmtpd/pkg/utils"
)

// ClientHandle wraps a traffic generation client and adds payer convenience
// methods for balance queries and deposit/withdrawal operations.
//
// NOTE: In the current e2e framework, clients publish directly to nodes and
// sign their own payer envelopes — making them act as both client AND payer.
// In production, the gateway is the actual payer; the client just sends
// messages through a gateway. The payer methods here exist for testing
// direct-publish scenarios only.
type ClientHandle struct {
	client client.Client
	env    *Environment
}

// newClientHandle creates a ClientHandle wrapping the given client with a
// reference to the environment for chain operations.
func newClientHandle(c client.Client, env *Environment) *ClientHandle {
	return &ClientHandle{
		client: c,
		env:    env,
	}
}

// --- Traffic (delegated to client.Client) ---

// PublishEnvelopes publishes the specified number of envelopes to the target node.
func (h *ClientHandle) PublishEnvelopes(ctx context.Context, count uint) error {
	return h.client.PublishEnvelopes(ctx, count)
}

// GenerateTraffic starts background traffic generation.
func (h *ClientHandle) GenerateTraffic(
	ctx context.Context,
	opts client.TrafficOptions,
) *client.TrafficGenerator {
	return h.client.GenerateTraffic(ctx, opts)
}

// Stop stops any active background traffic generation.
func (h *ClientHandle) Stop() {
	h.client.Stop()
}

// NodeID returns the on-chain nodeID of the node this client publishes to.
func (h *ClientHandle) NodeID() uint32 {
	return h.client.NodeID()
}

// --- Identity ---

// Address returns the Ethereum address derived from this client's payer key.
// In the current e2e framework this is the address that gets charged in the
// PayerRegistry when publishing directly to a node. In production, the
// gateway's address would be charged instead.
func (h *ClientHandle) Address() common.Address {
	privKey, err := utils.ParseEcdsaPrivateKey(h.client.PayerKey())
	if err != nil {
		panic("failed to parse client payer key: " + err.Error())
	}
	return crypto.PubkeyToAddress(privKey.PublicKey)
}

// PayerAddress returns the hex-encoded Ethereum address for this client's payer key.
// Convenience wrapper around Address().Hex() for use with string-based assertions.
func (h *ClientHandle) PayerAddress() string {
	return h.Address().Hex()
}

// --- Balance queries ---

// GetPayerBalance returns this client's balance in the PayerRegistry.
func (h *ClientHandle) GetPayerBalance(ctx context.Context) (*big.Int, error) {
	return h.env.GetPayerBalance(ctx, h.Address())
}

// GetFeeTokenBalance returns the fee token (xUSD) balance for this client's address.
func (h *ClientHandle) GetFeeTokenBalance(ctx context.Context) (*big.Int, error) {
	return h.env.GetFeeTokenBalance(ctx, h.Address())
}

// GetGasBalance returns the native ETH balance for this client's address.
func (h *ClientHandle) GetGasBalance(ctx context.Context) (*big.Int, error) {
	return h.env.GetGasBalance(ctx, h.Address())
}

// --- Payer operations ---
// NOTE: These methods exist because in the current e2e framework, clients
// publish directly to nodes and act as their own payer. In production, the
// gateway handles payer operations. See the TODO in pkg/e2e/README.md.

// Deposit mints fee tokens and deposits them into the PayerRegistry for this
// client's payer address. Handles the full flow: mint → wrap → approve → deposit.
func (h *ClientHandle) Deposit(ctx context.Context, amount *big.Int) error {
	return h.env.FundPayer(ctx, h.Address(), amount)
}

// RequestWithdrawal requests a withdrawal from the PayerRegistry for this client.
func (h *ClientHandle) RequestWithdrawal(ctx context.Context, amount *big.Int) error {
	return h.env.RequestPayerWithdrawal(ctx, h.client.PayerKey(), amount)
}

// CancelWithdrawal cancels a pending withdrawal from the PayerRegistry for this client.
func (h *ClientHandle) CancelWithdrawal(ctx context.Context) error {
	return h.env.CancelPayerWithdrawal(ctx, h.client.PayerKey())
}

// FinalizeWithdrawal finalizes a pending withdrawal, transferring funds to the recipient.
func (h *ClientHandle) FinalizeWithdrawal(
	ctx context.Context,
	recipient common.Address,
) error {
	return h.env.FinalizePayerWithdrawal(ctx, h.client.PayerKey(), recipient)
}
