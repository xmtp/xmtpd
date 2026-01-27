package vectorclock

import (
	"time"
)

// NOTE: If needed, these could be added to global config to allow easier configuration.

var defaultConfig = config{
	resolveStrategy: defaultResolveStrategy,
	syncTimeout:     defaultSyncTimeout,
}

type config struct {
	resolveStrategy ResolveStrategy
	syncTimeout     time.Duration
}

type ConfigOption func(*config)

func WithResolveStrategy(resolve ResolveStrategy) ConfigOption {
	return func(cfg *config) {
		cfg.resolveStrategy = resolve
	}
}

func WithSyncTimeout(d time.Duration) ConfigOption {
	return func(cfg *config) {
		cfg.syncTimeout = d
	}
}
