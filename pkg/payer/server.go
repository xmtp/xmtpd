package payer

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/config"
)

type payerServiceImpl struct {
	identityFn  IdentityFn
	authorizers []AuthorizePublishFn
	config      *config.PayerConfig
}

func (s *payerServiceImpl) Serve(ctx context.Context) error {
	return nil
}

func (s *payerServiceImpl) ServeUntilShutdown() error {
	return s.Serve(context.Background())
}
