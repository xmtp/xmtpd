package vectorclock

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/utils"
)

type ReadFunc func(context.Context) (map[uint32]uint64, error)

type VectorClock struct {
	lock *sync.RWMutex

	cfg config

	log             *zap.Logger
	readVectorClock ReadFunc
	vc              map[uint32]uint64
}

func New(log *zap.Logger, readfn ReadFunc, opts ...ConfigOption) *VectorClock {
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	vc := &VectorClock{
		lock:            &sync.RWMutex{},
		cfg:             cfg,
		log:             log.Named(utils.VectorClockLoggerName),
		readVectorClock: readfn,
		vc:              make(map[uint32]uint64),
	}

	return vc
}

// Start will do the initial vector clock sync with the DB, and start a sync loop in the background.
func (v *VectorClock) Start(ctx context.Context) error {
	err := v.forceSyncWithDB(ctx)
	if err != nil {
		return fmt.Errorf("could not run initial sync with DB: %w", err)
	}

	go v.runSyncLoop(ctx)

	return nil
}

// Ideally we have enough confidence in the in-memory implementation this sync loop can run very rarely (or not at all).
// However, initially, it might give us confidence to run this periodically and check correctness.
func (v *VectorClock) runSyncLoop(ctx context.Context) {
	v.log.Info("starting vector clock sync loop",
		zap.Int("resolve-strategy", int(v.cfg.resolveStrategy)),
		zap.Duration("sync-timeout", v.cfg.syncTimeout),
	)

	ticker := time.NewTicker(v.cfg.syncTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			v.log.Info("stopping vector clock sync loop")
			return
		case <-ticker.C:
			err := v.runIntegrityCheck(ctx)
			if err != nil {
				v.log.Error("integrity check failed", zap.Error(err))
			}
		}
	}
}

func (v *VectorClock) runIntegrityCheck(ctx context.Context) error {
	v.lock.RLock()

	dbvc, err := v.ReadFromDB(ctx)
	if err != nil {
		v.lock.RUnlock()
		return fmt.Errorf("could not read vector clock: %w", err)
	}

	err = v.compareAgainst(dbvc)

	v.lock.RUnlock()

	if err == nil {
		v.log.Debug("vector clock integrity ok")
		return nil
	}

	v.log.Error("vector clock mismatch detected", zap.Error(err))

	switch v.cfg.resolveStrategy {
	case ResolveCrash:
		v.log.Fatal("vector clock mismatch detected, halting", zap.Error(err))
	case ResolveReconcile:
		// TODO: What to do in case THIS again fails.
		_ = v.forceSyncWithDB(ctx)
		v.log.Info("vector clock force synced with DB")
	}

	return nil
}

func (v *VectorClock) ForceSync(ctx context.Context) error {
	return v.forceSyncWithDB(ctx)
}

func (v *VectorClock) forceSyncWithDB(ctx context.Context) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	// Force refetch now with the lock in case DB updated since last check.
	ref, err := v.ReadFromDB(ctx)
	if err != nil {
		return fmt.Errorf("could not read vector clock from DB: %w", err)
	}

	// Clear the map first in case there's non-overlapping keys.
	clear(v.vc)
	for id, seqID := range ref {
		v.vc[id] = seqID
	}

	return nil
}

func (v *VectorClock) compareAgainst(ref map[uint32]uint64) error {
	v.lock.RLock()
	defer v.lock.RUnlock()

	var errs []error
	if len(v.vc) != len(ref) {
		errs = append(
			errs,
			fmt.Errorf(
				"vector clocks have different lengths (me: %v, reference: %v)",
				len(v.vc),
				len(ref),
			),
		)
	}

	for id, seqID := range ref {
		if v.vc[id] == ref[id] {
			continue
		}

		errs = append(
			errs,
			fmt.Errorf("vector clock mismatch (id: %v, have: %v, want: %v)", id, v.vc[id], seqID),
		)
	}

	return errors.Join(errs...)
}

func (v *VectorClock) ReadFromDB(ctx context.Context) (map[uint32]uint64, error) {
	return v.readVectorClock(ctx)
}

func (v *VectorClock) Save(nodeID uint32, seqID uint64) {
	v.lock.Lock()
	defer v.lock.Unlock()

	// This code is not as dumb as it could be, as it checks if the to-be-saved sequence ID is larger than the current value.
	// This is something that could be done by the outer code; however, we have multiple tests that insert out-of-order sequence ID
	// so it seems we need to handle this scenario.

	existing, ok := v.vc[nodeID]
	if !ok {
		v.vc[nodeID] = seqID
		return
	}

	// Ignore lower seqID value.
	if existing > seqID {
		return
	}

	v.vc[nodeID] = seqID
}

func (v *VectorClock) Get(nodeID uint32) uint64 {
	v.lock.RLock()
	defer v.lock.RUnlock()

	// NOTE: Since sequenceID of 0 is not a thing, we do not
	// have to bother with returning (uint64, bool).

	seqID := v.vc[nodeID]
	return seqID
}

func (v *VectorClock) Values() map[uint32]uint64 {
	v.lock.RLock()
	defer v.lock.RUnlock()

	out := make(map[uint32]uint64)
	maps.Copy(out, v.vc)

	return out
}
