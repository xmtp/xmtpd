package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/prune"
	"github.com/xmtp/xmtpd/pkg/utils"
)

var Version string

var options config.PruneOptions

func main() {
	_, err := flags.Parse(&options)
	if err != nil {
		if err, ok := err.(*flags.Error); !ok || err.Type != flags.ErrHelp {
			fatal("Could not parse options: %s", err)
		}
		return
	}

	if Version == "" {
		Version = os.Getenv("VERSION")
		if Version == "" {
			fatal("Could not determine version")
		}
	}

	err = config.ValidatePruneOptions(options)
	if err != nil {
		fatal("Could not validate options: %s", err)
	}

	logger, _, err := utils.BuildLogger(options.Log)
	if err != nil {
		fatal("Could not build logger: %s", err)
	}
	logger = logger.Named("prune")

	logger.Info(fmt.Sprintf("Version: %s", Version))

	ctx := context.Background()

	dbInstance, err := db.ConnectToDB(ctx, logger,
		options.DB.WriterConnectionString,
		options.DB.WaitForDB,
		options.DB.ReadTimeout,
		options.DB.NameOverride,
	)
	if err != nil {
		fatal("Could not connect to DB: %s", err)
	}

	pruneExecutor := prune.NewPruneExecutor(ctx, logger, dbInstance)

	err = pruneExecutor.Run()
	if err != nil {
		fatal("Could not execute prune: %s", err)
	}
}

func fatal(msg string, args ...any) {
	log.Fatalf(msg, args...)
}
