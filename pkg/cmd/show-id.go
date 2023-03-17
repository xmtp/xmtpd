package cmd

import (
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/xmtp/xmtpd/pkg/node"
)

type ShowID struct {
	NodeKey string `long:"node-key" env:"XMTP_NODE_KEY" description:"P2P node identity private key in hex format" required:"true"`
}

func (c *ShowID) Execute(args []string) error {
	privKey, err := node.DecodePrivateKey(c.NodeKey)
	if err != nil {
		return err
	}

	host, err := libp2p.New(
		libp2p.Identity(privKey),
	)
	if err != nil {
		return err
	}
	defer host.Close()

	fmt.Println(host.ID().String())

	return nil
}
