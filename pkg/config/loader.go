package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/xmtp/xmtpd/pkg/config/environments"
)

const (
	// maxConfigSize is the maximum allowed size for configuration files (10KB).
	maxConfigSize = 10 << 10
	// httpTimeout is the timeout for HTTP requests when fetching remote configs.
	httpTimeout = 10 * time.Second
)

// ContractsSource specifies how to load contract configuration.
type ContractsSource struct {
	Environment string // Named environment: "testnet", "mainnet"
	FilePath    string // Local path, URL, or config://<env>
	JSONData    string // Raw JSON string
}

// LoadContractsConfig loads contract configuration from the specified source.
// Only one source should be specified; they are mutually exclusive.
func LoadContractsConfig(src ContractsSource) (*ContractsOptions, error) {
	sources := 0
	if src.Environment != "" {
		sources++
	}
	if src.FilePath != "" {
		sources++
	}
	if src.JSONData != "" {
		sources++
	}

	if sources == 0 {
		return nil, errors.New("one of environment, file-path, or json is required")
	}
	if sources > 1 {
		return nil, errors.New("environment, file-path, and json are mutually exclusive")
	}

	var data []byte
	var err error

	switch {
	case src.Environment != "":
		data, err = loadFromEnvironment(src.Environment)
	case src.FilePath != "":
		data, err = loadFromPath(src.FilePath)
	case src.JSONData != "":
		data = []byte(src.JSONData)
	}

	if err != nil {
		return nil, err
	}

	return parseChainConfig(data)
}

func loadFromEnvironment(name string) ([]byte, error) {
	var env environments.SmartContractEnvironment
	if err := env.UnmarshalFlag(name); err != nil {
		return nil, fmt.Errorf("invalid environment %q: %w", name, err)
	}
	return environments.GetEnvironmentConfig(env)
}

func loadFromPath(path string) ([]byte, error) {
	// Handle config:// URL scheme
	if strings.HasPrefix(path, "config://") {
		return loadFromEnvironment(strings.TrimPrefix(path, "config://"))
	}

	// Handle HTTP/HTTPS URLs
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return fetchURL(path)
	}

	// Handle file:// URLs
	if strings.HasPrefix(path, "file://") {
		path = strings.TrimPrefix(path, "file://")
		var err error
		path, err = url.PathUnescape(path)
		if err != nil {
			return nil, fmt.Errorf("invalid file URL path %q: %w", path, err)
		}
	}

	// Local file
	return os.ReadFile(path)
}

func fetchURL(url string) ([]byte, error) {
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", url, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch %s: status %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxConfigSize+1))
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", url, err)
	}
	if len(data) > maxConfigSize {
		return nil, fmt.Errorf("fetch %s: response exceeds %d bytes", url, maxConfigSize)
	}
	return data, nil
}

func parseChainConfig(data []byte) (*ContractsOptions, error) {
	var cfg ChainConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &ContractsOptions{
		AppChain: AppChainOptions{
			GroupMessageBroadcasterAddress:   cfg.GroupMessageBroadcaster,
			IdentityUpdateBroadcasterAddress: cfg.IdentityUpdateBroadcaster,
			GatewayAddress:                   cfg.AppChainGateway,
			ParameterRegistryAddress:         cfg.AppChainParameterRegistry,
			ChainID:                          int64(cfg.AppChainID),
			DeploymentBlock:                  uint64(cfg.AppChainDeploymentBlock),
			MaxChainDisconnectTime:           300 * time.Second,
			BackfillBlockPageSize:            500,
		},
		SettlementChain: SettlementChainOptions{
			NodeRegistryAddress:         cfg.NodeRegistry,
			RateRegistryAddress:         cfg.RateRegistry,
			ParameterRegistryAddress:    cfg.SettlementChainParameterRegistry,
			PayerRegistryAddress:        cfg.PayerRegistry,
			PayerReportManagerAddress:   cfg.PayerReportManager,
			GatewayAddress:              cfg.SettlementChainGateway,
			DistributionManagerAddress:  cfg.DistributionManager,
			UnderlyingFeeToken:          cfg.UnderlyingFeeToken,
			FeeToken:                    cfg.FeeToken,
			ChainID:                     int64(cfg.SettlementChainID),
			DeploymentBlock:             uint64(cfg.SettlementChainDeploymentBlock),
			NodeRegistryRefreshInterval: 60 * time.Second,
			RateRegistryRefreshInterval: 300 * time.Second,
			MaxChainDisconnectTime:      300 * time.Second,
			BackfillBlockPageSize:       500,
		},
	}, nil
}
