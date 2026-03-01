package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"sync/atomic"
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
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

const (
	// Minimal valid MLS PrivateMessage frame (non-commit).
	// Used as the prefix for GroupMessage payloads to pass server-side validation.
	minimalMLSFrameHex = "0001000210aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa000000000000000101000000"

	envelopePoolSize = 2000

	methodPublish           = "xmtp.xmtpv4.message_api.ReplicationApi.PublishPayerEnvelopes"
	methodQueryEnvelopes    = "xmtp.xmtpv4.message_api.ReplicationApi.QueryEnvelopes"
	methodGetInboxIds       = "xmtp.xmtpv4.message_api.ReplicationApi.GetInboxIds"
	methodGetNewestEnvelope = "xmtp.xmtpv4.message_api.ReplicationApi.GetNewestEnvelope"
)

// testCase defines a single performance test.
// Read-path tests set Method and JSONPayload directly.
// Write-path tests set Method, TopicKind, and PayloadSize to build signed envelopes.
type testCase struct {
	Name        string
	Method      string
	JSONPayload string // static JSON request body for read-path tests
	TopicKind   topic.TopicKind
	PayloadSize int
}

var testCases = []testCase{
	// Read path
	{
		Name:   "QueryEnvelopes",
		Method: methodQueryEnvelopes,
		JSONPayload: `{
			"query":{"originator_node_ids":[100]},
			"limit":5
		}`,
	},
	{
		Name:   "GetInboxIds",
		Method: methodGetInboxIds,
		JSONPayload: `{
			"requests":[{
				"identifier":"0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
				"identifier_kind":"IDENTIFIER_KIND_ETHEREUM"
			}]
		}`,
	},
	{
		Name:        "GetNewestEnvelope",
		Method:      methodGetNewestEnvelope,
		JSONPayload: `{"topics":["AAAAAAAAAAAAAAAAAAAAAA=="]}`,
	},
	// Write path
	{
		Name:        "Welcome",
		Method:      methodPublish,
		TopicKind:   topic.TopicKindWelcomeMessagesV1,
		PayloadSize: 256,
	},
	{
		Name:        "GroupMessage-256B",
		Method:      methodPublish,
		TopicKind:   topic.TopicKindGroupMessagesV1,
		PayloadSize: 256,
	},
	{
		Name:        "GroupMessage-512B",
		Method:      methodPublish,
		TopicKind:   topic.TopicKindGroupMessagesV1,
		PayloadSize: 512,
	},
	{
		Name:        "GroupMessage-1KB",
		Method:      methodPublish,
		TopicKind:   topic.TopicKindGroupMessagesV1,
		PayloadSize: 1024,
	},
	{
		Name:        "GroupMessage-5KB",
		Method:      methodPublish,
		TopicKind:   topic.TopicKindGroupMessagesV1,
		PayloadSize: 5120,
	},
}

type testResult struct {
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

type config struct {
	Addr        string
	NodeID      uint32
	Concurrency int
	Connections int
	Duration    time.Duration
	Insecure    bool
}

func makeGroupMessagePayload(size int) []byte {
	header, _ := hex.DecodeString(minimalMLSFrameHex)
	if size <= len(header) {
		return header
	}
	padding := make([]byte, size-len(header))
	_, _ = rand.Read(padding)
	return append(header, padding...)
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

func buildClientEnvelope(tc testCase) *envelopesProto.ClientEnvelope {
	topicID := randomBytes(16)
	targetTopic := topic.NewTopic(tc.TopicKind, topicID)

	aad := &envelopesProto.AuthenticatedData{
		TargetTopic: targetTopic.Bytes(),
	}

	switch tc.TopicKind {
	case topic.TopicKindGroupMessagesV1:
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

	case topic.TopicKindWelcomeMessagesV1:
		return &envelopesProto.ClientEnvelope{
			Aad: aad,
			Payload: &envelopesProto.ClientEnvelope_WelcomeMessage{
				WelcomeMessage: &apiv1.WelcomeMessageInput{
					Version: &apiv1.WelcomeMessageInput_V1_{
						V1: &apiv1.WelcomeMessageInput_V1{
							Data: randomBytes(tc.PayloadSize),
						},
					},
				},
			},
		}
	default:
		panic(fmt.Sprintf("unsupported topic kind: %v", tc.TopicKind))
	}
}

func buildPublishRequest(
	cfg *config,
	tc testCase,
	key *ecdsa.PrivateKey,
) ([]byte, error) {
	clientEnv := buildClientEnvelope(tc)
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

	req := &messageApi.PublishPayerEnvelopesRequest{
		PayerEnvelopes: []*envelopesProto.PayerEnvelope{payerEnv},
	}

	return proto.Marshal(req)
}

func msFromDuration(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

func runGhz(
	cfg *config,
	tc testCase,
	opts []runner.Option,
) (*testResult, error) {
	fmt.Printf(
		"Running for %s with %d workers, %d connections...\n",
		cfg.Duration, cfg.Concurrency, cfg.Connections,
	)

	if cfg.Insecure {
		opts = append(opts, runner.WithInsecure(true))
	}

	rep, err := runner.Run(tc.Method, cfg.Addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("ghz run: %w", err)
	}

	pr := printer.ReportPrinter{Out: os.Stdout, Report: rep}
	_ = pr.Print("summary")

	return collectResult(tc.Name, rep), nil
}

func runReadTest(cfg *config, tc testCase) (*testResult, error) {
	fmt.Printf("\n========== %s ==========\n", tc.Name)

	opts := []runner.Option{
		runner.WithConnections(uint(cfg.Connections)),
		runner.WithConcurrency(uint(cfg.Concurrency)),
		runner.WithRunDuration(cfg.Duration),
		runner.WithTimeout(5 * time.Second),
		runner.WithCountErrors(true),
		runner.WithDataFromJSON(tc.JSONPayload),
	}

	return runGhz(cfg, tc, opts)
}

func runWriteTest(cfg *config, tc testCase) (*testResult, error) {
	key, err := ethcrypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	fmt.Printf("\n========== %s ==========\n", tc.Name)
	fmt.Printf(
		"Generating %d pre-built envelopes...\n",
		envelopePoolSize,
	)

	payloads := make([][]byte, envelopePoolSize)
	for i := range payloads {
		payloads[i], err = buildPublishRequest(cfg, tc, key)
		if err != nil {
			return nil, fmt.Errorf("build request %d: %w", i, err)
		}
	}

	var idx atomic.Uint64
	opts := []runner.Option{
		runner.WithConnections(uint(cfg.Connections)),
		runner.WithConcurrency(uint(cfg.Concurrency)),
		runner.WithRunDuration(cfg.Duration),
		runner.WithTimeout(30 * time.Second),
		runner.WithCountErrors(true),
		runner.WithBinaryDataFunc(
			func(_ *desc.MethodDescriptor, _ *runner.CallData) []byte {
				i := idx.Add(1)
				return payloads[i%uint64(len(payloads))]
			},
		),
	}

	return runGhz(cfg, tc, opts)
}

func runTest(cfg *config, tc testCase) (*testResult, error) {
	if tc.JSONPayload != "" {
		return runReadTest(cfg, tc)
	}
	return runWriteTest(cfg, tc)
}

func collectResult(name string, rep *runner.Report) *testResult {
	result := &testResult{
		Name:       name,
		Count:      rep.Count,
		RPS:        rep.Rps,
		AvgLatency: msFromDuration(rep.Average),
		Errors:     make(map[string]int),
	}

	for _, ld := range rep.LatencyDistribution {
		switch ld.Percentage {
		case 50:
			result.P50Latency = msFromDuration(ld.Latency)
		case 95:
			result.P95Latency = msFromDuration(ld.Latency)
		case 99:
			result.P99Latency = msFromDuration(ld.Latency)
		}
	}

	if rep.Count > 1 {
		mean := float64(rep.Average)
		var sumSquares float64
		for _, detail := range rep.Details {
			diff := float64(detail.Latency) - mean
			sumSquares += diff * diff
		}
		result.StdDev = math.Sqrt(
			sumSquares/float64(rep.Count),
		) / float64(time.Millisecond)
	}

	for code, count := range rep.StatusCodeDist {
		if code == "OK" {
			result.OKCount = count
		} else {
			result.ErrCount += count
		}
	}
	for errMsg, count := range rep.ErrorDist {
		if len(errMsg) > 120 {
			errMsg = errMsg[:120] + "..."
		}
		result.Errors[errMsg] = count
	}
	if total := result.OKCount + result.ErrCount; total > 0 {
		result.ErrorPct = float64(result.ErrCount) / float64(total) * 100
	}

	return result
}

func parseFlags() (*config, []testCase, string) {
	cfg := &config{}
	flag.StringVar(
		&cfg.Addr, "addr",
		"grpc.testnet-staging.xmtp.network:443", "gRPC host:port",
	)
	nodeID := flag.Uint("node-id", 100, "Target originator node ID")
	flag.IntVar(&cfg.Concurrency, "c", 8, "Concurrent workers")
	flag.IntVar(&cfg.Connections, "conn", 4, "Client connections")
	flag.DurationVar(
		&cfg.Duration, "dur", 10*time.Second, "Duration per test",
	)
	flag.BoolVar(
		&cfg.Insecure, "insecure", false, "Use plaintext (no TLS)",
	)
	tests := flag.String(
		"tests",
		"all",
		"Comma-separated test names or 'all'",
	)
	outPath := flag.String(
		"out", "perf_results.json", "JSON results output path",
	)
	flag.Parse()

	cfg.NodeID = uint32(*nodeID)

	selected := testCases
	if *tests != "all" {
		nameSet := make(map[string]bool)
		for n := range strings.SplitSeq(*tests, ",") {
			nameSet[strings.TrimSpace(n)] = true
		}
		selected = nil
		for _, tc := range testCases {
			if nameSet[tc.Name] {
				selected = append(selected, tc)
			}
		}
	}

	return cfg, selected, *outPath
}

func printSummaryTable(results []testResult) {
	const (
		top    = "╔══════════════════════╦══════════╦══════════╦══════════╦══════════╦══════════╦════════╗"
		title  = "║                      Node API Latency by Message Type                              ║"
		sep    = "╠══════════════════════╬══════════╬══════════╬══════════╬══════════╬══════════╬════════╣"
		header = "║ Test                 ║ Count    ║ RPS      ║ Avg(ms)  ║ Stdev    ║ P99(ms)  ║ Err%   ║"
		bottom = "╚══════════════════════╩══════════╩══════════╩══════════╩══════════╩══════════╩════════╝"
	)

	fmt.Println()
	fmt.Println(top)
	fmt.Println(title)
	fmt.Println(sep)
	fmt.Println(header)
	fmt.Println(sep)
	for _, r := range results {
		fmt.Printf(
			"║ %-20s ║ %8d ║ %8.1f ║ %8.2f ║ %8.2f ║ %8.2f ║ %5.1f%% ║\n",
			r.Name, r.Count, r.RPS,
			r.AvgLatency, r.StdDev, r.P99Latency, r.ErrorPct,
		)
	}
	fmt.Println(bottom)
}

func main() {
	cfg, selectedTests, outPath := parseFlags()

	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║     XMTP D14N Node API Performance Test         ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")
	fmt.Printf("Target:      %s (node %d)\n", cfg.Addr, cfg.NodeID)
	fmt.Printf(
		"Concurrency: %d workers, %d connections\n",
		cfg.Concurrency, cfg.Connections,
	)
	fmt.Printf("Duration:    %s per test\n", cfg.Duration)
	fmt.Printf("Tests:       %d selected\n", len(selectedTests))

	var results []testResult
	for _, tc := range selectedTests {
		result, err := runTest(cfg, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR [%s]: %v\n", tc.Name, err)
			continue
		}
		results = append(results, *result)
	}

	printSummaryTable(results)

	b, _ := json.MarshalIndent(results, "", "  ")
	_ = os.WriteFile(outPath, b, 0o644)
	fmt.Printf("\nResults saved to %s\n", outPath)
}
