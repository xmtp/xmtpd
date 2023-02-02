package node

import (
	"context"

	"github.com/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/api"
	"go.uber.org/zap"
)

type Node struct {
	log       *zap.Logger
	ctx       context.Context
	ctxCancel context.CancelFunc
	api       *api.Server
}

func New(ctx context.Context, log *zap.Logger, opts *Options) (*Node, error) {
	n := &Node{
		log: log,
	}
	n.ctx, n.ctxCancel = context.WithCancel(ctx)
	var err error

	// Initialize API server/gateway.
	n.api, err = api.New(n.ctx, log, &opts.API)
	if err != nil {
		return nil, errors.Wrap(err, "initializing api")
	}

	return n, nil
}

func (n *Node) Close() {
	if n.api != nil {
		n.api.Close()
	}

	if n.ctxCancel != nil {
		n.ctxCancel()
	}
}
