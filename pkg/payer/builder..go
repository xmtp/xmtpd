package payer

import (
	"errors"

	"github.com/xmtp/xmtpd/pkg/config"
)

var (
	ErrMissingConfig = errors.New("config must be provided and not nil")
	ErrUnauthorized  = errors.New("unauthorized")
)

type PayerServiceBuilder struct {
	identityFn  IdentityFn
	authorizers []AuthorizePublishFn
	config      *config.PayerConfig
}

func NewPayerServiceBuilder(config *config.PayerConfig) IPayerServiceBuilder {
	return &PayerServiceBuilder{
		config: config,
	}
}

func (b *PayerServiceBuilder) WithIdentityFn(identityFn IdentityFn) IPayerServiceBuilder {
	b.identityFn = identityFn
	return b
}

func (b *PayerServiceBuilder) WithAuthorizers(
	authorizers ...AuthorizePublishFn,
) IPayerServiceBuilder {
	b.authorizers = authorizers
	return b
}

func (b *PayerServiceBuilder) Build() (PayerService, error) {
	if b.config == nil {
		return nil, ErrMissingConfig
	}

	if b.identityFn == nil {
		b.identityFn = IPIdentityFn
	}

	// Create a new PayerService with the configured options
	service := &payerServiceImpl{
		identityFn:  b.identityFn,
		authorizers: b.authorizers,
		config:      b.config,
	}

	return service, nil
}
