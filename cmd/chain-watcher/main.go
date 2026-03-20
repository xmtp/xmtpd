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
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xmtp/xmtpd/pkg/chainwatcher"
	"github.com/xmtp/xmtpd/pkg/config/environments"
	"go.uber.org/zap"
)

// envConfig represents the subset of environment JSON config we need.
type envConfig struct {
	PayerReportManager             string `json:"payerReportManager"`
	PayerRegistry                  string `json:"payerRegistry"`
	SettlementChainDeploymentBlock int    `json:"settlementChainDeploymentBlock"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	cfg := buildConfig(logger)

	// Metrics setup
	reg := prometheus.NewRegistry()
	chainwatcher.RegisterMetrics(reg)

	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "8008"
	}
	go serveMetrics(logger, reg, metricsPort)

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
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	logger.Info("received signal, shutting down", zap.String("signal", sig.String()))

	cancel()
	watcher.Stop()
}

// buildConfig creates a Config from environment variables.
// Supports two modes:
//  1. Explicit: PAYER_REPORT_MANAGER_ADDRESS + PAYER_REGISTRY_ADDRESS
//  2. Environment-based: XMTPD_CONTRACTS_ENVIRONMENT=testnet-staging (auto-loads addresses)
func buildConfig(logger *zap.Logger) chainwatcher.Config {
	cfg := chainwatcher.Config{
		SettlementChainRPCURL: requireEnv("SETTLEMENT_CHAIN_RPC_URL"),
		SettlementChainWSSURL: requireEnv("SETTLEMENT_CHAIN_WSS_URL"),
	}

	// Try environment-based config first
	if envName := os.Getenv("XMTPD_CONTRACTS_ENVIRONMENT"); envName != "" {
		var env environments.SmartContractEnvironment
		if err := env.UnmarshalFlag(envName); err != nil {
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
		cfg.PayerReportManagerAddress = ec.PayerReportManager
		cfg.PayerRegistryAddress = ec.PayerRegistry
		if ec.SettlementChainDeploymentBlock < 0 {
			logger.Fatal("settlementChainDeploymentBlock cannot be negative",
				zap.Int("value", ec.SettlementChainDeploymentBlock))
		}
		cfg.DeploymentBlock = uint64(ec.SettlementChainDeploymentBlock)
		logger.Info("loaded contract addresses from environment",
			zap.String("environment", envName),
			zap.String("payer_report_manager", cfg.PayerReportManagerAddress),
			zap.String("payer_registry", cfg.PayerRegistryAddress),
			zap.Uint64("deployment_block", cfg.DeploymentBlock),
		)
	} else {
		// Explicit addresses
		cfg.PayerReportManagerAddress = requireEnv("PAYER_REPORT_MANAGER_ADDRESS")
		cfg.PayerRegistryAddress = requireEnv("PAYER_REGISTRY_ADDRESS")
	}

	// Optional overrides
	if v := os.Getenv("DEPLOYMENT_BLOCK"); v != "" {
		block, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			logger.Fatal("failed to parse DEPLOYMENT_BLOCK", zap.String("value", v), zap.Error(err))
		}
		cfg.DeploymentBlock = block
	}

	if v := os.Getenv("MAX_CHAIN_DISCONNECT_TIME"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			logger.Fatal("failed to parse MAX_CHAIN_DISCONNECT_TIME", zap.String("value", v), zap.Error(err))
		}
		cfg.MaxChainDisconnectTime = d
	}

	if v := os.Getenv("BACKFILL_BLOCK_PAGE_SIZE"); v != "" {
		size, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			logger.Fatal("failed to parse BACKFILL_BLOCK_PAGE_SIZE", zap.String("value", v), zap.Error(err))
		}
		cfg.BackfillBlockPageSize = size
	}

	if v := os.Getenv("ACTIVE_ORIGINATOR_WINDOW"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			logger.Fatal("failed to parse ACTIVE_ORIGINATOR_WINDOW", zap.String("value", v), zap.Error(err))
		}
		cfg.ActiveOriginatorWindow = d
	}

	return cfg
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "required environment variable %s is not set\n", key)
		os.Exit(1)
	}
	return v
}

func serveMetrics(logger *zap.Logger, reg *prometheus.Registry, port string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
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
