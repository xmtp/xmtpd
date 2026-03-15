package chain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	reportManager "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/config/environments"
	"go.uber.org/zap"
)

// OnChainPayerReport is a simplified view of on-chain payer report state
// for use in E2E test assertions.
type OnChainPayerReport struct {
	StartSequenceID     uint64
	EndSequenceID       uint64
	EndMinuteSinceEpoch uint32
	IsSettled           bool
	FeesSettled         *big.Int
	Offset              uint32
	PayersMerkleRoot    [32]byte
	NodeIDs             []uint32
}

// Contracts provides read-only access to on-chain contract state for E2E tests.
// It wraps an ethclient connection and ABI bindings for the PayerReportManager
// contract deployed on the Anvil chain.
type Contracts struct {
	logger        *zap.Logger
	client        *ethclient.Client
	reportManager *reportManager.PayerReportManagerCaller
}

// NewContracts creates a Contracts reader connected to the given RPC URL.
// Contract addresses are loaded from the embedded Anvil configuration.
func NewContracts(
	ctx context.Context,
	logger *zap.Logger,
	rpcURL string,
) (*Contracts, error) {
	chainConfig, err := loadAnvilConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load anvil config: %w", err)
	}

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to chain at %s: %w", rpcURL, err)
	}

	addr := common.HexToAddress(chainConfig.PayerReportManager)
	rm, err := reportManager.NewPayerReportManagerCaller(addr, client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf(
			"failed to bind PayerReportManager at %s: %w",
			addr.Hex(), err,
		)
	}

	logger.Info("contracts reader initialized",
		zap.String("rpc_url", rpcURL),
		zap.String("payer_report_manager", addr.Hex()),
	)

	return &Contracts{
		logger:        logger,
		client:        client,
		reportManager: rm,
	}, nil
}

// Close releases the ethclient connection.
func (c *Contracts) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

// GetPayerReport reads a single payer report from the chain by originator
// node ID and report index.
func (c *Contracts) GetPayerReport(
	ctx context.Context,
	originatorNodeID uint32,
	index uint64,
) (*OnChainPayerReport, error) {
	report, err := c.reportManager.GetPayerReport(
		&bind.CallOpts{Context: ctx},
		originatorNodeID,
		new(big.Int).SetUint64(index),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get payer report (node=%d, index=%d): %w",
			originatorNodeID, index, err,
		)
	}

	return &OnChainPayerReport{
		StartSequenceID:     report.StartSequenceId,
		EndSequenceID:       report.EndSequenceId,
		EndMinuteSinceEpoch: report.EndMinuteSinceEpoch,
		IsSettled:           report.IsSettled,
		FeesSettled:         report.FeesSettled,
		Offset:              report.Offset,
		PayersMerkleRoot:    report.PayersMerkleRoot,
		NodeIDs:             report.NodeIds,
	}, nil
}

// WaitForSettledReport polls the chain until checkFn returns true for a payer
// report at the given originator node ID and index, or the context is cancelled.
// Use this to wait for settlement to complete on-chain.
func (c *Contracts) WaitForSettledReport(
	ctx context.Context,
	originatorNodeID uint32,
	index uint64,
	checkFn func(*OnChainPayerReport) bool,
	description string,
) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	var lastReport *OnChainPayerReport

	for {
		report, err := c.GetPayerReport(ctx, originatorNodeID, index)
		if err != nil {
			c.logger.Warn("error checking on-chain report",
				zap.Error(err),
			)
		} else {
			lastReport = report
			if checkFn(report) {
				c.logger.Info("on-chain report condition met",
					zap.String("condition", description),
					zap.Uint64("start_seq", report.StartSequenceID),
					zap.Uint64("end_seq", report.EndSequenceID),
					zap.Bool("is_settled", report.IsSettled),
				)
				return nil
			}
		}

		select {
		case <-ctx.Done():
			lastState := ""
			if lastReport != nil {
				lastState = fmt.Sprintf(
					" (last: start_seq=%d, end_seq=%d, settled=%v, offset=%d)",
					lastReport.StartSequenceID,
					lastReport.EndSequenceID,
					lastReport.IsSettled,
					lastReport.Offset,
				)
			}
			return fmt.Errorf(
				"timed out waiting for on-chain report (%s)%s: %w",
				description, lastState, ctx.Err(),
			)
		case <-ticker.C:
		}
	}
}

// loadAnvilConfig parses the embedded Anvil environment JSON to extract
// contract addresses.
func loadAnvilConfig() (*config.ChainConfig, error) {
	data, err := environments.GetEnvironmentConfig(environments.Anvil)
	if err != nil {
		return nil, err
	}

	var cfg config.ChainConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse anvil config: %w", err)
	}

	return &cfg, nil
}
