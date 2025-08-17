package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
)

func main() {
	addr := flag.String("addr", "localhost:5050", "gRPC host:port")
	call := flag.String("call",
		"xmtp.xmtpv4.message_api.ReplicationApi.QueryEnvelopes",
		"Fully-qualified method name (package.Service/Method or package.Service.Method)")

	concurrency := flag.Int("c", 8, "Concurrent workers")
	connections := flag.Int("conn", 1, "Concurrent client connections")
	runDur := flag.Duration("dur", 5*time.Second, "Total run duration (e.g. 5s, 1m)")
	callTO := flag.Duration("timeout", 2*time.Second, "Per-call timeout")

	body := flag.String("json", `{"query":{"originator_node_ids":[100]},"limit":5}`,
		"Request JSON to send")

	// output
	outPath := flag.String("out", "readpath_report.json", "Write JSON report to file")

	flag.Parse()

	var tmp interface{}
	if err := json.Unmarshal([]byte(*body), &tmp); err != nil {
		fmt.Fprintf(os.Stderr, "invalid -json payload: %v\n", err)
		os.Exit(2)
	}

	rep, err := runner.Run(
		*call,
		*addr,
		runner.WithInsecure(true),
		runner.WithConnections(uint(*connections)),
		runner.WithConcurrency(uint(*concurrency)),
		runner.WithRunDuration(*runDur),
		runner.WithTimeout(*callTO),
		runner.WithDataFromJSON(*body),
		runner.WithCountErrors(true),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "run error: %v\n", err)
		os.Exit(1)
	}

	pr := printer.ReportPrinter{Out: os.Stdout, Report: rep}
	_ = pr.Print("pretty") // or "summary", "json", "html", etc.

	if b, err := rep.MarshalJSON(); err == nil {
		_ = os.WriteFile(*outPath, b, 0644)
		fmt.Printf("\nSaved report to %s\n", *outPath)
	} else {
		fmt.Fprintf(os.Stderr, "marshal report: %v\n", err)
	}
}
