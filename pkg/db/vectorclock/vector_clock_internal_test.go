package vectorclock

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestVectorClock_CompareAgainstReference(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		vc := newSimpleVectorClock(t)

		err := vc.compareAgainst(map[uint32]uint64{})
		require.NoError(t, err)
	})
	t.Run("equal", func(t *testing.T) {
		var (
			vc = newSimpleVectorClock(t)

			ref = map[uint32]uint64{
				100: 1000,
				200: 2000,
				300: 3000,
				400: 4000,
			}
		)

		for k, v := range ref {
			vc.Save(k, v)
		}

		err := vc.compareAgainst(ref)
		require.NoError(t, err)
	})
	t.Run("different value detected", func(t *testing.T) {
		var (
			vc = newSimpleVectorClock(t)

			ref = map[uint32]uint64{
				100: 1000,
				200: 2000,
				300: 3000,
				400: 4000,
			}
		)

		for k, v := range ref {
			vc.Save(k, v)
		}

		// Change a single value and make sure it's reported.
		vc.Save(100, 1001)

		err := vc.compareAgainst(ref)
		require.Error(t, err)
	})
	t.Run("missing value detected", func(t *testing.T) {
		var (
			vc = newSimpleVectorClock(t)

			ref = map[uint32]uint64{
				100: 1000,
				200: 2000,
				300: 3000,
				400: 4000,
			}
		)

		for k, v := range ref {
			// Skip saving one value
			if k == 400 {
				continue
			}
			vc.Save(k, v)
		}

		err := vc.compareAgainst(ref)
		require.Error(t, err)
	})
	t.Run("extra value detected", func(t *testing.T) {
		var (
			vc = newSimpleVectorClock(t)

			ref = map[uint32]uint64{
				100: 1000,
				200: 2000,
				300: 3000,
				400: 4000,
			}
		)

		for k, v := range ref {
			vc.Save(k, v)
		}

		vc.Save(500, 5000)

		err := vc.compareAgainst(ref)
		require.Error(t, err)
	})
}

func newSimpleVectorClock(t *testing.T) *VectorClock {
	nopReadFunc := func(ctx context.Context) (map[uint32]uint64, error) {
		return nil, errors.New("no data")
	}

	return New(zap.NewNop(), nopReadFunc)
}
