package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/server"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var Version string

var options config.ServerOptions

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

	if options.Version {
		fmt.Printf("Version: %s\n", Version)
		return
	}

	err = config.ValidateServerOptions(options)
	if err != nil {
		fatal("Could not validate options: %s", err)
	}

	logger, _, err := utils.BuildLogger(options.Log)
	if err != nil {
		fatal("Could not build logger: %s", err)
	}
	logger = logger.Named("replication")

	logger.Info(fmt.Sprintf("Version: %s", Version))

	version, err := semver.NewVersion(Version)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not parse semver version (%s): %s", Version, err))
	}

	if options.Tracing.Enable {
		logger.Info("starting tracer")
		tracing.Start(Version, logger)
		defer func() {
			logger.Info("stopping tracer")
			tracing.Stop()
		}()
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	doneC := make(chan bool, 1)
	tracing.GoPanicWrap(ctx, &wg, "main", func(ctx context.Context) {
		var dbInstance *sql.DB
		if options.Replication.Enable || options.Sync.Enable || options.Indexer.Enable ||
			options.Payer.Enable {
			namespace := options.DB.NameOverride
			if namespace == "" {
				namespace = utils.BuildNamespace(options)
			}
			dbInstance, err = db.NewNamespacedDB(
				ctx,
				logger,
				options.DB.WriterConnectionString,
				namespace,
				options.DB.WaitForDB,
				options.DB.ReadTimeout,
			)
			if err != nil {
				logger.Fatal("initializing database", zap.Error(err))
			}
		}

		ethclient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
		if err != nil {
			logger.Fatal("initializing blockchain client", zap.Error(err))
		}

		chainRegistry, err := registry.NewSmartContractRegistry(
			ctx,
			ethclient,
			logger,
			options.Contracts,
		)
		if err != nil {
			logger.Fatal("initializing smart contract registry", zap.Error(err))
		}
		err = chainRegistry.Start()
		if err != nil {
			logger.Fatal("starting smart contract registry", zap.Error(err))
		}

		s, err := server.NewReplicationServer(
			ctx,
			logger,
			options,
			chainRegistry,
			dbInstance,
			fmt.Sprintf("0.0.0.0:%d", options.API.Port),
			fmt.Sprintf("0.0.0.0:%d", options.API.HTTPPort),
			version,
		)
		if err != nil {
			logger.Fatal("initializing server", zap.Error(err))
		}

		s.WaitForShutdown(10 * time.Second)
		doneC <- true
	})
	<-doneC

	cancel()
	wg.Wait()
}

func fatal(msg string, args ...any) {
	log.Fatalf(msg, args...)
}
