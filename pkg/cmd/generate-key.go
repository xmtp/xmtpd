package cmd

import (
	"fmt"
	"os"

	"github.com/xmtp/xmtpd/pkg/node"
)

type GenerateKey struct {
	OutPath   string `long:"out" short:"o" description:"File path to output the node key (default: stdout)"`
	Overwrite bool   `long:"overwrite" description:"Overwrite output file if it already exists (default: false)"`
}

func (c *GenerateKey) Execute(args []string) error {
	pk, err := node.GeneratePrivateKey()
	if err != nil {
		return err
	}

	hex, err := node.PrivateKeyToHex(pk)
	if err != nil {
		return err
	}

	if c.OutPath == "" {
		fmt.Println(hex)
		return nil
	}

	_, err = os.Stat(c.OutPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if !c.Overwrite {
		return fmt.Errorf("output file %q already exists", c.OutPath)
	}

	f, err := os.Create(c.OutPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(hex)
	return err
}
