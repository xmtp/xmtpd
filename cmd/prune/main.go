package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/vectorclock"
	"github.com/xmtp/xmtpd/pkg/prune"
	"github.com/xmtp/xmtpd/pkg/utils"
)

var Version string

var options config.PruneOptions

func main() {
	_, err := flags.Parse(&options)
	if err != nil {
		if err, ok := err.(*flags.Error); !ok || err.Type != flags.ErrHelp {
			fatal("could not parse options: %s", err)
		}
		return
	}

	if Version == "" {
		Version = os.Getenv("VERSION")
		if Version == "" {
			fatal("could not determine version")
		}
	}

	logger, _, err := utils.BuildLogger(options.Log)
	if err != nil {
		fatal("could not build logger: %s", err)
	}
	logger = logger.Named(utils.PrunerLoggerName)

	logger.Info(fmt.Sprintf("version: %s", Version))

	validator := config.NewOptionsValidator(logger)

	err = validator.ParseJSONConfig(&options.Contracts)
	if err != nil {
		fatal("could not parse JSON contracts config: %s", err)
	}

	err = validator.ValidatePruneOptions(options)
	if err != nil {
		fatal("could not validate options: %s", err)
	}

	ctx := context.Background()

	namespace := options.DB.NameOverride
	if namespace == "" {
		namespace = utils.BuildNamespace(
			options.Signer.PrivateKey,
			options.Contracts.SettlementChain.NodeRegistryAddress,
		)
	}

	dbInstance, err := db.ConnectToDB(ctx, logger,
		options.DB.WriterConnectionString,
		namespace,
		options.DB.WaitForDB,
		options.DB.ReadTimeout,
		nil,
	)
	if err != nil {
		fatal("could not connect to DB: %s", err)
	}

	vc := vectorclock.New(logger, db.GetVectorClockReader(dbInstance))
	db := db.NewDBHandler(dbInstance, vc)
	pruneExecutor := prune.NewPruneExecutor(ctx, logger, db, &options.PruneConfig)

	err = pruneExecutor.Run()
	if err != nil {
		fatal("could not execute prune: %s", err)
	}
}

func fatal(msg string, args ...any) {
	log.Fatalf(msg, args...)
}
