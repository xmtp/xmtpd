package main

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/debug"
	"github.com/xmtp/xmtpd/pkg/fees"
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

	if options.Version {
		fmt.Printf("version: %s\n", Version)
		return
	}

	logger, _, err := utils.BuildLogger(options.Log)
	if err != nil {
		fatal("could not build logger: %s", err)
	}

	logger = logger.Named(utils.BaseLoggerName)
	logger.Info(fmt.Sprintf("version: %s", Version))

	version, err := semver.NewVersion(Version)
	if err != nil {
		logger.Error(
			"could not parse semver version",
			zap.String("version", Version),
			zap.Error(err),
		)
	}

	// consolidate API options
	//nolint:staticcheck
	if options.Replication.Enable && !options.API.Enable {
		logger.Warn("--replication.enable is deprecated, use --api.enable instead")
		options.API.Enable = true
	}

	validator := config.NewOptionsValidator(logger)
	err = validator.ValidateServerOptions(&options)
	if err != nil {
		fatal("could not validate options: %s", err)
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
		promReg := prometheus.NewRegistry()

		var (
			dbh     *db.Handler
			readDB  *sql.DB
			writeDB *sql.DB
		)

		if options.Debug.Enable {
			pprofServer := debug.NewServer(options.Debug.Port)
			go func() {
				err := pprofServer.Start(ctx)
				if err != nil {
					logger.Fatal("starting pprof server", zap.Error(err))
				}
			}()

			defer func() {
				_ = pprofServer.Shutdown(ctx)
			}()
		}

		if options.API.Enable ||
			options.Sync.Enable ||
			options.Indexer.Enable ||
			options.MigrationServer.Enable ||
			options.PayerReport.Enable {

			namespace := cmp.Or(
				options.DB.NameOverride,
				utils.BuildNamespace(
					options.Signer.PrivateKey,
					options.Contracts.SettlementChain.NodeRegistryAddress,
				),
			)

			writeDB, err = db.NewNamespacedDB(
				ctx,
				logger,
				options.DB.WriterConnectionString,
				namespace,
				options.DB.WaitForDB,
				options.DB.ReadTimeout,
				promReg,
			)
			if err != nil {
				logger.Fatal("initializing writer database", zap.Error(err))
			}

			var dbopts []db.HandlerOption

			// If we have a separate reader DB initialize it here.
			if options.DB.ReaderConnectionString != "" {

				readDB, err = db.NewNamespacedDB(
					ctx,
					logger,
					options.DB.ReaderConnectionString,
					namespace,
					options.DB.WaitForDB,
					options.DB.ReadTimeout,
					promReg,
				)
				if err != nil {
					logger.Fatal("initializing reader database", zap.Error(err))
				}

				// Instruct db handler to include a read replica.
				dbopts = append(dbopts, db.WithReadReplica(readDB))
			}

			dbh = db.NewDBHandler(writeDB, dbopts...)
		}

		settlementChainClient, err := blockchain.NewRPCClient(
			ctx,
			options.Contracts.SettlementChain.RPCURL,
		)
		if err != nil {
			logger.Fatal("initializing blockchain client", zap.Error(err))
		}

		chainRegistry, err := registry.NewSmartContractRegistry(
			ctx,
			settlementChainClient,
			logger,
			&options.Contracts,
		)
		if err != nil {
			logger.Fatal("initializing smart contract registry", zap.Error(err))
		}
		err = chainRegistry.Start()
		if err != nil {
			logger.Fatal("starting smart contract registry", zap.Error(err))
		}

		feeCalculator, err := setupFeeCalculator(
			ctx,
			settlementChainClient,
			logger,
			&options.Contracts,
		)
		if err != nil {
			logger.Fatal("initializing fee calculator", zap.Error(err))
		}

		s, err := server.NewBaseServer(
			server.WithContext(ctx),
			server.WithLogger(logger),
			server.WithServerOptions(&options),
			server.WithNodeRegistry(chainRegistry),
			server.WithDB(dbh),
			server.WithFeeCalculator(feeCalculator),
			server.WithServerVersion(version),
			server.WithPromReg(promReg),
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

func setupFeeCalculator(
	ctx context.Context,
	ethclient bind.ContractCaller,
	logger *zap.Logger,
	contractsOptions *config.ContractsOptions,
) (*fees.FeeCalculator, error) {
	ratesFetcher, err := fees.NewContractRatesFetcher(ctx, ethclient, logger, contractsOptions)
	if err != nil {
		return nil, err
	}
	err = ratesFetcher.Start()
	if err != nil {
		return nil, err
	}

	return fees.NewFeeCalculator(ratesFetcher), nil
}

func fatal(msg string, args ...any) {
	log.Fatalf(msg, args...)
}
