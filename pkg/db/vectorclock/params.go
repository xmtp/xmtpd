package vectorclock

import (
	"time"
)

const (
	defaultSyncTimeout     = 5 * time.Second
	defaultResolveStrategy = ResolveReconcile
)

type ResolveStrategy int

const (
	// If we detect that our vector clock is out of sync, what should we do?
	ResolveReconcile ResolveStrategy = 0 // Overwrite our vector clock with the DB version.
	ResolveCrash     ResolveStrategy = 1 // Return a fatal error.
)
