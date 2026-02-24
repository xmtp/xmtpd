package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	dbpkg "github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/seeds"
)

func databaseCmd() *cobra.Command {
	var (
		dbConnectionString string
		namespace          string
		numEnvelopes       int
		numOriginators     int
		numTopics          int
		numPayers          int
		blobSize           int
	)

	cmd := &cobra.Command{
		Use:          "database",
		Short:        "Populate the local database with production-like test data",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			logger, err := cliLogger()
			if err != nil {
				return err
			}

			ctx := context.Background()

			dbConn, err := dbpkg.NewNamespacedDB(
				ctx,
				logger,
				dbConnectionString,
				namespace,
				30*time.Second,
				10*time.Second,
				nil,
			)
			if err != nil {
				return fmt.Errorf("connect to database: %w", err)
			}
			defer func() {
				_ = dbConn.Close()
			}()

			cfg := seeds.Config{
				NumEnvelopes:    numEnvelopes,
				NumOriginators:  numOriginators,
				NumTopics:       numTopics,
				NumPayers:       numPayers,
				BlobSize:        blobSize,
				LogInterval:     10_000,
				NumUsageMinutes: 40,
			}

			result, err := seeds.SeedEnvelopes(ctx, dbConn, cfg, logger)
			if err != nil {
				return fmt.Errorf("seed envelopes: %w", err)
			}

			if err := seeds.SeedUsage(ctx, dbConn, result, cfg, logger); err != nil {
				return fmt.Errorf("seed usage: %w", err)
			}

			return nil
		},
		Example: `
xmtpd-cli generate populate --namespace <namespace>
xmtpd-cli generate populate --namespace <namespace> --envelopes 1000000 --payers 20
`,
	}

	cmd.Flags().StringVar(
		&dbConnectionString,
		"db-connection-string",
		"postgres://postgres:xmtp@localhost:8765/postgres",
		"PostgreSQL connection string (points to the control database)",
	)

	cmd.Flags().StringVar(
		&namespace,
		"namespace",
		"",
		"database namespace — must match the running node's namespace (see XMTPD_DB_NAME_OVERRIDE in .env)",
	)
	_ = cmd.MarkFlagRequired("namespace")

	cmd.Flags().IntVar(&numEnvelopes, "envelopes", 1_000_000, "number of gateway envelopes to seed")
	cmd.Flags().IntVar(&numOriginators, "originators", 3, "number of originator node IDs")
	cmd.Flags().IntVar(&numTopics, "topics", 100, "number of distinct topics")
	cmd.Flags().IntVar(&numPayers, "payers", 5, "number of payer accounts")
	cmd.Flags().IntVar(&blobSize, "blob-size", 500, "size of each envelope blob in bytes")

	return cmd
}
