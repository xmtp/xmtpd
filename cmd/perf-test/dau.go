package main

// Realistic CB Wallet DAU blast test.
//
// Improvements over basic blast:
// 1. Data cardinality: 10K unique addresses for GetInboxIds (cache-busting)
// 2. Seeded data: publishes to 500 topics before querying them
// 3. Background streaming: concurrent SubscribeAllEnvelopes connections
// 4. Mixed write types: GroupMessage (80%) + WelcomeMessage (20%) through gateway
// 5. Varied payload sizes: 256B/512B/1KB mix
// 6. Topic-based QueryEnvelopes (real production path, not node-ID filter)

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
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

const (
	dauAddressPoolSize = 10000 // unique ETH addresses for GetInboxIds
	dauTopicPoolSize   = 500   // unique topics to seed and query
	dauSeedPerTopic    = 5     // messages per topic during seed phase
)

func generateAddressPool(n int) []string {
	pool := make([]string, n)
	for i := range n {
		pool[i] = "0x" + hex.EncodeToString(randomBytes(20))
	}
	return pool
}

func generateTopicPool(n int, kind topic.TopicKind) []*topic.Topic {
	pool := make([]*topic.Topic, n)
	for i := range n {
		pool[i] = topic.NewTopic(kind, randomBytes(16))
	}
	return pool
}

func runDAUBlast(cfg *config, concurrency int, duration time.Duration) (*blastSnapshot, error) {
	numConns := max(cfg.Connections, 1)
	// Scale background streams with concurrency: ~15% of goroutines, min 10, max 200
	numStreams := max(min(concurrency*15/100, 200), 10)

	fmt.Printf("\n╔══════════════════════════════════════════════════╗\n")
	fmt.Printf("║  REALISTIC DAU BLAST — %d goroutines × %d conns  ║\n", concurrency, numConns)
	fmt.Printf("╚══════════════════════════════════════════════════╝\n")
	fmt.Printf("Duration: %s | Background streams: %d\n", duration, numStreams)
	if cfg.GatewayAddr != "" {
		fmt.Printf("Gateway:  %s (writes route through gateway)\n", cfg.GatewayAddr)
	}
	fmt.Printf("Realistic improvements:\n")
	fmt.Printf("  - %d unique addresses for GetInboxIds (cache-busting)\n", dauAddressPoolSize)
	fmt.Printf("  - %d seeded topics for QueryEnvelopes/GetNewest (real data)\n", dauTopicPoolSize)
	fmt.Printf("  - %d SubscribeAllEnvelopes streams (server resource pressure)\n", numStreams)
	fmt.Printf(
		"  - Mixed writes: GroupMsg 48%% (app→node) + Commits 32%% (→blockchain) + Welcome 15%% + KeyPkg 5%%\n",
	)
	fmt.Printf("  - 40%% of group messages are MLS Commits (routed through blockchain)\n")
	fmt.Printf("  - Key package rotations included\n")
	fmt.Printf("  - Varied payload sizes: 256B/512B/1KB\n")
	fmt.Printf("  - Topic-based queries (not originator node ID)\n")
	fmt.Printf(
		"CB traffic mix: 65%% GetInboxIds, 14%% QueryEnvelopes, 3%% GetNewest, 3%% Write, 15%% Streams\n\n",
	)

	// --- Pre-generate data pools ---
	fmt.Print("Generating data pools... ")
	addresses := generateAddressPool(dauAddressPoolSize)
	topics := generateTopicPool(dauTopicPoolSize, topic.TopicKindGroupMessagesV1)
	fmt.Printf("%d addresses, %d topics\n", len(addresses), len(topics))

	// --- Connection pool ---
	conns := make([]*grpc.ClientConn, numConns)
	for i := range numConns {
		c, err := newGRPCConn(cfg)
		if err != nil {
			for j := range i {
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

	// === Phase 1: Seed topics with real data ===
	fmt.Printf(
		"\n--- Phase 1: Seeding %d topics × %d msgs ---\n",
		dauTopicPoolSize,
		dauSeedPerTopic,
	)
	seedStart := time.Now()
	var seedOK, seedErr atomic.Uint64

	seedCtx, seedCancel := context.WithTimeout(context.Background(), 60*time.Second)
	var seedWg sync.WaitGroup
	seedWorkers := min(16, numConns)
	topicsPerWorker := dauTopicPoolSize / seedWorkers

	for w := range seedWorkers {
		seedWg.Add(1)
		conn := conns[w%numConns]
		go func(wIdx int) {
			defer seedWg.Done()
			pub := messageApi.NewPublishApiClient(conn)
			start := wIdx * topicsPerWorker
			end := start + topicsPerWorker
			if wIdx == seedWorkers-1 {
				end = dauTopicPoolSize
			}
			for t := start; t < end; t++ {
				for range dauSeedPerTopic {
					req, e := buildPublishRequestForTopic(cfg, key, topics[t], 256)
					if e != nil {
						seedErr.Add(1)
						continue
					}
					if _, e = pub.PublishPayerEnvelopes(seedCtx, req); e != nil {
						seedErr.Add(1)
						continue
					}
					seedOK.Add(1)
				}
			}
		}(w)
	}
	seedWg.Wait()
	seedCancel()
	fmt.Printf("Seeded: %d OK, %d errors in %s\n",
		seedOK.Load(), seedErr.Load(), time.Since(seedStart).Round(time.Millisecond))

	// === Phase 2: Open background streaming connections ===
	fmt.Printf("\n--- Phase 2: Opening %d SubscribeAllEnvelopes streams ---\n", numStreams)

	blastCtx, blastCancel := context.WithTimeout(context.Background(), duration+30*time.Second)
	defer blastCancel()
	deadline := time.Now().Add(duration)

	var streamWg sync.WaitGroup
	var streamRecv, streamErrs atomic.Uint64

	for s := range numStreams {
		streamWg.Add(1)
		conn := conns[s%numConns]
		go func() {
			defer streamWg.Done()
			client := messageApi.NewNotificationApiClient(conn)
			stream, serr := client.SubscribeAllEnvelopes(
				blastCtx, &messageApi.SubscribeAllEnvelopesRequest{},
			)
			if serr != nil {
				streamErrs.Add(1)
				return
			}
			for {
				if time.Now().After(deadline) {
					return
				}
				resp, rerr := stream.Recv()
				if rerr != nil {
					if blastCtx.Err() != nil {
						return
					}
					streamErrs.Add(1)
					return
				}
				if resp != nil && len(resp.GetEnvelopes()) > 0 {
					streamRecv.Add(uint64(len(resp.GetEnvelopes())))
				}
			}
		}()
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Streams ready (errors: %d)\n", streamErrs.Load())

	// === Phase 3: Blast with realistic data ===
	fmt.Printf("\n--- Phase 3: Realistic blast ---\n")

	adjustedTotal := cbFracGetInboxIds + cbFracQueryEnvelopes + cbFracGetNewest + cbFracWrite
	type apiAlloc struct {
		name       string
		goroutines int
	}
	allocs := []apiAlloc{
		{
			"GetInboxIds",
			max(int(math.Round(float64(concurrency)*cbFracGetInboxIds/adjustedTotal)), 1),
		},
		{
			"QueryEnvelopes",
			max(int(math.Round(float64(concurrency)*cbFracQueryEnvelopes/adjustedTotal)), 1),
		},
		{
			"GetNewestEnvelope",
			max(int(math.Round(float64(concurrency)*cbFracGetNewest/adjustedTotal)), 1),
		},
		{"Writes-Mixed", max(int(math.Round(float64(concurrency)*cbFracWrite/adjustedTotal)), 1)},
	}

	for _, a := range allocs {
		fmt.Printf("  %-22s -> %d goroutines\n", a.name, a.goroutines)
	}
	fmt.Printf("  %-22s -> %d streams (background)\n", "SubscribeAll", numStreams)
	fmt.Println()

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
		for range a.goroutines {
			conn := conns[connIdx%numConns]
			connIdx++
			wg.Go(func() {
				dauBlastLoop(
					blastCtx, conn, cfg, key, a.name, tracker, deadline,
					&globalRateLimited, gwClient, addresses, topics,
				)
			})
		}
	}

	// Progress reporter
	wg.Go(func() {
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
				rps := float64(totalOK) / elapsed
				errPct := float64(0)
				if totalOK+totalErr > 0 {
					errPct = float64(totalErr) / float64(totalOK+totalErr) * 100
				}
				fmt.Printf(
					"  [%3.0fs] %d OK | %d err (%.1f%%) | %.0f req/s | streams: %d msgs recv'd\n",
					elapsed,
					totalOK,
					totalErr,
					errPct,
					rps,
					streamRecv.Load(),
				)
				if globalRateLimited.Load() {
					fmt.Println("  !! RATE LIMITED -- 429 detected")
				}
			case <-blastCtx.Done():
				return
			}
			if time.Now().After(deadline) {
				return
			}
		}
	})

	wg.Wait()
	blastCancel()
	streamWg.Wait()

	elapsed := duration

	// Build results
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

	// Add stream stats
	snap.APIs = append(snap.APIs, testResult{
		Name:     "SubscribeAll",
		Count:    streamRecv.Load(),
		RPS:      float64(streamRecv.Load()) / elapsed.Seconds(),
		OKCount:  int(streamRecv.Load()),
		ErrCount: int(streamErrs.Load()),
		Errors:   make(map[string]int),
	})

	snap.Aggregate.TotalCount = totalOK + totalErr
	snap.Aggregate.TotalRPS = float64(totalOK) / elapsed.Seconds()
	snap.Aggregate.TotalErrors = int(totalErr)
	if totalOK+totalErr > 0 {
		snap.Aggregate.ErrorPct = float64(totalErr) / float64(totalOK+totalErr) * 100
	}

	// Print results
	fmt.Println()
	fmt.Printf("═══ REALISTIC DAU BLAST RESULTS (%d goroutines × %d conns, %d streams) ═══\n",
		concurrency, numConns, numStreams)
	fmt.Printf("Actual: %.0f req/s | Errors: %d (%.1f%%) | Stream msgs: %d\n",
		snap.Aggregate.TotalRPS, totalErr, snap.Aggregate.ErrorPct, streamRecv.Load())
	if snap.RateLimited {
		fmt.Println("!! RATE LIMITED")
	}
	fmt.Println()

	fmt.Println(
		"╔════════════════════════╦══════════╦══════════╦══════════╦══════════╦══════════╦════════╗",
	)
	fmt.Println(
		"║ API                    ║ Actual   ║ Avg(ms)  ║ P50(ms)  ║ P99(ms)  ║ StdDev  ║ Err%   ║",
	)
	fmt.Println(
		"╠════════════════════════╬══════════╬══════════╬══════════╬══════════╬══════════╬════════╣",
	)
	for _, r := range snap.APIs {
		fmt.Printf("║ %-22s ║ %8.1f ║ %8.2f ║ %8.2f ║ %8.2f ║ %7.2f ║ %5.1f%% ║\n",
			r.Name, r.RPS, r.AvgLatency, r.P50Latency, r.P99Latency, r.StdDev, r.ErrorPct)
	}
	fmt.Println(
		"╚════════════════════════╩══════════╩══════════╩══════════╩══════════╩══════════╩════════╝",
	)

	// DAU calculation
	maxRPS := snap.Aggregate.TotalRPS
	rpsPerUser := cbCallsPerUserPerDay / cbSecondsPerDay
	maxDAU := maxRPS / rpsPerUser
	fmt.Printf("\n═══ DAU ESTIMATE ═══\n")
	fmt.Printf("Peak throughput: %.0f req/s (realistic workload)\n", maxRPS)
	fmt.Printf(
		"CB user profile: %.0f calls/day = %.4f req/s/user\n",
		cbCallsPerUserPerDay,
		rpsPerUser,
	)
	fmt.Printf("MAX SUPPORTED DAU: %.0f (~%.0fK)\n", maxDAU, maxDAU/1000)
	fmt.Printf("\nComparison: basic blast hit ~23K req/s (optimistic, cached, no blockchain)\n")
	fmt.Printf(
		"This test: diverse queries, real data, %d streams, 40%% of group msgs → blockchain\n",
		numStreams,
	)

	return snap, nil
}

func dauBlastLoop(
	ctx context.Context,
	conn *grpc.ClientConn,
	cfg *config,
	key *ecdsa.PrivateKey,
	apiName string,
	tracker *apiLatencyTracker,
	deadline time.Time,
	rateLimited *atomic.Bool,
	gwClient *gatewayClient,
	addresses []string,
	topics []*topic.Topic,
) {
	queryClient := messageApi.NewQueryApiClient(conn)
	publishClient := messageApi.NewPublishApiClient(conn)

	writeTopicIdx := uint64(0)
	var callCount uint64

	for {
		if time.Now().After(deadline) {
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
		}

		callCount++
		start := time.Now()
		var callErr error

		switch apiName {
		case "GetInboxIds":
			// Cycle through address pool — each goroutine hits different addresses
			addr := addresses[callCount%uint64(len(addresses))]
			_, callErr = queryClient.GetInboxIds(ctx, &messageApi.GetInboxIdsRequest{
				Requests: []*messageApi.GetInboxIdsRequest_Request{{
					Identifier:     addr,
					IdentifierKind: associations.IdentifierKind_IDENTIFIER_KIND_ETHEREUM,
				}},
			})

		case "QueryEnvelopes":
			// Query seeded topics by topic bytes (real production path)
			tp := topics[callCount%uint64(len(topics))]
			_, callErr = queryClient.QueryEnvelopes(ctx, &messageApi.QueryEnvelopesRequest{
				Query: &messageApi.EnvelopesQuery{
					Topics: [][]byte{tp.Bytes()},
				},
				Limit: 10,
			})

		case "GetNewestEnvelope":
			// Query seeded topics
			tp := topics[callCount%uint64(len(topics))]
			_, callErr = queryClient.GetNewestEnvelope(ctx, &messageApi.GetNewestEnvelopeRequest{
				Topics: [][]byte{tp.Bytes()},
			})

		case "Writes-Mixed":
			writeTopic := topics[writeTopicIdx%uint64(len(topics))]

			if gwClient != nil {
				// Write mix (per 20 calls):
				//   0-9:   group message application (48% → node)
				//   10-15: group message commit (32% → blockchain)
				//   16-18: welcome message (15% → node)
				//   19:    key package upload (5% → node)
				slot := callCount % 20
				switch {
				case slot < 10:
					// 48% application group messages — vary payload sizes
					slot60 := callCount % 60
					size := 1024 // 15%
					if slot60 < 36 {
						size = 256 // 60%
					} else if slot60 < 51 {
						size = 512 // 25%
					}
					callErr = gwClient.publishClientEnvelope(
						ctx, topic.TopicKindGroupMessagesV1, writeTopic, size,
					)
				case slot < 16:
					// 32% commit group messages → blockchain via GroupMessageBroadcaster
					size := 256
					if callCount%3 == 0 {
						size = 512
					}
					callErr = gwClient.publishCommitEnvelope(ctx, writeTopic, size)
				case slot < 19:
					// 15% welcome messages
					welcTopic := topic.NewTopic(
						topic.TopicKindWelcomeMessagesV1, randomBytes(32),
					)
					callErr = gwClient.publishClientEnvelope(
						ctx, topic.TopicKindWelcomeMessagesV1, welcTopic, 256,
					)
				default:
					// 5% key package rotations
					callErr = gwClient.publishKeyPackageEnvelope(ctx)
				}
			} else {
				// Direct to node — fallback (no commit/blockchain path)
				req, buildErr := buildPublishRequestForTopic(cfg, key, writeTopic, 256)
				if buildErr != nil {
					tracker.recordErr("build: " + buildErr.Error())
					continue
				}
				_, callErr = publishClient.PublishPayerEnvelopes(ctx, req)
			}

			// Rotate write topic every 50 writes (user switches conversations)
			if callCount%50 == 0 {
				writeTopicIdx++
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
