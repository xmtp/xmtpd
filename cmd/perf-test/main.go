package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/jhump/protoreflect/desc"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	apiv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	messageApi "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/protobuf/proto"
)

const MINIMAL_APPLICATION_PAYLOAD = "0001000210aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa000000000000000101000000"

type TestCase struct {
	Name        string
	TopicKind   topic.TopicKind
	PayloadSize int
}

var testCases = []TestCase{
	{Name: "Welcome", TopicKind: topic.TOPIC_KIND_WELCOME_MESSAGES_V1, PayloadSize: 256},
	{Name: "KeyPackage", TopicKind: topic.TOPIC_KIND_KEY_PACKAGES_V1, PayloadSize: 256},
	{Name: "GroupMessage-256B", TopicKind: topic.TOPIC_KIND_GROUP_MESSAGES_V1, PayloadSize: 256},
	{Name: "GroupMessage-512B", TopicKind: topic.TOPIC_KIND_GROUP_MESSAGES_V1, PayloadSize: 512},
	{Name: "GroupMessage-1KB", TopicKind: topic.TOPIC_KIND_GROUP_MESSAGES_V1, PayloadSize: 1024},
	{Name: "GroupMessage-5KB", TopicKind: topic.TOPIC_KIND_GROUP_MESSAGES_V1, PayloadSize: 5120},
}

type TestResult struct {
	Name       string         `json:"name"`
	Count      uint64         `json:"count"`
	RPS        float64        `json:"rps"`
	AvgLatency float64        `json:"avg_latency_ms"`
	P50Latency float64        `json:"p50_latency_ms"`
	P95Latency float64        `json:"p95_latency_ms"`
	P99Latency float64        `json:"p99_latency_ms"`
	StdDev     float64        `json:"stddev_ms"`
	ErrorPct   float64        `json:"error_pct"`
	OKCount    int            `json:"ok_count"`
	ErrCount   int            `json:"err_count"`
	Errors     map[string]int `json:"errors,omitempty"`
}

func hashPayerSignatureInput(originatorID uint32, unsignedClientEnvelope []byte) []byte {
	targetBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(targetBytes, originatorID)
	return ethcrypto.Keccak256(
		[]byte(constants.TARGET_ORIGINATOR_SEPARATION_LABEL),
		targetBytes,
		[]byte(constants.PAYER_DOMAIN_SEPARATION_LABEL),
		unsignedClientEnvelope,
	)
}

func signClientEnvelope(originatorID uint32, unsignedClientEnvelope []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	hash := hashPayerSignatureInput(originatorID, unsignedClientEnvelope)
	return ethcrypto.Sign(hash, key)
}

func makeGroupMessagePayload(size int) []byte {
	header, _ := hex.DecodeString(MINIMAL_APPLICATION_PAYLOAD)
	if size <= len(header) {
		return header
	}
	padding := make([]byte, size-len(header))
	rand.Read(padding)
	return append(header, padding...)
}

func buildClientEnvelope(tc TestCase) *envelopesProto.ClientEnvelope {
	topicID := make([]byte, 16)
	rand.Read(topicID)
	targetTopic := topic.NewTopic(tc.TopicKind, topicID)

	aad := &envelopesProto.AuthenticatedData{
		TargetTopic: targetTopic.Bytes(),
		DependsOn:   &envelopesProto.Cursor{},
	}

	switch tc.TopicKind {
	case topic.TOPIC_KIND_GROUP_MESSAGES_V1:
		return &envelopesProto.ClientEnvelope{
			Aad: aad,
			Payload: &envelopesProto.ClientEnvelope_GroupMessage{
				GroupMessage: &apiv1.GroupMessageInput{
					Version: &apiv1.GroupMessageInput_V1_{
						V1: &apiv1.GroupMessageInput_V1{
							Data: makeGroupMessagePayload(tc.PayloadSize),
						},
					},
				},
			},
		}

	case topic.TOPIC_KIND_WELCOME_MESSAGES_V1:
		data := make([]byte, tc.PayloadSize)
		rand.Read(data)
		return &envelopesProto.ClientEnvelope{
			Aad: aad,
			Payload: &envelopesProto.ClientEnvelope_WelcomeMessage{
				WelcomeMessage: &apiv1.WelcomeMessageInput{
					Version: &apiv1.WelcomeMessageInput_V1_{
						V1: &apiv1.WelcomeMessageInput_V1{
							Data: data,
						},
					},
				},
			},
		}

	case topic.TOPIC_KIND_KEY_PACKAGES_V1:
		data := make([]byte, tc.PayloadSize)
		rand.Read(data)
		return &envelopesProto.ClientEnvelope{
			Aad: aad,
			Payload: &envelopesProto.ClientEnvelope_UploadKeyPackage{
				UploadKeyPackage: &apiv1.UploadKeyPackageRequest{
					KeyPackage: &apiv1.KeyPackageUpload{
						KeyPackageTlsSerialized: data,
					},
				},
			},
		}
	}
	return nil
}

func buildRequest(nodeID uint32, tc TestCase, key *ecdsa.PrivateKey) ([]byte, error) {
	clientEnv := buildClientEnvelope(tc)
	clientEnvBytes, err := proto.Marshal(clientEnv)
	if err != nil {
		return nil, fmt.Errorf("marshal client envelope: %w", err)
	}

	sig, err := signClientEnvelope(nodeID, clientEnvBytes, key)
	if err != nil {
		return nil, fmt.Errorf("sign envelope: %w", err)
	}

	payerEnv := &envelopesProto.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvBytes,
		PayerSignature: &associations.RecoverableEcdsaSignature{
			Bytes: sig,
		},
		TargetOriginator:     nodeID,
		MessageRetentionDays: uint32(constants.DEFAULT_STORAGE_DURATION_DAYS),
	}

	req := &messageApi.PublishPayerEnvelopesRequest{
		PayerEnvelopes: []*envelopesProto.PayerEnvelope{payerEnv},
	}

	return proto.Marshal(req)
}

func runTest(addr string, nodeID uint32, tc TestCase, concurrency, connections int, duration time.Duration, insecure bool) (*TestResult, error) {
	key, err := ethcrypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	poolSize := 2000
	fmt.Printf("\n========== %s ==========\n", tc.Name)
	fmt.Printf("Generating %d pre-built envelopes...\n", poolSize)

	payloads := make([][]byte, poolSize)
	for i := range payloads {
		payloads[i], err = buildRequest(nodeID, tc, key)
		if err != nil {
			return nil, fmt.Errorf("build request %d: %w", i, err)
		}
	}
	fmt.Printf("Running for %s with %d workers, %d connections...\n", duration, concurrency, connections)

	idx := 0
	opts := []runner.Option{
		runner.WithConnections(uint(connections)),
		runner.WithConcurrency(uint(concurrency)),
		runner.WithRunDuration(duration),
		runner.WithTimeout(30 * time.Second),
		runner.WithCountErrors(true),
		runner.WithBinaryDataFunc(func(_ *desc.MethodDescriptor, _ *runner.CallData) []byte {
			i := idx % len(payloads)
			idx++
			return payloads[i]
		}),
	}
	if insecure {
		opts = append(opts, runner.WithInsecure(true))
	}

	rep, err := runner.Run(
		"xmtp.xmtpv4.message_api.ReplicationApi.PublishPayerEnvelopes",
		addr,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("ghz run: %w", err)
	}

	pr := printer.ReportPrinter{Out: os.Stdout, Report: rep}
	_ = pr.Print("summary")

	result := &TestResult{
		Name:       tc.Name,
		Count:      rep.Count,
		RPS:        rep.Rps,
		AvgLatency: float64(rep.Average) / float64(time.Millisecond),
		Errors:     make(map[string]int),
	}

	for _, ld := range rep.LatencyDistribution {
		switch ld.Percentage {
		case 50:
			result.P50Latency = float64(ld.Latency) / float64(time.Millisecond)
		case 95:
			result.P95Latency = float64(ld.Latency) / float64(time.Millisecond)
		case 99:
			result.P99Latency = float64(ld.Latency) / float64(time.Millisecond)
		}
	}

	// Calculate stddev from histogram data
	if rep.Count > 1 {
		mean := float64(rep.Average)
		var sumSquares float64
		for _, detail := range rep.Details {
			diff := float64(detail.Latency) - mean
			sumSquares += diff * diff
		}
		result.StdDev = math.Sqrt(sumSquares/float64(rep.Count)) / float64(time.Millisecond)
	}

	for code, count := range rep.StatusCodeDist {
		if code == "OK" {
			result.OKCount = count
		} else {
			result.ErrCount += count
		}
	}
	for errMsg, count := range rep.ErrorDist {
		// Truncate long error messages
		if len(errMsg) > 120 {
			errMsg = errMsg[:120] + "..."
		}
		result.Errors[errMsg] = count
	}
	if result.OKCount+result.ErrCount > 0 {
		result.ErrorPct = float64(result.ErrCount) / float64(result.OKCount+result.ErrCount) * 100
	}

	return result, nil
}

func main() {
	addr := flag.String("addr", "grpc.testnet-staging.xmtp.network:443", "gRPC host:port")
	nodeID := flag.Uint("node-id", 100, "Target originator node ID")
	concurrency := flag.Int("c", 8, "Concurrent workers")
	connections := flag.Int("conn", 4, "Client connections")
	dur := flag.Duration("dur", 10*time.Second, "Duration per test")
	insecure := flag.Bool("insecure", false, "Use plaintext (no TLS)")
	tests := flag.String("tests", "all", "Comma-separated: Welcome,KeyPackage,GroupMessage-256B,GroupMessage-512B,GroupMessage-1KB,GroupMessage-5KB or 'all'")
	outPath := flag.String("out", "perf_results.json", "JSON results output path")
	flag.Parse()

	selectedTests := testCases
	if *tests != "all" {
		names := strings.Split(*tests, ",")
		nameSet := make(map[string]bool)
		for _, n := range names {
			nameSet[strings.TrimSpace(n)] = true
		}
		selectedTests = nil
		for _, tc := range testCases {
			if nameSet[tc.Name] {
				selectedTests = append(selectedTests, tc)
			}
		}
	}

	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║     XMTP D14N Node API Performance Test         ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")
	fmt.Printf("Target:      %s (node %d)\n", *addr, *nodeID)
	fmt.Printf("Concurrency: %d workers, %d connections\n", *concurrency, *connections)
	fmt.Printf("Duration:    %s per test\n", *dur)
	fmt.Printf("Tests:       %d selected\n", len(selectedTests))

	var results []TestResult

	for _, tc := range selectedTests {
		result, err := runTest(*addr, uint32(*nodeID), tc, *concurrency, *connections, *dur, *insecure)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR [%s]: %v\n", tc.Name, err)
			continue
		}
		results = append(results, *result)
	}

	// Summary table
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                         (Staging) Node API Latency by Message Type                  ║")
	fmt.Println("╠══════════════════════╦══════════╦══════════╦══════════╦══════════╦══════════╦════════╣")
	fmt.Println("║ Test                 ║ Count    ║ RPS      ║ Avg(ms)  ║ Stdev    ║ P99(ms)  ║ Err%   ║")
	fmt.Println("╠══════════════════════╬══════════╬══════════╬══════════╬══════════╬══════════╬════════╣")
	for _, r := range results {
		fmt.Printf("║ %-20s ║ %8d ║ %8.1f ║ %8.2f ║ %8.2f ║ %8.2f ║ %5.1f%% ║\n",
			r.Name, r.Count, r.RPS, r.AvgLatency, r.StdDev, r.P99Latency, r.ErrorPct)
	}
	fmt.Println("╚══════════════════════╩══════════╩══════════╩══════════╩══════════╩══════════╩════════╝")

	b, _ := json.MarshalIndent(results, "", "  ")

	if err := os.WriteFile(*outPath, b, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "\nERROR: failed to save results to %s: %v\n", *outPath, err)
	} else {
		fmt.Printf("\nResults saved to %s\n", *outPath)
	}
}
