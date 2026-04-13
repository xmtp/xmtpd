package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"fmt"
	"math"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	apiv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	messageApi "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

// streamStats collects per-stream metrics from a subscriber goroutine.
type streamStats struct {
	MessagesRecv uint64
	Latencies    []time.Duration // delivery latencies (publish → recv)
	Errors       []string
	FirstMsgAt   time.Duration // time from stream open to first message
}

// publishTracker records the most recent publish timestamp per topic.
type publishTracker struct {
	mu    sync.RWMutex
	times map[string]time.Time // topic key → publish time
}

func newPublishTracker() *publishTracker {
	return &publishTracker{times: make(map[string]time.Time)}
}

func (pt *publishTracker) record(topicKey string, t time.Time) {
	pt.mu.Lock()
	pt.times[topicKey] = t
	pt.mu.Unlock()
}

func (pt *publishTracker) lookup(topicKey string) (time.Time, bool) {
	pt.mu.RLock()
	t, ok := pt.times[topicKey]
	pt.mu.RUnlock()
	return t, ok
}

func newGRPCConn(cfg *config) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if cfg.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		creds := credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	return grpc.NewClient(cfg.Addr, opts...)
}

// buildPublishRequestForTopic builds a signed PublishPayerEnvelopes request
// targeting a specific topic.
func buildPublishRequestForTopic(
	cfg *config,
	key *ecdsa.PrivateKey,
	targetTopic *topic.Topic,
	payloadSize int,
) (*messageApi.PublishPayerEnvelopesRequest, error) {
	aad := &envelopesProto.AuthenticatedData{
		TargetTopic: targetTopic.Bytes(),
	}

	clientEnv := &envelopesProto.ClientEnvelope{
		Aad: aad,
		Payload: &envelopesProto.ClientEnvelope_GroupMessage{
			GroupMessage: &apiv1.GroupMessageInput{
				Version: &apiv1.GroupMessageInput_V1_{
					V1: &apiv1.GroupMessageInput_V1{
						Data: makeGroupMessagePayload(payloadSize),
					},
				},
			},
		},
	}

	clientEnvBytes, err := proto.Marshal(clientEnv)
	if err != nil {
		return nil, fmt.Errorf("marshal client envelope: %w", err)
	}

	sig, err := utils.SignClientEnvelope(cfg.NodeID, clientEnvBytes, key)
	if err != nil {
		return nil, fmt.Errorf("sign envelope: %w", err)
	}

	payerEnv := &envelopesProto.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvBytes,
		PayerSignature: &associations.RecoverableEcdsaSignature{
			Bytes: sig,
		},
		TargetOriginator:     cfg.NodeID,
		MessageRetentionDays: constants.DefaultStorageDurationDays,
	}

	return &messageApi.PublishPayerEnvelopesRequest{
		PayerEnvelopes: []*envelopesProto.PayerEnvelope{payerEnv},
	}, nil
}

// runStreamTest runs a streaming performance test.
// -c controls the number of concurrent subscriber streams.
// -conn controls the number of publisher goroutines.
// -pub-rate controls the aggregate publish rate in msg/s.
func runStreamTest(cfg *config, tc testCase) (*testResult, error) {
	fmt.Printf("\n========== %s (streaming) ==========\n", tc.Name)

	conn, err := newGRPCConn(cfg)
	if err != nil {
		return nil, fmt.Errorf("grpc connect: %w", err)
	}
	defer func() { _ = conn.Close() }()

	numStreams := max(cfg.Concurrency, 1)
	numPublishers := max(cfg.Connections, 1)

	// Generate one random topic per stream.
	topics := make([]*topic.Topic, numStreams)
	topicKeys := make([]string, numStreams)
	for i := range topics {
		topicID := randomBytes(16)
		topics[i] = topic.NewTopic(topic.TopicKindGroupMessagesV1, topicID)
		topicKeys[i] = string(topics[i].Bytes())
	}

	tracker := newPublishTracker()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Duration+10*time.Second)
	defer cancel()

	// Signal for coordinated start: subscribers open first, then publishers begin.
	var subscribersReady sync.WaitGroup
	subscribersReady.Add(numStreams)

	deadline := time.Now().Add(cfg.Duration)

	fmt.Printf("Streams: %d, Publishers: %d, Pub rate: %d msg/s, Duration: %s\n",
		numStreams, numPublishers, cfg.PubRate, cfg.Duration)

	// --- Subscriber goroutines ---
	var wg sync.WaitGroup
	allStats := make([]streamStats, numStreams)

	for i := range numStreams {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			stats := &allStats[idx]
			topicBytes := topics[idx].Bytes()
			topicKey := topicKeys[idx]

			runSubscribeTopicsStream(
				ctx, conn, topicBytes, topicKey, deadline,
				tracker, stats, &subscribersReady,
			)
		}(i)
	}

	// Wait for all subscribers to open their streams.
	subscribersReady.Wait()
	time.Sleep(500 * time.Millisecond) // let server-side stream registration settle

	// --- Publisher goroutines ---
	var totalPublished atomic.Uint64
	var publishErrors atomic.Uint64

	key, err := ethcrypto.GenerateKey()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("generate key: %w", err)
	}

	// Rate limiting: distribute total pub-rate across publisher goroutines.
	perPublisherRate := max(cfg.PubRate/numPublishers, 1)

	publishClient := messageApi.NewPublishApiClient(conn)

	for p := range numPublishers {
		wg.Add(1)
		go func(pubIdx int) {
			defer wg.Done()
			ticker := time.NewTicker(time.Second / time.Duration(perPublisherRate))
			defer ticker.Stop()

			topicIdx := pubIdx % numStreams
			for {
				select {
				case <-ticker.C:
					if time.Now().After(deadline) {
						return
					}
					tIdx := topicIdx % numStreams
					topicIdx++

					req, buildErr := buildPublishRequestForTopic(
						cfg, key, topics[tIdx], 256,
					)
					if buildErr != nil {
						publishErrors.Add(1)
						continue
					}

					pubTime := time.Now()
					_, pubErr := publishClient.PublishPayerEnvelopes(ctx, req)
					if pubErr != nil {
						publishErrors.Add(1)
						continue
					}
					tracker.record(topicKeys[tIdx], pubTime)
					totalPublished.Add(1)

				case <-ctx.Done():
					return
				}
			}
		}(p)
	}

	wg.Wait()

	return aggregateStreamStats(
		tc.Name, cfg.Duration, allStats,
		totalPublished.Load(), publishErrors.Load(),
	), nil
}

func runSubscribeTopicsStream(
	ctx context.Context,
	conn *grpc.ClientConn,
	topicBytes []byte,
	topicKey string,
	deadline time.Time,
	tracker *publishTracker,
	stats *streamStats,
	ready *sync.WaitGroup,
) {
	client := messageApi.NewQueryApiClient(conn)
	stream, err := client.SubscribeTopics(ctx, &messageApi.SubscribeTopicsRequest{
		Filters: []*messageApi.SubscribeTopicsRequest_TopicFilter{
			{Topic: topicBytes},
		},
	})
	if err != nil {
		stats.Errors = append(stats.Errors, err.Error())
		ready.Done()
		return
	}
	ready.Done()

	streamOpen := time.Now()
	firstMsg := false
	for !time.Now().After(deadline) {

		resp, recvErr := stream.Recv()
		recvTime := time.Now()
		if recvErr != nil {
			if ctx.Err() != nil || time.Now().After(deadline) {
				break
			}
			stats.Errors = append(stats.Errors, recvErr.Error())
			break
		}
		if envs := resp.GetEnvelopes(); envs != nil {
			stats.MessagesRecv += uint64(len(envs.GetEnvelopes()))
		}
		if !firstMsg {
			stats.FirstMsgAt = recvTime.Sub(streamOpen)
			firstMsg = true
		}
		if pubTime, ok := tracker.lookup(topicKey); ok {
			stats.Latencies = append(stats.Latencies, recvTime.Sub(pubTime))
		}
	}
}

func aggregateStreamStats(
	name string,
	duration time.Duration,
	stats []streamStats,
	published uint64,
	pubErrors uint64,
) *testResult {
	var totalRecv uint64
	var allLatencies []time.Duration
	errMap := make(map[string]int)

	for i := range stats {
		totalRecv += stats[i].MessagesRecv
		allLatencies = append(allLatencies, stats[i].Latencies...)
		for _, e := range stats[i].Errors {
			if len(e) > 120 {
				e = e[:120] + "..."
			}
			errMap[e]++
		}
	}

	streamErrCount := 0
	for _, c := range errMap {
		streamErrCount += c
	}

	result := &testResult{
		Name:     name,
		Count:    totalRecv,
		RPS:      float64(totalRecv) / duration.Seconds(),
		OKCount:  int(totalRecv),
		ErrCount: streamErrCount + int(pubErrors),
		Errors:   errMap,
	}

	if total := result.OKCount + result.ErrCount; total > 0 {
		result.ErrorPct = float64(result.ErrCount) / float64(total) * 100
	}

	if len(allLatencies) > 0 {
		slices.Sort(allLatencies)

		var sum float64
		for _, l := range allLatencies {
			sum += msFromDuration(l)
		}
		result.AvgLatency = sum / float64(len(allLatencies))
		result.P50Latency = msFromDuration(percentile(allLatencies, 50))
		result.P95Latency = msFromDuration(percentile(allLatencies, 95))
		result.P99Latency = msFromDuration(percentile(allLatencies, 99))

		var sumSquares float64
		for _, l := range allLatencies {
			diff := msFromDuration(l) - result.AvgLatency
			sumSquares += diff * diff
		}
		result.StdDev = math.Sqrt(sumSquares / float64(len(allLatencies)))
	}

	fmt.Printf(
		"  Published: %d | Received: %d | Pub errors: %d | Stream errors: %d\n",
		published, totalRecv, pubErrors, streamErrCount,
	)
	if len(allLatencies) > 0 {
		fmt.Printf(
			"  Delivery latency — Avg: %.1fms | P50: %.1fms | P95: %.1fms | P99: %.1fms\n",
			result.AvgLatency, result.P50Latency, result.P95Latency, result.P99Latency,
		)
	} else {
		fmt.Println("  No messages received on streams (0 latency samples)")
	}

	return result
}

func percentile(sorted []time.Duration, pct int) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := max(int(math.Ceil(float64(pct)/100.0*float64(len(sorted))))-1, 0)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
