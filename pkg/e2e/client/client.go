// Package client provides traffic generation for E2E tests.
//
// Clients are created via the environment and bound to a specific node:
//
//	env.NewClient(100)                            // create client for node 100
//	env.Client(100).PublishEnvelopes(ctx, 10)     // publish 10 envelopes
//	env.Client(100).GenerateTraffic(ctx, opts)    // background traffic
//	env.Client(100).Stop()                        // stop traffic + cleanup
package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/stress"
	"go.uber.org/zap"
)

// TrafficOptions configures background traffic generation.
type TrafficOptions struct {
	// BatchSize is the number of envelopes to publish per batch.
	BatchSize uint
	// Duration is how long to generate traffic before stopping automatically.
	Duration time.Duration
}

// Client is the interface for E2E traffic generation.
// Implementations are bound to a specific node and can publish envelopes
// synchronously or generate background traffic.
type Client interface {
	// PublishEnvelopes publishes the specified number of group message envelopes
	// to the target node via gRPC. Blocks until all envelopes are published.
	PublishEnvelopes(ctx context.Context, count uint) error

	// GenerateTraffic starts background traffic generation that continuously
	// publishes envelopes in batches until Stop is called or the duration elapses.
	// Returns a TrafficGenerator to monitor and stop the generation.
	GenerateTraffic(ctx context.Context, opts TrafficOptions) *TrafficGenerator

	// Stop stops any active background traffic generation and releases resources.
	Stop()

	// NodeID returns the on-chain nodeID of the node this client publishes to.
	NodeID() uint32

	// PayerKey returns the private key hex string used for signing payer envelopes.
	PayerKey() string

	// Name returns the unique identifier for this client within the environment.
	Name() string

	// Address returns the Ethereum address derived from this client's payer key.
	// This address appears in the payers table and can be used to verify per-payer
	// attribution via GetUnsettledUsage.
	Address() common.Address
}

// Options configures a new client instance.
type Options struct {
	// NodeAddr is the host-accessible address of a node (e.g. http://localhost:XXXXX).
	NodeAddr string
	// PayerKey is the private key used to sign payer envelopes.
	PayerKey string
	// OriginatorID is the on-chain nodeID of the target node.
	OriginatorID uint32
	// Name is the unique identifier for this client within the environment.
	Name string
}

type client struct {
	logger  *zap.Logger
	opts    Options
	address common.Address
	mu      sync.Mutex
	traffic *TrafficGenerator
}

// New creates a new Client bound to the node specified in opts.
func New(logger *zap.Logger, opts Options) Client {
	if opts.PayerKey == "" {
		panic("PayerKey must be provided — use keys.ClientKey()")
	}

	// Parse the payer key and derive the Ethereum address once at construction.
	keyHex := opts.PayerKey
	if len(keyHex) >= 2 && keyHex[:2] == "0x" {
		keyHex = keyHex[2:]
	}
	privateKey, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		panic(fmt.Sprintf("invalid PayerKey: %v", err))
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &client{
		logger:  logger,
		opts:    opts,
		address: address,
	}
}

// NodeID returns the on-chain nodeID of the node this client publishes to.
func (c *client) NodeID() uint32 {
	return c.opts.OriginatorID
}

// PayerKey returns the private key hex string used for signing payer envelopes.
func (c *client) PayerKey() string {
	return c.opts.PayerKey
}

// Name returns the unique identifier for this client within the environment.
func (c *client) Name() string {
	return c.opts.Name
}

// Address returns the Ethereum address derived from this client's payer key.
func (c *client) Address() common.Address {
	return c.address
}

// PublishEnvelopes publishes the specified number of group message envelopes
// to the target node via gRPC. This uses the stress.EnvelopesGenerator
// which creates properly signed payer envelopes with random data.
func (c *client) PublishEnvelopes(ctx context.Context, count uint) error {
	gen, err := stress.NewEnvelopesGenerator(
		c.opts.NodeAddr,
		c.opts.PayerKey,
		c.opts.OriginatorID,
		stress.ProtocolConnectGRPC,
	)
	if err != nil {
		return fmt.Errorf("failed to create envelopes generator: %w", err)
	}
	defer func() {
		_ = gen.Close()
	}()

	_, err = gen.PublishGroupMessageEnvelopes(ctx, count, "256B")
	if err != nil {
		return fmt.Errorf("failed to publish envelopes: %w", err)
	}

	return nil
}

// GenerateTraffic starts a goroutine that continuously publishes envelopes
// in batches until Stop is called, the context is cancelled, or the duration elapses.
// Returns a TrafficGenerator to monitor errors. Only one background traffic
// generation can be active per client; calling GenerateTraffic again stops the previous one.
func (c *client) GenerateTraffic(
	ctx context.Context,
	opts TrafficOptions,
) *TrafficGenerator {
	c.mu.Lock()
	// Stop any existing traffic generation
	if c.traffic != nil {
		c.traffic.Stop()
	}

	genCtx, cancel := context.WithTimeout(ctx, opts.Duration)

	gen := &TrafficGenerator{cancel: cancel}
	gen.wg.Add(1)
	c.traffic = gen
	c.mu.Unlock()

	go func() {
		defer gen.wg.Done()
		defer cancel()

		startTime := time.Now()

		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		c.logger.Info("background traffic generation started",
			zap.Uint("batch_size", opts.BatchSize),
			zap.Duration("duration", opts.Duration),
			zap.Uint32("node_id", c.opts.OriginatorID),
			zap.String("node_addr", c.opts.NodeAddr),
		)

		for {
			select {
			case <-genCtx.Done():
				c.logger.Info("background traffic generation stopped")
				return

			case <-ticker.C:
				remaining := opts.Duration - time.Since(startTime)
				c.logger.Info("background traffic generation ongoing",
					zap.Duration("remaining", remaining.Truncate(time.Second)),
					zap.Uint32("node_id", c.opts.OriginatorID),
					zap.String("node_addr", c.opts.NodeAddr),
				)

			default:
				if err := c.PublishEnvelopes(genCtx, opts.BatchSize); err != nil {
					// context cancellation is expected when Stop is called
					if genCtx.Err() != nil {
						return
					}
					gen.setErr(err)
					c.logger.Error("background traffic generation error",
						zap.Error(err),
					)
					return
				}
			}
		}
	}()

	return gen
}

// Stop stops any active background traffic generation and waits for it to finish.
func (c *client) Stop() {
	c.mu.Lock()
	t := c.traffic
	c.mu.Unlock()
	if t != nil {
		t.Stop()
	}
}

// TrafficGenerator manages a background traffic generation goroutine.
// Call Stop to cancel the generation and wait for it to finish.
// Check Err after Stop to see if the generation encountered an error.
type TrafficGenerator struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.Mutex
	err    error
}

// Err returns the first error from the background generation, if any.
// Safe to call after Stop returns.
func (g *TrafficGenerator) Err() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.err
}

// Stop cancels the background generation and waits for it to finish.
func (g *TrafficGenerator) Stop() {
	g.cancel()
	g.wg.Wait()
}

func (g *TrafficGenerator) setErr(err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.err == nil {
		g.err = err
	}
}
