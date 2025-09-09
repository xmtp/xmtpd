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
	RPC             string
	PrivateKey      string
	Nonce           *uint64
	Async           bool
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

	cmd.WriteString(fmt.Sprintf(" --rpc-url %s", c.RPC))
	cmd.WriteString(fmt.Sprintf(" --private-key %s", c.PrivateKey))

	if c.Nonce != nil {
		cmd.WriteString(fmt.Sprintf(" --nonce %d", *c.Nonce))
	}

	if c.Async {
		cmd.WriteString(" --async")
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
