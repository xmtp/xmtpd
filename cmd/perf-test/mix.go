package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	messageApi "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/grpc"
)

// CB Wallet user traffic profile (fraction of total API calls).
const (
	cbFracGetInboxIds      = 0.65
	cbFracQueryEnvelopes   = 0.14
	cbFracGetNewest        = 0.03 // subset of query/sync
	cbFracSubscribe        = 0.15
	cbFracWrite            = 0.03
	cbCallsPerUserPerDay   = 470.0
	cbSecondsPerDay        = 86400.0
	cbMsgsPerUserPerDay    = 3.4
	cbAvgSessionStreamSecs = 60.0
)

// mixAPIWorker defines a single API workload in the mix.
type mixAPIWorker struct {
	Name      string
	TargetRPS float64
	IsWrite   bool
	IsStream  bool
}

// mixSnapshot captures per-API metrics for a single DAU tier.
type mixSnapshot struct {
	DAU       int          `json:"dau"`
	TargetRPS float64      `json:"target_rps_total"`
	Duration  string       `json:"duration"`
	Timestamp string       `json:"timestamp"`
	APIs      []testResult `json:"apis"`
	Aggregate struct {
		TotalRPS    float64 `json:"total_rps"`
		TotalCount  uint64  `json:"total_count"`
		TotalErrors int     `json:"total_errors"`
		ErrorPct    float64 `json:"error_pct"`
	} `json:"aggregate"`
	RateLimited bool   `json:"rate_limited"`
	StoppedAt   string `json:"stopped_at,omitempty"`
}

// dauToRPS converts a DAU count to aggregate req/s.
func dauToRPS(dau int) float64 {
	return float64(dau) * cbCallsPerUserPerDay / cbSecondsPerDay
}

// buildMixWorkers creates the API workload distribution for a given DAU.
func buildMixWorkers(dau int) []mixAPIWorker {
	totalRPS := dauToRPS(dau)
	return []mixAPIWorker{
		{Name: "GetInboxIds", TargetRPS: totalRPS * cbFracGetInboxIds},
		{Name: "QueryEnvelopes", TargetRPS: totalRPS * cbFracQueryEnvelopes},
		{Name: "GetNewestEnvelope", TargetRPS: totalRPS * cbFracGetNewest},
		// Skip live subscribe (broken on staging) — note in results
		{Name: "GroupMessage-256B", TargetRPS: totalRPS * cbFracWrite, IsWrite: true},
	}
}

// apiLatencyTracker collects per-request latencies for one API.
type apiLatencyTracker struct {
	mu        sync.Mutex
	latencies []time.Duration
	okCount   atomic.Uint64
	errCount  atomic.Uint64
	errors    sync.Map // error string -> *atomic.Int64
	is429     atomic.Bool
}

func (t *apiLatencyTracker) recordOK(d time.Duration) {
	t.okCount.Add(1)
	t.mu.Lock()
	t.latencies = append(t.latencies, d)
	t.mu.Unlock()
}

func (t *apiLatencyTracker) recordErr(errMsg string) {
	t.errCount.Add(1)
	v, _ := t.errors.LoadOrStore(errMsg, &atomic.Int64{})
	v.(*atomic.Int64).Add(1)
}

func (t *apiLatencyTracker) toResult(name string, elapsed time.Duration) testResult {
	t.mu.Lock()
	lats := make([]time.Duration, len(t.latencies))
	copy(lats, t.latencies)
	t.mu.Unlock()

	ok := int(t.okCount.Load())
	errs := int(t.errCount.Load())
	total := uint64(ok + errs)

	result := testResult{
		Name:     name,
		Count:    total,
		RPS:      float64(ok) / elapsed.Seconds(),
		OKCount:  ok,
		ErrCount: errs,
		Errors:   make(map[string]int),
	}

	t.errors.Range(func(key, value any) bool {
		msg := key.(string)
		cnt := value.(*atomic.Int64).Load()
		if len(msg) > 120 {
			msg = msg[:120] + "..."
		}
		result.Errors[msg] = int(cnt)
		return true
	})

	if ok+errs > 0 {
		result.ErrorPct = float64(errs) / float64(ok+errs) * 100
	}

	if len(lats) > 0 {
		slices.Sort(lats)
		var sum float64
		for _, l := range lats {
			sum += msFromDuration(l)
		}
		result.AvgLatency = sum / float64(len(lats))
		result.P50Latency = msFromDuration(percentile(lats, 50))
		result.P95Latency = msFromDuration(percentile(lats, 95))
		result.P99Latency = msFromDuration(percentile(lats, 99))

		var sumSq float64
		for _, l := range lats {
			diff := msFromDuration(l) - result.AvgLatency
			sumSq += diff * diff
		}
		result.StdDev = math.Sqrt(sumSq / float64(len(lats)))
	}

	return result
}

// runMixedWorkload runs the CB user simulation at a given DAU level.
func runMixedWorkload(cfg *config, dau int, duration time.Duration) (*mixSnapshot, error) {
	workers := buildMixWorkers(dau)
	totalRPS := dauToRPS(dau)

	fmt.Printf("\n╔══════════════════════════════════════════════════╗\n")
	fmt.Printf("║  CB User Mixed Workload — %6d DAU             ║\n", dau)
	fmt.Printf("╚══════════════════════════════════════════════════╝\n")
	fmt.Printf("Target: %.0f req/s total | Duration: %s\n", totalRPS, duration)
	for _, w := range workers {
		fmt.Printf("  %-22s → %6.1f req/s\n", w.Name, w.TargetRPS)
	}
	fmt.Printf("  %-22s → %6.1f req/s (SKIPPED — staging bug)\n", "SubscribeTopics", totalRPS*cbFracSubscribe)
	fmt.Println()

	// Connection pool — multiple gRPC connections to avoid HTTP/2 multiplexing bottleneck
	numConns := max(cfg.Connections, 1)
	conns := make([]*grpc.ClientConn, numConns)
	for i := range numConns {
		c, err := newGRPCConn(cfg)
		if err != nil {
			// Close any already-opened connections
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
	fmt.Printf("Connections: %d\n", numConns)

	// Signing key for writes
	key, err := ethcrypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration+10*time.Second)
	defer cancel()

	deadline := time.Now().Add(duration)
	trackers := make(map[string]*apiLatencyTracker)
	var wg sync.WaitGroup

	// Global 429 detector — if any API hits 429, flag it
	var globalRateLimited atomic.Bool

	connIdx := 0
	for _, w := range workers {
		if w.TargetRPS < 0.5 {
			continue // skip negligible load
		}

		tracker := &apiLatencyTracker{}
		trackers[w.Name] = tracker

		// Calculate goroutine count: 1 goroutine per ~50 req/s, min 1
		numGoroutines := max(int(math.Ceil(w.TargetRPS/50)), 1)
		perGoroutineRPS := w.TargetRPS / float64(numGoroutines)

		for g := range numGoroutines {
			// Round-robin connections across goroutines
			conn := conns[connIdx%numConns]
			connIdx++
			wg.Add(1)
			go func(wk mixAPIWorker, tr *apiLatencyTracker, grps float64, gIdx int, c *grpc.ClientConn) {
				defer wg.Done()
				runMixAPILoop(ctx, c, cfg, key, wk, tr, grps, deadline, &globalRateLimited)
			}(w, tracker, perGoroutineRPS, g, conn)
		}
	}

	// Progress reporter
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(15 * time.Second)
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
				fmt.Printf("  [%3.0fs] %d OK | %d err | %.0f actual req/s (target %.0f)\n",
					elapsed, totalOK, totalErr, actualRPS, totalRPS)
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

	// Build snapshot
	snap := &mixSnapshot{
		DAU:         dau,
		TargetRPS:   totalRPS,
		Duration:    duration.String(),
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		RateLimited: globalRateLimited.Load(),
	}

	var totalOK, totalErr uint64
	for _, w := range workers {
		tr, ok := trackers[w.Name]
		if !ok {
			continue
		}
		result := tr.toResult(w.Name, elapsed)
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

	// Print summary
	fmt.Println()
	fmt.Printf("═══ %dK DAU Results ═══\n", dau/1000)
	fmt.Printf("Target: %.0f req/s | Actual: %.0f req/s | Errors: %d (%.1f%%)\n",
		totalRPS, snap.Aggregate.TotalRPS, totalErr, snap.Aggregate.ErrorPct)
	if snap.RateLimited {
		fmt.Println("⚠ RATE LIMITED during this tier")
	}
	fmt.Println()

	const (
		mixTop    = "╔════════════════════════╦════════╦══════════╦══════════╦══════════╦══════════╦══════════╦════════╗"
		mixHeader = "║ API                    ║ Target ║ Actual   ║ Avg(ms)  ║ P50(ms)  ║ P99(ms)  ║ StdDev  ║ Err%   ║"
		mixSep    = "╠════════════════════════╬════════╬══════════╬══════════╬══════════╬══════════╬══════════╬════════╣"
		mixBot    = "╚════════════════════════╩════════╩══════════╩══════════╩══════════╩══════════╩══════════╩════════╝"
	)

	fmt.Println(mixTop)
	fmt.Println(mixHeader)
	fmt.Println(mixSep)
	for i, w := range workers {
		if i >= len(snap.APIs) {
			break
		}
		r := snap.APIs[i]
		fmt.Printf("║ %-22s ║ %5.0f  ║ %8.1f ║ %8.2f ║ %8.2f ║ %8.2f ║ %7.2f ║ %5.1f%% ║\n",
			w.Name, w.TargetRPS, r.RPS, r.AvgLatency, r.P50Latency, r.P99Latency, r.StdDev, r.ErrorPct)
	}
	fmt.Println(mixBot)

	return snap, nil
}

// runMixAPILoop runs rate-limited gRPC calls for a single API.
func runMixAPILoop(
	ctx context.Context,
	conn *grpc.ClientConn,
	cfg *config,
	key *ecdsa.PrivateKey,
	w mixAPIWorker,
	tracker *apiLatencyTracker,
	targetRPS float64,
	deadline time.Time,
	rateLimited *atomic.Bool,
) {
	if targetRPS < 0.1 {
		return
	}

	interval := time.Duration(float64(time.Second) / targetRPS)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	replicationClient := messageApi.NewReplicationApiClient(conn)
	publishClient := messageApi.NewPublishApiClient(conn)

	// Pre-generate a write topic for this goroutine
	topicID := randomBytes(16)
	writeTopic := topic.NewTopic(topic.TopicKindGroupMessagesV1, topicID)

	for {
		select {
		case <-ticker.C:
			if time.Now().After(deadline) {
				return
			}

			start := time.Now()
			var callErr error

			switch w.Name {
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
				req, buildErr := buildPublishRequestForTopic(cfg, key, writeTopic, 256)
				if buildErr != nil {
					tracker.recordErr("build: " + buildErr.Error())
					continue
				}
				_, callErr = publishClient.PublishPayerEnvelopes(ctx, req)
			}

			latency := time.Since(start)

			if callErr != nil {
				errMsg := callErr.Error()
				tracker.recordErr(errMsg)
				if contains429(errMsg) {
					rateLimited.Store(true)
					tracker.is429.Store(true)
				}
			} else {
				tracker.recordOK(latency)
			}

		case <-ctx.Done():
			return
		}
	}
}

func contains429(s string) bool {
	return strings.Contains(s, "429") || strings.Contains(s, "Too Many Requests")
}
