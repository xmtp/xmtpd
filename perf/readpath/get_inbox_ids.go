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
	call := flag.String("call", "xmtp.xmtpv4.message_api.ReplicationApi.GetInboxIds", "Fully-qualified method")
	c := flag.Int("c", 16, "Concurrent workers")
	conn := flag.Int("conn", 4, "Client connections")
	dur := flag.Duration("dur", 30*time.Second, "Run duration")
	to := flag.Duration("timeout", 2*time.Second, "Per-call timeout")

	id := flag.String("id", "0x70997970C51812dc3A010C7d01b50e0d17dc79C8", "identifier")
	kind := flag.String("kind", "IDENTIFIER_KIND_EVM_ADDRESS", "identifier_kind enum")
	out := flag.String("out", "get_inbox_ids_report.json", "Report path")
	flag.Parse()

	bodyObj := map[string]any{
		"requests": []any{
			map[string]any{
				"identifier":       *id,
				"identifier_kind":  *kind,
			},
		},
	}
	body, _ := json.Marshal(bodyObj)

	rep, err := runner.Run(
		*call, *addr,
		runner.WithInsecure(true),
		runner.WithConnections(uint(*conn)),
		runner.WithConcurrency(uint(*c)),
		runner.WithRunDuration(*dur),
		runner.WithTimeout(*to),
		runner.WithDataFromJSON(string(body)),
		runner.WithCountErrors(true),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "run error: %v\n", err)
		os.Exit(1)
	}

	pr := printer.ReportPrinter{Out: os.Stdout, Report: rep}
	_ = pr.Print("pretty")
	if b, err := rep.MarshalJSON(); err == nil {
		_ = os.WriteFile(*out, b, 0644)
		fmt.Println("\nSaved report to", *out)
	}
}
