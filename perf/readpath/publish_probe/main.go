package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
)

func b64(n int) string {
	buf := make([]byte, n)
	_, _ = rand.Read(buf)
	return base64.StdEncoding.EncodeToString(buf)
}

func main() {
	addr := flag.String("addr", "grpc.testnet-staging.xmtp.network:443", "gRPC host:port")
	call := flag.String("call", "xmtp.xmtpv4.message_api.ReplicationApi.PublishPayerEnvelopes", "Fully-qualified method")
	c := flag.Int("c", 16, "Concurrent workers")
	conn := flag.Int("conn", 4, "Client connections")
	dur := flag.Duration("dur", 30*time.Second, "Run duration")
	to := flag.Duration("timeout", 3*time.Second, "Per-call timeout")
	insecure := flag.Bool("insecure", false, "Use plaintext (no TLS)")

	originator := flag.Int("originator", 100, "target_originator (node id)")
	retention := flag.Int("retention_days", 7, "message_retention_days")
	ueBytes := flag.Int("ue_bytes", 256, "size of unsigned_client_envelope bytes")
	sigBytes := flag.Int("sig_bytes", 65, "size of payer_signature bytes (65 is typical)")

	out := flag.String("out", "publish_probe_report.json", "Report path")
	flag.Parse()

	bodyObj := map[string]any{
		"payer_envelopes": []any{
			map[string]any{
				"unsigned_client_envelope": b64(*ueBytes),
				"payer_signature":          map[string]any{"bytes": b64(*sigBytes)},
				"target_originator":        *originator,
				"message_retention_days":   *retention,
			},
		},
	}
	bodyBytes, _ := json.Marshal(bodyObj)

	opts := []runner.Option{
		runner.WithConnections(uint(*conn)),
		runner.WithConcurrency(uint(*c)),
		runner.WithRunDuration(*dur),
		runner.WithTimeout(*to),
		runner.WithDataFromJSON(string(bodyBytes)),
		runner.WithCountErrors(true),
	}
	if *insecure {
		opts = append(opts, runner.WithInsecure(true))
	}

	rep, err := runner.Run(*call, *addr, opts...)
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
