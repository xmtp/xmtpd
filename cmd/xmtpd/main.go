package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/node"
	"github.com/xmtp/xmtpd/pkg/store/bolt"
	memstore "github.com/xmtp/xmtpd/pkg/store/mem"
	postgresstore "github.com/xmtp/xmtpd/pkg/store/postgres"
	"github.com/xmtp/xmtpd/pkg/zap"
)

// GitCommit should be included in the binary via -ldflags=-X ${COMMIT}
var GitCommit string

var opts node.Options

func main() {
	// Initialize options.
	_, err := flags.Parse(&opts)
	if err != nil {
		if err, ok := err.(*flags.Error); !ok || err.Type != flags.ErrHelp {
			fatal("error parsing options: %s", err)
		}
		return
	}

	// Initialize logger.
	log, err := zap.NewLogger(&opts.Log)
	if err != nil {
		fatal("error building logger: %s", err)
	}
	log.Info("running", zap.String("git-commit", GitCommit))

	ctx := context.New(context.Background(), log)

	// Initialize datastore.
	var store node.NodeStore
	switch {
	case opts.Store.Type == "postgres":
		log.Info("using postgres store")
		db, err := postgresstore.NewDB(&opts.Store.Postgres)
		if err != nil {
			fatal("error initializing postgres: %s", err)
		}
		store, err = postgresstore.NewNodeStore(ctx, db)
		if err != nil {
			fatal("error initializing postgres: %s", err)
		}
	case opts.Store.Type == "bolt":
		log.Info("using bolt store")
		store, err = bolt.NewNodeStore(ctx, &opts.Store.Bolt)
		if err != nil {
			fatal("error opening bolt store: %s", err)
		}
	default:
		log.Info("using memory store")
		store = memstore.NewNodeStore(ctx)
	}

	// Initialize node.
	node, err := node.New(ctx, store, &opts)
	if err != nil {
		log.Fatal("error initializing node", zap.Error(err))
	}
	defer node.Close()

	// Wait for shutdown.
	sig := waitForEndSignal()
	log.Info("ending", zap.String("signal", sig.String()))
}

func waitForEndSignal() os.Signal {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	return <-sigC
}

func fatal(msg string, args ...any) {
	log.Fatalf(msg, args...)
}
