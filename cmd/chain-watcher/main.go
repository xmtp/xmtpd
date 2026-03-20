// chain-watcher is a standalone service that monitors payer report events on
// the XMTP settlement chain and emits Prometheus metrics.
//
// It subscribes to PayerReportSubmitted, PayerReportSubsetSettled, and
// UsageSettled events and derives health signals: submission cadence,
// settlement latency, node participation, envelope gaps, and fee flow.
//
// Exposes a /metrics endpoint for scraping by AMP (Amazon Managed Prometheus).
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
	"github.com/xmtp/xmtpd/pkg/chainwatcher"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/config/environments"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var Version string

// ChainWatcherOptions holds go-flags options for the chain-watcher binary.
type ChainWatcherOptions struct {
	Log     config.LogOptions `group:"Log Options" namespace:"log"`
	Version bool              `short:"v" long:"version" description:"Output binary version and exit"`

	Contracts string `long:"contracts-environment" env:"XMTPD_CONTRACTS_ENVIRONMENT" description:"Named contract environment (e.g. testnet-dev)"`

	SettlementChainRPCURL string `long:"settlement-chain-rpc-url" env:"SETTLEMENT_CHAIN_RPC_URL" description:"Settlement chain HTTP RPC URL" required:"true"`
	SettlementChainWSSURL string `long:"settlement-chain-wss-url" env:"SETTLEMENT_CHAIN_WSS_URL" description:"Settlement chain WebSocket URL" required:"true"`

	PayerReportManagerAddress string `long:"payer-report-manager-address" env:"PAYER_REPORT_MANAGER_ADDRESS" description:"PayerReportManager contract address"`
	PayerRegistryAddress      string `long:"payer-registry-address" env:"PAYER_REGISTRY_ADDRESS" description:"PayerRegistry contract address"`

	DeploymentBlock        uint64        `long:"deployment-block" env:"DEPLOYMENT_BLOCK" description:"Block number to start backfill from"`
	MaxChainDisconnectTime time.Duration `long:"max-chain-disconnect-time" env:"MAX_CHAIN_DISCONNECT_TIME" description:"Max time before considering chain disconnected" default:"5m"`
	BackfillBlockPageSize  uint64        `long:"backfill-block-page-size" env:"BACKFILL_BLOCK_PAGE_SIZE" description:"Number of blocks per backfill page" default:"500"`
	ActiveOriginatorWindow time.Duration `long:"active-originator-window" env:"ACTIVE_ORIGINATOR_WINDOW" description:"Sliding window for active originator tracking" default:"150m"`

	MetricsPort string `long:"metrics-port" env:"METRICS_PORT" description:"Port for metrics/health HTTP server" default:"8008"`
}

// envConfig represents the subset of environment JSON config we need.
type envConfig struct {
	PayerReportManager             string `json:"payerReportManager"`
	PayerRegistry                  string `json:"payerRegistry"`
	SettlementChainDeploymentBlock int    `json:"settlementChainDeploymentBlock"`
}

var options ChainWatcherOptions

func main() {
	_, err := flags.Parse(&options)
	if err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) &&
			flagsErr.Type == flags.ErrHelp {
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

	cfg := buildConfig(logger)

	// Metrics setup
	reg := prometheus.NewRegistry()
	chainwatcher.RegisterMetrics(reg)
	go serveMetrics(logger, reg, options.MetricsPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher, err := chainwatcher.New(ctx, logger, cfg)
	if err != nil {
		logger.Fatal("failed to create chain watcher", zap.Error(err))
	}

	if err := watcher.Start(); err != nil {
		logger.Fatal("failed to start chain watcher", zap.Error(err))
	}

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	sig := <-sigCh
	logger.Info(
		"received signal, shutting down",
		zap.String("signal", sig.String()),
	)

	cancel()
	watcher.Stop()
}

// buildConfig creates a chainwatcher.Config from parsed options.
func buildConfig(logger *zap.Logger) chainwatcher.Config {
	cfg := chainwatcher.Config{
		SettlementChainRPCURL: options.SettlementChainRPCURL,
		SettlementChainWSSURL: options.SettlementChainWSSURL,
		MaxChainDisconnectTime: options.MaxChainDisconnectTime,
		BackfillBlockPageSize:  options.BackfillBlockPageSize,
		ActiveOriginatorWindow: options.ActiveOriginatorWindow,
		DeploymentBlock:        options.DeploymentBlock,
	}

	// Load contract addresses from named environment if provided.
	if options.Contracts != "" {
		var env environments.SmartContractEnvironment
		if err := env.UnmarshalFlag(options.Contracts); err != nil {
			logger.Fatal(
				"invalid contracts environment",
				zap.Error(err),
			)
		}
		data, err := environments.GetEnvironmentConfig(env)
		if err != nil {
			logger.Fatal(
				"failed to load environment config",
				zap.Error(err),
			)
		}
		var ec envConfig
		if err := json.Unmarshal(data, &ec); err != nil {
			logger.Fatal(
				"failed to parse environment config",
				zap.Error(err),
			)
		}
		cfg.PayerReportManagerAddress = ec.PayerReportManager
		cfg.PayerRegistryAddress = ec.PayerRegistry
		if ec.SettlementChainDeploymentBlock < 0 {
			logger.Fatal(
				"settlementChainDeploymentBlock cannot be negative",
				zap.Int("value", ec.SettlementChainDeploymentBlock),
			)
		}
		// Only override deployment block from environment config
		// if not explicitly set via flag/env.
		if cfg.DeploymentBlock == 0 {
			cfg.DeploymentBlock = uint64(
				ec.SettlementChainDeploymentBlock,
			)
		}
		logger.Info("loaded contract addresses from environment",
			zap.String("environment", options.Contracts),
			zap.String(
				"payer_report_manager",
				cfg.PayerReportManagerAddress,
			),
			zap.String(
				"payer_registry",
				cfg.PayerRegistryAddress,
			),
			zap.Uint64("deployment_block", cfg.DeploymentBlock),
		)
	} else {
		// Explicit addresses required when not using named environment.
		if options.PayerReportManagerAddress == "" {
			fatal("--payer-report-manager-address is required " +
				"when --contracts-environment is not set")
		}
		if options.PayerRegistryAddress == "" {
			fatal("--payer-registry-address is required " +
				"when --contracts-environment is not set")
		}
		cfg.PayerReportManagerAddress = options.PayerReportManagerAddress
		cfg.PayerRegistryAddress = options.PayerRegistryAddress
	}

	return cfg
}

func serveMetrics(
	logger *zap.Logger,
	reg *prometheus.Registry,
	port string,
) {
	mux := http.NewServeMux()
	mux.Handle(
		"/metrics",
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}),
	)
	mux.HandleFunc(
		"/health",
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		},
	)

	addr := net.JoinHostPort("0.0.0.0", port)
	logger.Info("serving metrics", zap.String("address", addr))
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		logger.Error("metrics server error", zap.Error(err))
	}
}

func fatal(msg string, args ...any) {
	log.Fatalf(msg, args...)
}
