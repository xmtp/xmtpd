package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// GitCommit should be included in the binary via -ldflags=-X ${COMMIT}
var GitCommit string

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	log.Info("running", zap.String("git-commit", GitCommit))

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
