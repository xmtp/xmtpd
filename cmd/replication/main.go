package main

import (
	"context"
	"fmt"
	"log"
	"sync"

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

var Commit string = "unknown"

var options config.ServerOptions

func main() {
	if _, err := flags.Parse(&options); err != nil {
		if err, ok := err.(*flags.Error); !ok || err.Type != flags.ErrHelp {
			fatal("Could not parse options: %s", err)
		}
		return
	}

	if options.Version {
		fmt.Printf("Version: %s\n", Commit)
		return
	}

	logger, _, err := utils.BuildLogger(options.Log)
	if err != nil {
		fatal("Could not build logger: %s", err)
	}
	logger = logger.Named("replication")

	logger.Info(fmt.Sprintf("Version: %s", Commit))
	if options.Tracing.Enable {
		logger.Info("starting tracer")
		tracing.Start(Commit, logger)
		defer func() {
			logger.Info("stopping tracer")
			tracing.Stop()
		}()
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	doneC := make(chan bool, 1)
	tracing.GoPanicWrap(ctx, &wg, "main", func(ctx context.Context) {
		db, err := db.NewDB(
			ctx,
			options.DB.WriterConnectionString,
			options.DB.WaitForDB,
			options.DB.ReadTimeout,
		)

		if err != nil {
			logger.Fatal("initializing database", zap.Error(err))
		}

		ethclient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
		if err != nil {
			logger.Fatal("initializing blockchain client", zap.Error(err))
		}

		chainRegistry, err := registry.NewSmartContractRegistry(
			ethclient,
			logger,
			options.Contracts,
		)
		if err != nil {
			logger.Fatal("initializing smart contract registry", zap.Error(err))
		}
		err = chainRegistry.Start(ctx)
		if err != nil {
			logger.Fatal("starting smart contract registry", zap.Error(err))
		}

		signer, err := blockchain.NewPrivateKeySigner(
			options.Payer.PrivateKey,
			options.Contracts.ChainID,
		)
		if err != nil {
			logger.Fatal("initializing signer", zap.Error(err))
		}

		blockchainPublisher, err := blockchain.NewBlockchainPublisher(
			logger,
			ethclient,
			signer,
			options.Contracts,
		)
		if err != nil {
			logger.Fatal("initializing message publisher", zap.Error(err))
		}

		s, err := server.NewReplicationServer(
			ctx,
			logger,
			options,
			chainRegistry,
			db,
			blockchainPublisher,
		)
		if err != nil {
			log.Fatal("initializing server", zap.Error(err))
		}

		s.WaitForShutdown()
		doneC <- true
	})
	<-doneC

	cancel()
	wg.Wait()
}

func fatal(msg string, args ...any) {
	log.Fatalf(msg, args...)
}
