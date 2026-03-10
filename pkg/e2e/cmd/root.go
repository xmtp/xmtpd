// Package cmd provides the CLI for the E2E test framework.
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/e2e/runner"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "xmtpd-e2e",
	Short: "E2E testing CLI for the XMTP Network",
	Long:  "xmtpd-e2e runs end-to-end tests against the XMTP Network, including nodes, gateways, chains, and clients.",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run E2E tests",
	RunE:  runE2E,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available E2E tests",
	RunE:  listTests,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("could not execute: %s", err)
	}
}

func init() {
	viper.SetEnvPrefix("XMTPD_E2E")
	viper.AutomaticEnv()

	pf := rootCmd.PersistentFlags()
	pf.String("log-level", "info", "log level: debug, info, warn, error")
	pf.String("log-encoding", "console", "log encoding: console, json")
	pf.String("output-format", "human", "output format: human, json")

	for _, name := range []string{"log-level", "log-encoding", "output-format"} {
		_ = viper.BindPFlag(name, pf.Lookup(name))
	}

	rf := runCmd.Flags()
	rf.StringSlice("test", nil, "run only specified test(s) by name")
	rf.String("xmtpd-image", "ghcr.io/xmtp/xmtpd:latest", "docker image for xmtpd nodes")
	rf.String("gateway-image", "ghcr.io/xmtp/xmtpd-gateway:latest", "docker image for gateways")
	rf.String("chain-image", "ghcr.io/xmtp/contracts:latest", "docker image for chain")
	rf.String("cli-image", "ghcr.io/xmtp/xmtpd-cli:latest", "docker image for cli")

	for _, name := range []string{"test", "xmtpd-image", "gateway-image", "chain-image", "cli-image"} {
		_ = viper.BindPFlag(name, rf.Lookup(name))
	}

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(listCmd)

	rootCmd.SetOut(os.Stderr)
}

func runE2E(cmd *cobra.Command, _ []string) error {
	logger, err := buildLogger()
	if err != nil {
		return err
	}
	defer func() {
		_ = logger.Sync()
	}()

	cfg := runner.Config{
		XmtpdImage:   viper.GetString("xmtpd-image"),
		GatewayImage: viper.GetString("gateway-image"),
		TestFilter:   viper.GetStringSlice("test"),
		OutputFormat: viper.GetString("output-format"),
		ChainImage:   viper.GetString("chain-image"),
		CLIImage:     viper.GetString("cli-image"),
	}

	r := runner.New(logger, cfg)

	return r.Run(cmd.Context())
}

func listTests(_ *cobra.Command, _ []string) error {
	logger, err := buildLogger()
	if err != nil {
		return err
	}

	var (
		r      = runner.New(logger, runner.Config{})
		tests  = r.Tests()
		format = viper.GetString("output-format")
	)

	if format == "json" {
		return listTestsJSON(tests)
	}

	listTestsTable(tests)

	return nil
}

func listTestsJSON(tests []runner.TestInfo) error {
	data, err := json.MarshalIndent(tests, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func listTestsTable(tests []runner.TestInfo) {
	if len(tests) == 0 {
		fmt.Println("No tests available.")
		return
	}

	nameWidth := len("NAME")
	for _, t := range tests {
		if len(t.Name) > nameWidth {
			nameWidth = len(t.Name)
		}
	}
	nameWidth += 2

	header := fmt.Sprintf("%-*s  %s", nameWidth, "NAME", "DESCRIPTION")
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)+10))

	for _, t := range tests {
		fmt.Printf("%-*s  %s\n", nameWidth, t.Name, t.Description)
	}
}

func buildLogger() (*zap.Logger, error) {
	l, _, err := utils.BuildLogger(config.LogOptions{
		LogLevel:    viper.GetString("log-level"),
		LogEncoding: viper.GetString("log-encoding"),
	})
	if err != nil || l == nil {
		return nil, err
	}
	return l.Named("xmtpd.e2e"), nil
}
