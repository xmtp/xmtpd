package types

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/e2e/gateway"
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

// Address returns the host-accessible address for this gateway.
func (h *GatewayHandle) Address() string {
	return h.gateway.InternalAddr()
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
