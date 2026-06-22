package payer

import (
	"time"
)

const (
	defaultPublishTimeout = 30 * time.Second
	defaultPublishRetries = 5
)

type Config struct {
	PublishTimeout time.Duration
	PublishRetries uint

	// No-blockchain experiment: route commits and identity updates to dedicated nodes
	NoBlockchain   bool
	CommitNodeID   uint32
	IdentityNodeID uint32
}

var defaultConfig = Config{
	PublishTimeout: defaultPublishTimeout,
	PublishRetries: defaultPublishRetries,
}

type Option func(*Config)

func WithPublishTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.PublishTimeout = d
	}
}

func WithPublishRetries(n uint) Option {
	return func(cfg *Config) {
		cfg.PublishRetries = n
	}
}

func WithNoBlockchain(commitNodeID, identityNodeID uint32) Option {
	return func(cfg *Config) {
		cfg.NoBlockchain = true
		cfg.CommitNodeID = commitNodeID
		cfg.IdentityNodeID = identityNodeID
	}
}
