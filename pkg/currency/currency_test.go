package currency

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConversion(t *testing.T) {
	initial, err := FromDollars(1.25)
	require.NoError(t, err)
	require.Equal(t, PicoDollar(1250000000000), initial)

	converted := initial.toDollarsTestOnly()
	require.Equal(t, 1.25, converted)
}

func TestString(t *testing.T) {
	initial, err := FromDollars(1.25)
	require.NoError(t, err)
	require.Equal(t, "1.250000000000", initial.String())
}

func TestToMicroDollars(t *testing.T) {
	initial, err := FromDollars(1.25)
	require.NoError(t, err)
	require.EqualValues(t, int64(1250000), initial.ToMicroDollars())
}
