package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	messageApi "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/grpc"
)

// blastSnapshot captures results from a blast test.
type blastSnapshot struct {
	Concurrency int          `json:"concurrency"`
	Connections int          `json:"connections"`
	Duration    string       `json:"duration"`
	Timestamp   string       `json:"timestamp"`
	APIs        []testResult `json:"apis"`
	Aggregate   struct {
		TotalRPS    float64 `json:"total_rps"`
		TotalCount  uint64  `json:"total_count"`
		TotalErrors int     `json:"total_errors"`
		ErrorPct    float64 `json:"error_pct"`
	} `json:"aggregate"`
	RateLimited bool `json:"rate_limited"`
}

// runBlastWorkload sends requests as fast as possible across N goroutines and M connections.
// No rate limiting — pure saturation test to find the server's max throughput.
func runBlastWorkload(cfg *config, concurrency int, duration time.Duration) (*blastSnapshot, error) {
	numConns := max(cfg.Connections, 1)

	fmt.Printf("\n╔══════════════════════════════════════════════════╗\n")
	fmt.Printf("║  BLAST MODE — %d goroutines × %d conns           ║\n", concurrency, numConns)
	fmt.Printf("╚══════════════════════════════════════════════════╝\n")
	fmt.Printf("Duration: %s | No rate limiting — max throughput test\n", duration)
	if cfg.GatewayAddr != "" {
		fmt.Printf("Gateway:  %s (writes route through gateway)\n", cfg.GatewayAddr)
	}
	fmt.Printf("CB traffic mix: 65%% GetInboxIds, 14%% QueryEnvelopes, 3%% GetNewest, 3%% Write\n\n")

	// Connection pool
	conns := make([]*grpc.ClientConn, numConns)
	for i := range numConns {
		c, err := newGRPCConn(cfg)
		if err != nil {
			for j := 0; j < i; j++ {
				_ = conns[j].Close()
			}
			return nil, fmt.Errorf("grpc connect [%d]: %w", i, err)
		}
		conns[i] = c
	}
	defer func() {
		for _, c := range conns {
			_ = c.Close()
		}
	}()

	key, err := ethcrypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration+10*time.Second)
	defer cancel()

	deadline := time.Now().Add(duration)

	// Distribute goroutines across APIs using CB traffic fractions
	// (skip SubscribeTopics — 15%)
	type apiAlloc struct {
		name       string
		goroutines int
	}
	adjustedTotal := cbFracGetInboxIds + cbFracQueryEnvelopes + cbFracGetNewest + cbFracWrite
	allocs := []apiAlloc{
		{"GetInboxIds", max(int(math.Round(float64(concurrency)*cbFracGetInboxIds/adjustedTotal)), 1)},
		{"QueryEnvelopes", max(int(math.Round(float64(concurrency)*cbFracQueryEnvelopes/adjustedTotal)), 1)},
		{"GetNewestEnvelope", max(int(math.Round(float64(concurrency)*cbFracGetNewest/adjustedTotal)), 1)},
		{"GroupMessage-256B", max(int(math.Round(float64(concurrency)*cbFracWrite/adjustedTotal)), 1)},
	}

	for _, a := range allocs {
		fmt.Printf("  %-22s → %d goroutines\n", a.name, a.goroutines)
	}
	fmt.Println()

	// Create gateway client if configured (for routing writes through the gateway)
	var gwClient *gatewayClient
	if cfg.GatewayAddr != "" {
		gwClient = newGatewayClient(cfg.GatewayAddr)
	}

	trackers := make(map[string]*apiLatencyTracker)
	var globalRateLimited atomic.Bool
	var wg sync.WaitGroup

	connIdx := 0
	for _, a := range allocs {
		tracker := &apiLatencyTracker{}
		trackers[a.name] = tracker

		for g := range a.goroutines {
			conn := conns[connIdx%numConns]
			connIdx++
			wg.Add(1)
			go func(apiName string, tr *apiLatencyTracker, c *grpc.ClientConn, gIdx int) {
				defer wg.Done()
				blastLoop(ctx, c, cfg, key, apiName, tr, deadline, &globalRateLimited, gwClient)
			}(a.name, tracker, conn, g)
		}
	}

	// Progress reporter
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		start := time.Now()
		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(start).Seconds()
				var totalOK, totalErr uint64
				for _, tr := range trackers {
					totalOK += tr.okCount.Load()
					totalErr += tr.errCount.Load()
				}
				actualRPS := float64(totalOK) / elapsed
				errPct := float64(0)
				if totalOK+totalErr > 0 {
					errPct = float64(totalErr) / float64(totalOK+totalErr) * 100
				}
				fmt.Printf("  [%3.0fs] %d OK | %d err (%.1f%%) | %.0f req/s\n",
					elapsed, totalOK, totalErr, errPct, actualRPS)
				if globalRateLimited.Load() {
					fmt.Println("  ⚠ RATE LIMITED — 429 detected")
				}
			case <-ctx.Done():
				return
			}
			if time.Now().After(deadline) {
				return
			}
		}
	}()

	wg.Wait()
	elapsed := duration

	snap := &blastSnapshot{
		Concurrency: concurrency,
		Connections: numConns,
		Duration:    duration.String(),
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		RateLimited: globalRateLimited.Load(),
	}

	var totalOK, totalErr uint64
	for _, a := range allocs {
		tr := trackers[a.name]
		result := tr.toResult(a.name, elapsed)
		snap.APIs = append(snap.APIs, result)
		totalOK += uint64(result.OKCount)
		totalErr += uint64(result.ErrCount)
	}

	snap.Aggregate.TotalCount = totalOK + totalErr
	snap.Aggregate.TotalRPS = float64(totalOK) / elapsed.Seconds()
	snap.Aggregate.TotalErrors = int(totalErr)
	if totalOK+totalErr > 0 {
		snap.Aggregate.ErrorPct = float64(totalErr) / float64(totalOK+totalErr) * 100
	}

	fmt.Println()
	fmt.Printf("═══ BLAST RESULTS (%d goroutines × %d conns) ═══\n", concurrency, numConns)
	fmt.Printf("Actual: %.0f req/s | Errors: %d (%.1f%%)\n",
		snap.Aggregate.TotalRPS, totalErr, snap.Aggregate.ErrorPct)
	if snap.RateLimited {
		fmt.Println("⚠ RATE LIMITED")
	}
	fmt.Println()

	fmt.Println("╔════════════════════════╦══════════╦══════════╦══════════╦══════════╦══════════╦════════╗")
	fmt.Println("║ API                    ║ Actual   ║ Avg(ms)  ║ P50(ms)  ║ P99(ms)  ║ StdDev  ║ Err%   ║")
	fmt.Println("╠════════════════════════╬══════════╬══════════╬══════════╬══════════╬══════════╬════════╣")
	for _, r := range snap.APIs {
		fmt.Printf("║ %-22s ║ %8.1f ║ %8.2f ║ %8.2f ║ %8.2f ║ %7.2f ║ %5.1f%% ║\n",
			r.Name, r.RPS, r.AvgLatency, r.P50Latency, r.P99Latency, r.StdDev, r.ErrorPct)
	}
	fmt.Println("╚════════════════════════╩══════════╩══════════╩══════════╩══════════╩══════════╩════════╝")

	return snap, nil
}

// blastLoop sends requests as fast as possible with no rate limiting.
func blastLoop(
	ctx context.Context,
	conn *grpc.ClientConn,
	cfg *config,
	key *ecdsa.PrivateKey,
	apiName string,
	tracker *apiLatencyTracker,
	deadline time.Time,
	rateLimited *atomic.Bool,
	gwClient *gatewayClient,
) {
	replicationClient := messageApi.NewReplicationApiClient(conn)
	publishClient := messageApi.NewPublishApiClient(conn)

	topicID := randomBytes(16)
	writeTopic := topic.NewTopic(topic.TopicKindGroupMessagesV1, topicID)

	for {
		if time.Now().After(deadline) {
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
		}

		start := time.Now()
		var callErr error

		switch apiName {
		case "GetInboxIds":
			_, callErr = replicationClient.GetInboxIds(ctx, &messageApi.GetInboxIdsRequest{
				Requests: []*messageApi.GetInboxIdsRequest_Request{{
					Identifier:     "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
					IdentifierKind: associations.IdentifierKind_IDENTIFIER_KIND_ETHEREUM,
				}},
			})
		case "QueryEnvelopes":
			_, callErr = replicationClient.QueryEnvelopes(ctx, &messageApi.QueryEnvelopesRequest{
				Query: &messageApi.EnvelopesQuery{OriginatorNodeIds: []uint32{100}},
				Limit: 5,
			})
		case "GetNewestEnvelope":
			_, callErr = replicationClient.GetNewestEnvelope(ctx, &messageApi.GetNewestEnvelopeRequest{
				Topics: [][]byte{[]byte("AAAAAAAAAAAAAAAAAAAAAA==")},
			})
		case "GroupMessage-256B":
			if gwClient != nil {
				// Route writes through gateway (CB Wallet production path)
				callErr = gwClient.publishClientEnvelope(
					ctx, topic.TopicKindGroupMessagesV1, writeTopic, 256,
				)
			} else {
				// Direct to node (bypass gateway)
				req, buildErr := buildPublishRequestForTopic(cfg, key, writeTopic, 256)
				if buildErr != nil {
					tracker.recordErr("build: " + buildErr.Error())
					continue
				}
				_, callErr = publishClient.PublishPayerEnvelopes(ctx, req)
			}
		}

		latency := time.Since(start)

		if callErr != nil {
			errMsg := callErr.Error()
			tracker.recordErr(errMsg)
			if contains429(errMsg) {
				rateLimited.Store(true)
			}
		} else {
			tracker.recordOK(latency)
		}
	}
}

// Helper reused from mix.go
func blastPercentile(sorted []time.Duration, pct int) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := max(int(math.Ceil(float64(pct)/100.0*float64(len(sorted))))-1, 0)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

