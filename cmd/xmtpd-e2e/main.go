package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/e2e"
	"github.com/xmtp/xmtpd/pkg/zap"
)

// GitCommit should be included in the binary via -ldflags=-X ${COMMIT}
var GitCommit string

type Root struct {
	Log zap.Options `group:"Log options" namespace:"log"`
	E2E e2e.Options `group:"E2E options"`
}

func main() {
	root := Root{
		E2E: e2e.Options{
			GitCommit: GitCommit,
		},
	}
	parser := flags.NewParser(&root, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				return
			}
		}
		fatal(err)
	}

	log, err := zap.NewLogger(&root.Log)
	if err != nil {
		fatal(err)
	}
	ctx := context.New(context.Background(), log)

	_, err = e2e.New(ctx, &root.E2E)
	if err != nil {
		log.Fatal("running e2e", zap.Error(err))
	}
}

func fatal(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
