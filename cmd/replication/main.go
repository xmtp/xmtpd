package main

import (
	"context"
	"log"
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/server"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Commit string

var options config.ServerOptions

func main() {
	if _, err := flags.Parse(&options); err != nil {
		if err, ok := err.(*flags.Error); !ok || err.Type != flags.ErrHelp {
			fatal("Could not parse options: %s", err)
		}
		return
	}

	logger, _, err := buildLogger(options)
	if err != nil {
		fatal("Could not build logger: %s", err)
	}

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

		privateKey, err := utils.ParseEcdsaPrivateKey(options.SignerPrivateKey)
		if err != nil {
			log.Fatal("parsing private key", zap.Error(err))
		}

		s, err := server.NewReplicationServer(
			ctx,
			logger,
			options,
			// TODO:nm replace with real node registry
			registry.NewFixedNodeRegistry(
				[]registry.Node{
					{
						NodeID:        1,
						SigningKey:    &privateKey.PublicKey,
						IsHealthy:     true,
						HttpAddress:   "http://example.com",
						IsValidConfig: true,
					},
				},
			),
			db,
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

func buildLogger(options config.ServerOptions) (*zap.Logger, *zap.Config, error) {
	atom := zap.NewAtomicLevel()
	level := zapcore.InfoLevel
	err := level.Set(options.LogLevel)
	if err != nil {
		return nil, nil, err
	}
	atom.SetLevel(level)

	cfg := zap.Config{
		Encoding:         options.LogEncoding,
		Level:            atom,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			NameKey:      "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, nil, err
	}

	logger = logger.Named("replication")

	return logger, &cfg, nil
}
