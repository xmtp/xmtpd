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
