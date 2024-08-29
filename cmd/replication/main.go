package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/server"
	"github.com/xmtp/xmtpd/pkg/tracing"
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
	addEnvVars()

	log, _, err := buildLogger(options)
	if err != nil {
		fatal("Could not build logger: %s", err)
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
			log.Fatal("initializing database", zap.Error(err))
		}

		s, err := server.NewReplicationServer(
			ctx,
			log,
			options,
			registry.NewFixedNodeRegistry([]registry.Node{}),
			db,
		)
		if err != nil {
			log.Fatal("initializing server", zap.Error(err))
		}
		s.WaitForShutdown()
		doneC <- true
	})

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	select {
	case sig := <-sigC:
		log.Info("ending on signal", zap.String("signal", sig.String()))
	case <-doneC:
	}
	cancel()
	wg.Wait()
}

func addEnvVars() {
	if connStr, hasConnstr := os.LookupEnv("WRITER_DB_CONNECTION_STRING"); hasConnstr {
		options.DB.WriterConnectionString = connStr
	}

	if connStr, hasConnstr := os.LookupEnv("READER_DB_CONNECTION_STRING"); hasConnstr {
		options.DB.ReaderConnectionString = connStr
	}

	if privKey, hasPrivKey := os.LookupEnv("SIGNER_PRIVATE_KEY"); hasPrivKey {
		options.SignerPrivateKey = privKey
	}
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
	log, err := cfg.Build()
	if err != nil {
		return nil, nil, err
	}

	log = log.Named("replication")

	return log, &cfg, nil
}
