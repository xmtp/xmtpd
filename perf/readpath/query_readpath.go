package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/bojand/ghz/runner"
)

func main() {
	addr := flag.String("addr", "localhost:5050", "gRPC address of the node")
	method := flag.String("method", "xmtp.xmtpv4.message_api.ReplicationApi/QueryEnvelopes", "Fully-qualified gRPC method")
	concurrency := flag.Int("c", 32, "Number of concurrent workers")
	connections := flag.Int("conns", 0, "Number of client connections (0 = same as -c)")
	duration := flag.String("dur", "30s", "Test duration (e.g. 30s, 2m)")
	callTimeout := flag.Duration("timeout", 5*time.Second, "Per-call timeout")
	nodeID := flag.Int("node", 100, "originator node id to filter on")
	limit := flag.Int("limit", 5, "Query limit")
	summaryPath := flag.String("summary", "", "Optional path to write JSON summary (ghz format)")
	flag.Parse()

	if *connections == 0 {
		*connections = *concurrency
	}

	body := fmt.Sprintf(`{"query":{"originator_node_ids":[%d]},"limit":%d}`, *nodeID, *limit)

	rep, err := runner.Run(
		*method,
		*addr,
		runner.WithInsecure(true),
		runner.WithConcurrency(*concurrency),
		runner.WithConnections(*connections),
		runner.WithDuration(*duration),
		runner.WithDataFromJSON(body),
		runner.WithCallTimeout(*callTimeout),
		runner.WithName("query-readpath"),
		runner.WithReflect(true),
		runner.WithStatusCodeAndError(true),
	)
	if err != nil {
		log.Fatalf("run error: %v", err)
	}

	if *summaryPath != "" {
		if err := rep.Export(*summaryPath); err != nil {
			log.Printf("warning: failed to write summary: %v", err)
		}
	}

	fmt.Printf("\n=== QueryEnvelopes @ %s ===\n", *addr)
	fmt.Printf("method: %s\n", *method)
	fmt.Printf("concurrency: %d  connections: %d  duration: %s\n", *concurrency, *connections, *duration)
	fmt.Printf("requests: %d  rps: %.2f\n", rep.Count, rep.Rps)
	fmt.Printf("fastest: %v  slowest: %v  avg: %v\n", rep.Fastest, rep.Slowest, rep.Average)

	var p50, p90, p99 time.Duration
	for _, b := range rep.LatencyDistribution {
		switch int(b.Percentile + 0.5) {
		case 50:
			p50 = b.Latency
		case 90:
			p90 = b.Latency
		case 99:
			p99 = b.Latency
		}
	}
	if p50 > 0 || p90 > 0 || p99 > 0 {
		fmt.Printf("p50: %v  p90: %v  p99: %v\n", p50, p90, p99)
	}

	if len(rep.ErrorDist) > 0 {
		fmt.Printf("errors:\n")
		for msg, n := range rep.ErrorDist {
			fmt.Printf("  %7d  %s\n", n, msg)
		}
	} else {
		fmt.Println("errors: none")
	}
	fmt.Println()
}
