package stress

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type CastSendCommand struct {
	ContractAddress string
	Function        string
	FunctionArgs    []string
	Rpc             string
	PrivateKey      string
	Nonce           *int
}

func buildCastSendCommand(c *CastSendCommand) string {
	var cmd strings.Builder
	cmd.WriteString(fmt.Sprintf(
		"cast send '%s' '%s'",
		c.ContractAddress,
		c.Function,
	))

	for _, arg := range c.FunctionArgs {
		cmd.WriteString(fmt.Sprintf(" '%s'", arg))
	}

	cmd.WriteString(fmt.Sprintf(" --rpc-url %s", c.Rpc))
	cmd.WriteString(fmt.Sprintf(" --private-key %s", c.PrivateKey))

	if c.Nonce != nil {
		cmd.WriteString(fmt.Sprintf(" --nonce %d", *c.Nonce))
	}

	return cmd.String()
}

func (c *CastSendCommand) Run(ctx context.Context) error {
	cli := buildCastSendCommand(c)

	ctxwt, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctxwt, "bash", "-c", cli)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cast send failed: %s", errBuf.String())
	}

	return nil
}
