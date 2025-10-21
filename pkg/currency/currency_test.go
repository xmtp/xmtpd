package currency_test

import (
	"math/big"
	"testing"

	"github.com/xmtp/xmtpd/pkg/currency"

	"github.com/stretchr/testify/require"
)

func TestConversion(t *testing.T) {
	initial, err := currency.FromDollars(1.25)
	require.NoError(t, err)
	require.Equal(t, currency.PicoDollar(1250000000000), initial)

	converted := initial.ToDollars()
	require.Equal(t, 1.25, converted)
}

func TestString(t *testing.T) {
	initial, err := currency.FromDollars(1.25)
	require.NoError(t, err)
	require.Equal(t, "1.250000000000", initial.String())
}

func TestToMicroDollars(t *testing.T) {
	initial, err := currency.FromDollars(1.25)
	require.NoError(t, err)
	require.EqualValues(t, int64(1250000), initial.ToMicroDollars())
}

// FromWei: formatting & trimming behavior.
func TestFromWei_Basic(t *testing.T) {
	// 1 ether, 18 decimals -> "1"
	require.Equal(
		t,
		"1",
		currency.FromWei(big.NewInt(1).Mul(big.NewInt(1e18), big.NewInt(1)), 18),
	)

	// 1.5 ether -> "1.5"
	v := big.NewInt(0).Add(
		new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
		new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(17), nil)),
	)
	require.Equal(t, "1.5", currency.FromWei(v, 18))

	// 1234567 with 6 decimals -> "1.234567"
	require.Equal(t, "1.234567", currency.FromWei(big.NewInt(1_234_567), 6))

	// 1000000 with 6 decimals -> "1"
	require.Equal(t, "1", currency.FromWei(big.NewInt(1_000_000), 6))

	// 1000010 with 6 decimals -> "1.00001" (trims trailing zero only)
	require.Equal(t, "1.00001", currency.FromWei(big.NewInt(1_000_010), 6))

	// 0 with any decimals -> "0"
	require.Equal(t, "0", currency.FromWei(big.NewInt(0), 18))
}
