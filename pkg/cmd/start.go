package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/node"
	"github.com/xmtp/xmtpd/pkg/store/bolt"
	memstore "github.com/xmtp/xmtpd/pkg/store/mem"
	postgresstore "github.com/xmtp/xmtpd/pkg/store/postgres"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Start struct {
	node.Options `group:"node options"`
	AdminOptions `group:"admin options" namespace:"admin"`

	GitCommit string
}

func (c *Start) Execute(args []string) error {
	// Initialize logger.
	log, err := zap.NewLogger(&c.Log)
	if err != nil {
		return errors.Wrap(err, "error building logger")
	}
	log.Info("running", zap.String("git-commit", c.GitCommit))

	ctx := context.New(context.Background(), log)

	// Initialize admin interface (metrics/debug/...).
	metrics := startAdmin(ctx, &c.AdminOptions)

	// Initialize datastore.
	var store node.NodeStore
	switch {
	case c.Store.Type == "postgres":
		log.Info("using postgres store")
		db, err := postgresstore.NewDB(&c.Store.Postgres)
		if err != nil {
			return errors.Wrap(err, "error initializing postgres")
		}
		store, err = postgresstore.NewNodeStore(ctx, db)
		if err != nil {
			return errors.Wrap(err, "error initializing postgres")
		}
	case c.Store.Type == "bolt":
		log.Info("using bolt store")
		db, err := bolt.NewDB(&c.Store.Bolt)
		if err != nil {
			return errors.Wrap(err, "error opening bolt db")
		}
		store, err = bolt.NewNodeStore(ctx, db, &c.Store.Bolt)
		if err != nil {
			return errors.Wrap(err, "error creating bolt store")
		}
	default:
		log.Info("using memory store")
		store = memstore.NewNodeStore(ctx)
	}

	// Initialize node.
	node, err := node.New(ctx, metrics, store, &c.Options)
	if err != nil {
		ctx.Close()
		return errors.Wrap(err, "error initializing node")
	}
	defer node.Close()

	// Wait for shutdown.
	sig := waitForEndSignal()
	log.Info("ending", zap.String("signal", sig.String()))

	return nil
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
