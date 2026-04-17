// network-watcher subscribes to every registered xmtpd node's metadata-API
// sync-cursor stream and exposes Prometheus metrics describing global sync
// state. It is designed to be scraped alongside chain-watcher.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/net/http2"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/config/environments"
	"github.com/xmtp/xmtpd/pkg/networkwatcher"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
)

var Version string

// NetworkWatcherOptions holds go-flags options for the network-watcher binary.
type NetworkWatcherOptions struct {
	Log config.LogOptions `group:"Log Options" namespace:"log"`

	Version bool `short:"v" long:"version"`

	Contracts string `long:"contracts-env" env:"XMTPD_CONTRACTS_ENVIRONMENT"`

	RPCURL string `long:"rpc-url" env:"SETTLEMENT_CHAIN_RPC_URL" required:"true"`

	NodeRegistryAddress string `long:"node-registry-address" env:"NODE_REGISTRY_ADDRESS"`

	NodeRegistryRefreshInterval time.Duration `long:"node-registry-refresh-interval" env:"NODE_REGISTRY_REFRESH_INTERVAL" default:"60s"`

	ReconnectMinBackoff time.Duration `long:"reconnect-min-backoff" env:"RECONNECT_MIN_BACKOFF" default:"1s"`
	ReconnectMaxBackoff time.Duration `long:"reconnect-max-backoff" env:"RECONNECT_MAX_BACKOFF" default:"30s"`

	MetricsPort string `long:"metrics-port" env:"METRICS_PORT" default:"8009"`
}

// envConfig is the subset of the contracts-env JSON we need.
type envConfig struct {
	NodeRegistry string `json:"nodeRegistry"`
}

var options NetworkWatcherOptions

func main() {
	_, err := flags.Parse(&options)
	if err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && flagsErr.Type == flags.ErrHelp {
			return
		}
		fatal("could not parse options: %s", err)
	}

	if Version == "" {
		Version = os.Getenv("VERSION")
	}
	if options.Version {
		fmt.Printf("version: %s\n", Version)
		return
	}

	logger, _, err := utils.BuildLogger(options.Log)
	if err != nil {
		fatal("could not build logger: %s", err)
	}
	defer func() { _ = logger.Sync() }()

	if Version != "" {
		logger.Info("version: " + Version)
	}

	resolveContractAddresses(logger)

	reg := prometheus.NewRegistry()
	networkwatcher.RegisterMetrics(reg)

	go serveMetrics(logger, reg, options.MetricsPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	settlementClient, err := blockchain.NewRPCClient(ctx, options.RPCURL)
	if err != nil {
		logger.Fatal("initializing blockchain client", zap.Error(err))
	}

	contractsOpts := config.ContractsOptions{
		SettlementChain: config.SettlementChainOptions{
			NodeRegistryAddress:         options.NodeRegistryAddress,
			NodeRegistryRefreshInterval: options.NodeRegistryRefreshInterval,
			RPCURL:                      options.RPCURL,
		},
	}

	chainRegistry, err := registry.NewSmartContractRegistry(
		ctx,
		settlementClient,
		logger,
		&contractsOpts,
	)
	if err != nil {
		logger.Fatal("initializing smart contract registry", zap.Error(err))
	}
	if err := chainRegistry.Start(); err != nil {
		logger.Fatal("starting smart contract registry", zap.Error(err))
	}
	defer chainRegistry.Stop()

	watcher, err := networkwatcher.NewWatcher(networkwatcher.WatcherConfig{
		Registry:   chainRegistry,
		Logger:     logger.Named("xmtpd.network-watcher"),
		HTTPClient: newStreamingHTTPClient(logger),
		MinBackoff: options.ReconnectMinBackoff,
		MaxBackoff: options.ReconnectMaxBackoff,
	})
	if err != nil {
		logger.Fatal("creating network watcher", zap.Error(err))
	}
	if err := watcher.Start(ctx); err != nil {
		logger.Fatal("starting network watcher", zap.Error(err))
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	sig := <-sigCh
	logger.Info("received signal, shutting down", zap.String("signal", sig.String()))

	cancel()
	watcher.Stop()
}

func resolveContractAddresses(logger *zap.Logger) {
	if options.Contracts == "" {
		if options.NodeRegistryAddress == "" {
			fatal("--node-registry-address is required when --contracts-env is not set")
		}
		return
	}

	var env environments.SmartContractEnvironment
	if err := env.UnmarshalFlag(options.Contracts); err != nil {
		logger.Fatal("invalid contracts environment", zap.Error(err))
	}
	data, err := environments.GetEnvironmentConfig(env)
	if err != nil {
		logger.Fatal("failed to load environment config", zap.Error(err))
	}
	var ec envConfig
	if err := json.Unmarshal(data, &ec); err != nil {
		logger.Fatal("failed to parse environment config", zap.Error(err))
	}
	if options.NodeRegistryAddress == "" {
		options.NodeRegistryAddress = ec.NodeRegistry
	}
	logger.Info(
		"loaded contract addresses from environment",
		zap.String("environment", options.Contracts),
		zap.String("node_registry", options.NodeRegistryAddress),
	)
}

// newStreamingHTTPClient returns an http.Client tuned for long-lived
// server-streaming RPCs. On HTTPS connections it sends HTTP/2 PING frames
// on otherwise-idle streams so that intermediaries (LBs, reverse proxies)
// with idle timeouts don't silently terminate them. Plain http:// falls
// back to HTTP/1.1 — sufficient for dev/local setups where idle timeouts
// aren't a concern.
func newStreamingHTTPClient(logger *zap.Logger) *http.Client {
	t := &http.Transport{ForceAttemptHTTP2: true}
	h2t, err := http2.ConfigureTransports(t)
	if err != nil {
		logger.Warn("could not configure HTTP/2 keepalive on transport", zap.Error(err))
	} else {
		h2t.ReadIdleTimeout = 15 * time.Second
		h2t.PingTimeout = 10 * time.Second
	}
	return &http.Client{Transport: t}
}

func serveMetrics(logger *zap.Logger, reg *prometheus.Registry, port string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{EnableOpenMetrics: true}))
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := net.JoinHostPort("0.0.0.0", port)
	logger.Info("serving metrics", zap.String("address", addr))
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("metrics server error", zap.Error(err))
	}
}

func fatal(msg string, args ...any) {
	log.Fatalf(msg, args...)
}
