package chainwatcher

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	payerregistryabi "github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	payerreportmanagerabi "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	rpcstreamer "github.com/xmtp/xmtpd/pkg/indexer/rpc_streamer"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

const (
	payerReportManagerContractID = "chain-watcher-payer-report-manager"
	payerRegistryContractID      = "chain-watcher-payer-registry"

	// Default sliding window for active originator tracking.
	// Should be >= 2x the expected payer report worker cycle time
	// (production default is 60 min, so 150 min gives comfortable margin).
	defaultActiveOriginatorWindow = 150 * time.Minute

	// Interval for cleaning up stale entries from in-memory maps.
	staleEntryCleanupInterval = 10 * time.Minute

	// Maximum age for submission time tracking entries before cleanup.
	// Reports not settled within this window are considered stuck.
	maxSubmissionTrackingAge = 4 * time.Hour
)

// Config holds the configuration for the chain watcher.
type Config struct {
	SettlementChainRPCURL string
	SettlementChainWSSURL string

	PayerReportManagerAddress string
	PayerRegistryAddress      string

	DeploymentBlock        uint64
	MaxChainDisconnectTime time.Duration
	BackfillBlockPageSize  uint64
	ActiveOriginatorWindow time.Duration
}

// Watcher subscribes to settlement chain events and emits Prometheus metrics.
type Watcher struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	logger *zap.Logger

	rpcClient *ethclient.Client
	wsClient  *ethclient.Client

	payerReportManagerContract *payerreportmanagerabi.PayerReportManager
	payerRegistryContract      *payerregistryabi.PayerRegistry

	payerReportManagerABI *abi.ABI
	payerRegistryABI      *abi.ABI

	streamer *rpcstreamer.RPCLogStreamer

	activeOriginatorWindow time.Duration

	// State tracking
	mu                  sync.RWMutex
	lastSubmissionTime  time.Time
	lastEndSeqByNode    map[uint32]uint64    // track envelope gaps
	submissionTimeByKey map[string]time.Time // track submission→settlement latency
	activeOriginators   map[uint32]time.Time // sliding window
}

// New creates a new chain watcher.
func New(ctx context.Context, logger *zap.Logger, cfg Config) (*Watcher, error) {
	ctx, cancel := context.WithCancel(ctx)
	success := false
	defer func() {
		if !success {
			cancel()
		}
	}()

	rpcClient, err := blockchain.NewRPCClient(ctx, cfg.SettlementChainRPCURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		if !success {
			rpcClient.Close()
		}
	}()

	wsClient, err := blockchain.NewWebsocketClient(ctx, cfg.SettlementChainWSSURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		if !success {
			wsClient.Close()
		}
	}()

	payerReportManagerAddr := common.HexToAddress(cfg.PayerReportManagerAddress)
	payerRegistryAddr := common.HexToAddress(cfg.PayerRegistryAddress)

	prmContract, err := payerreportmanagerabi.NewPayerReportManager(
		payerReportManagerAddr,
		rpcClient,
	)
	if err != nil {
		return nil, err
	}

	prContract, err := payerregistryabi.NewPayerRegistry(
		payerRegistryAddr,
		rpcClient,
	)
	if err != nil {
		return nil, err
	}

	prmABI, err := payerreportmanagerabi.PayerReportManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	prABI, err := payerregistryabi.PayerRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	prmTopics := []common.Hash{
		prmABI.Events["PayerReportSubmitted"].ID,
		prmABI.Events["PayerReportSubsetSettled"].ID,
	}

	prTopics := []common.Hash{
		prABI.Events["UsageSettled"].ID,
	}

	maxDisconnect := cfg.MaxChainDisconnectTime
	if maxDisconnect == 0 {
		maxDisconnect = 5 * time.Minute
	}

	pageSize := cfg.BackfillBlockPageSize
	if pageSize == 0 {
		pageSize = 500
	}

	originatorWindow := cfg.ActiveOriginatorWindow
	if originatorWindow == 0 {
		originatorWindow = defaultActiveOriginatorWindow
	}

	streamer, err := rpcstreamer.NewRPCLogStreamer(
		ctx,
		rpcClient,
		wsClient,
		logger.Named("streamer"),
		rpcstreamer.WithContractConfig(&rpcstreamer.ContractConfig{
			ID:                payerReportManagerContractID,
			FromBlockNumber:   cfg.DeploymentBlock,
			Address:           payerReportManagerAddr,
			Topics:            prmTopics,
			MaxDisconnectTime: maxDisconnect,
		}),
		rpcstreamer.WithContractConfig(&rpcstreamer.ContractConfig{
			ID:                payerRegistryContractID,
			FromBlockNumber:   cfg.DeploymentBlock,
			Address:           payerRegistryAddr,
			Topics:            prTopics,
			MaxDisconnectTime: maxDisconnect,
		}),
		rpcstreamer.WithBackfillBlockPageSize(pageSize),
	)
	if err != nil {
		return nil, err
	}

	success = true
	return &Watcher{
		ctx:                        ctx,
		cancel:                     cancel,
		logger:                     logger.Named("chain-watcher"),
		rpcClient:                  rpcClient,
		wsClient:                   wsClient,
		payerReportManagerContract: prmContract,
		payerRegistryContract:      prContract,
		payerReportManagerABI:      prmABI,
		payerRegistryABI:           prABI,
		streamer:                   streamer,
		activeOriginatorWindow:     originatorWindow,
		lastEndSeqByNode:           make(map[uint32]uint64),
		submissionTimeByKey:        make(map[string]time.Time),
		activeOriginators:          make(map[uint32]time.Time),
	}, nil
}

// Start begins watching chain events and emitting metrics.
func (w *Watcher) Start() error {
	if err := w.streamer.Start(); err != nil {
		return err
	}

	prmChan := w.streamer.GetEventChannel(payerReportManagerContractID)
	prChan := w.streamer.GetEventChannel(payerRegistryContractID)

	tracing.GoPanicWrap(w.ctx, &w.wg, "payer-report-manager-events", func(ctx context.Context) {
		w.processPayerReportManagerEvents(ctx, prmChan)
	})

	tracing.GoPanicWrap(w.ctx, &w.wg, "payer-registry-events", func(ctx context.Context) {
		w.processPayerRegistryEvents(ctx, prChan)
	})

	tracing.GoPanicWrap(w.ctx, &w.wg, "submission-lag-ticker", func(ctx context.Context) {
		w.runSubmissionLagTicker(ctx)
	})

	tracing.GoPanicWrap(w.ctx, &w.wg, "stale-entry-cleanup", func(ctx context.Context) {
		w.runStaleEntryCleanup(ctx)
	})

	w.logger.Info("chain watcher started")
	return nil
}

// Stop gracefully shuts down the watcher.
func (w *Watcher) Stop() {
	w.logger.Info("stopping chain watcher")
	w.streamer.Stop()
	w.rpcClient.Close()
	w.wsClient.Close()
	w.cancel()
	w.wg.Wait()
	w.logger.Info("chain watcher stopped")
}

func (w *Watcher) processPayerReportManagerEvents(
	ctx context.Context,
	ch <-chan types.Log,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case log, ok := <-ch:
			if !ok {
				return
			}
			if len(log.Topics) == 0 {
				continue
			}
			event, err := w.payerReportManagerABI.EventByID(log.Topics[0])
			if err != nil {
				// Events from other contracts sharing the address — skip silently.
				continue
			}
			switch event.Name {
			case "PayerReportSubmitted":
				w.handlePayerReportSubmitted(log)
			case "PayerReportSubsetSettled":
				w.handlePayerReportSubsetSettled(log)
			}
		}
	}
}

func (w *Watcher) processPayerRegistryEvents(
	ctx context.Context,
	ch <-chan types.Log,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case log, ok := <-ch:
			if !ok {
				return
			}
			if len(log.Topics) == 0 {
				continue
			}
			event, err := w.payerRegistryABI.EventByID(log.Topics[0])
			if err != nil {
				continue
			}
			switch event.Name {
			case "UsageSettled":
				w.handleUsageSettled(log)
			}
		}
	}
}

func (w *Watcher) handlePayerReportSubmitted(log types.Log) {
	parsed, err := w.payerReportManagerContract.ParsePayerReportSubmitted(log)
	if err != nil {
		w.logger.Error("failed to parse PayerReportSubmitted", zap.Error(err))
		return
	}

	nodeID := parsed.OriginatorNodeId
	nodeLabel := nodeIDLabel(nodeID)
	now := time.Now()

	// Core counter
	reportSubmittedTotal.WithLabelValues(nodeLabel).Inc()
	eventsProcessedTotal.WithLabelValues("PayerReportSubmitted").Inc()

	// Envelope range
	if parsed.EndSequenceId < parsed.StartSequenceId {
		w.logger.Error("invalid envelope range: end < start",
			zap.Uint64("start_seq", parsed.StartSequenceId),
			zap.Uint64("end_seq", parsed.EndSequenceId),
		)
		return
	}
	envelopeCount := parsed.EndSequenceId - parsed.StartSequenceId
	envelopeRangeTotal.WithLabelValues(nodeLabel).Add(float64(envelopeCount))

	// Envelope gap detection
	w.mu.Lock()
	if lastEnd, exists := w.lastEndSeqByNode[nodeID]; exists {
		gap := int64(parsed.StartSequenceId) - int64(lastEnd)
		envelopeRangeGap.WithLabelValues(nodeLabel).Set(float64(gap))
	}
	w.lastEndSeqByNode[nodeID] = parsed.EndSequenceId

	// Track submission time for settlement latency
	key := submissionKey(nodeID, parsed.PayerReportIndex)
	w.submissionTimeByKey[key] = now

	// Active originator tracking
	w.activeOriginators[nodeID] = now
	w.lastSubmissionTime = now
	w.mu.Unlock()

	// Attesting node count
	attestingNodeCount.WithLabelValues(nodeLabel).Set(float64(len(parsed.SigningNodeIds)))

	// Update active originator gauge
	w.updateActiveOriginatorCount()

	w.logger.Info("PayerReportSubmitted",
		zap.Uint32("originator_node_id", nodeID),
		zap.Uint64("start_seq", parsed.StartSequenceId),
		zap.Uint64("end_seq", parsed.EndSequenceId),
		zap.Uint64("envelope_count", envelopeCount),
		zap.Int("signing_nodes", len(parsed.SigningNodeIds)),
		zap.Uint32s("node_ids", parsed.NodeIds),
	)
}

func (w *Watcher) handlePayerReportSubsetSettled(log types.Log) {
	parsed, err := w.payerReportManagerContract.ParsePayerReportSubsetSettled(log)
	if err != nil {
		w.logger.Error("failed to parse PayerReportSubsetSettled", zap.Error(err))
		return
	}

	nodeID := parsed.OriginatorNodeId
	nodeLabel := nodeIDLabel(nodeID)

	eventsProcessedTotal.WithLabelValues("PayerReportSubsetSettled").Inc()

	// Track fees — use big.Float conversion to avoid int64 overflow on uint96 values.
	if parsed.FeesSettled != nil && parsed.FeesSettled.Sign() > 0 {
		fees, _ := new(big.Float).SetInt(parsed.FeesSettled).Float64()
		feesSettledPicodollars.WithLabelValues(nodeLabel).Add(fees)
	}

	// Only count as fully settled when remaining == 0
	if parsed.Remaining == 0 {
		reportSettledTotal.WithLabelValues(nodeLabel).Inc()

		// Calculate submission → settlement latency
		key := submissionKey(nodeID, parsed.PayerReportIndex)
		w.mu.Lock()
		if subTime, exists := w.submissionTimeByKey[key]; exists {
			latency := time.Since(subTime).Seconds()
			submissionToSettlementSeconds.WithLabelValues(nodeLabel).Observe(latency)
			delete(w.submissionTimeByKey, key)
		}
		w.mu.Unlock()

		w.logger.Info("PayerReportSubsetSettled (fully settled)",
			zap.Uint32("originator_node_id", nodeID),
			zap.String("payer_report_index", parsed.PayerReportIndex.String()),
			zap.Uint32("count", parsed.Count),
			zap.String("fees_settled", parsed.FeesSettled.String()),
		)
	} else {
		w.logger.Info("PayerReportSubsetSettled (partial)",
			zap.Uint32("originator_node_id", nodeID),
			zap.String("payer_report_index", parsed.PayerReportIndex.String()),
			zap.Uint32("count", parsed.Count),
			zap.Uint32("remaining", parsed.Remaining),
			zap.String("fees_settled", parsed.FeesSettled.String()),
		)
	}
}

func (w *Watcher) handleUsageSettled(log types.Log) {
	parsed, err := w.payerRegistryContract.ParseUsageSettled(log)
	if err != nil {
		w.logger.Error("failed to parse UsageSettled", zap.Error(err))
		return
	}

	usageSettledTotal.Inc()
	eventsProcessedTotal.WithLabelValues("UsageSettled").Inc()

	w.logger.Debug("UsageSettled",
		zap.String("payer", parsed.Payer.Hex()),
		zap.String("amount", parsed.Amount.String()),
	)
}

func (w *Watcher) runSubmissionLagTicker(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.mu.RLock()
			if !w.lastSubmissionTime.IsZero() {
				lag := time.Since(w.lastSubmissionTime).Seconds()
				timeSinceLastSubmissionSeconds.Set(lag)
			}
			w.mu.RUnlock()
		}
	}
}

func (w *Watcher) updateActiveOriginatorCount() {
	w.mu.RLock()
	defer w.mu.RUnlock()
	cutoff := time.Now().Add(-w.activeOriginatorWindow)
	count := 0
	for _, lastSeen := range w.activeOriginators {
		if lastSeen.After(cutoff) {
			count++
		}
	}
	activeOriginatorNodes.Set(float64(count))
}

// runStaleEntryCleanup periodically removes expired entries from in-memory maps
// to prevent unbounded memory growth (e.g., reports that never settle).
func (w *Watcher) runStaleEntryCleanup(ctx context.Context) {
	ticker := time.NewTicker(staleEntryCleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.cleanupStaleEntries()
		}
	}
}

func (w *Watcher) cleanupStaleEntries() {
	now := time.Now()
	w.mu.Lock()
	defer w.mu.Unlock()

	// Clean up submission tracking for reports that never settled.
	for key, t := range w.submissionTimeByKey {
		if now.Sub(t) > maxSubmissionTrackingAge {
			delete(w.submissionTimeByKey, key)
		}
	}

	// Clean up expired active originator entries.
	cutoff := now.Add(-w.activeOriginatorWindow)
	for nodeID, lastSeen := range w.activeOriginators {
		if lastSeen.Before(cutoff) {
			delete(w.activeOriginators, nodeID)
		}
	}
}

func submissionKey(nodeID uint32, reportIndex *big.Int) string {
	return nodeIDLabel(nodeID) + ":" + reportIndex.String()
}
