package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/cmd"
)

// GitCommit should be included in the binary via -ldflags=-X ${COMMIT}
var GitCommit string

type Root struct {
	Start       cmd.Start       `command:"start"`
	GenerateKey cmd.GenerateKey `command:"generate-key"`
	ShowID      cmd.ShowID      `command:"show-id"`
	Version     cmd.Version     `command:"version"`
}

func main() {
	root := Root{
		Start: cmd.Start{
			GitCommit: GitCommit,
		},
		Version: cmd.Version{
			GitCommit: GitCommit,
		},
	}
	parser := flags.NewParser(&root, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}
