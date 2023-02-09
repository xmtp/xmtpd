package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	messagev1 "github.com/xmtp/xmtpd/pkg/api/message/v1"
	membroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/mem"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	memsyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/mem"
	"github.com/xmtp/xmtpd/pkg/node"
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

	// Initialize store.
	store := memstore.New(log)
	defer store.Close()

	// Initialize broadcaster.
	bc := membroadcaster.New(log)
	defer bc.Close()

	// Initialize syncer.
	syncer := memsyncer.New(log, store)
	defer bc.Close()

	// Initialize messagev1 service.
	messagev1, err := messagev1.New(log, store, bc, syncer)
	if err != nil {
		log.Fatal("error initializing messagev1", zap.Error(err))
	}
	defer messagev1.Close()

	// Initialize node.
	ctx := context.Background()
	node, err := node.New(ctx, log, messagev1, &opts)
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
