// Package currency implements the type PicoDollar and MicroDollar, and its conversions.
package currency

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
)

// PicoDollar is a type to represent currency with 12 decimal precision
type (
	PicoDollar  int64
	MicroDollar int64
)

const (
	MicroDollarsPerDollar = 1e6
	PicoDollarsPerDollar  = 1e12
)

// FromDollars converts a dollar amount (as a float) to Picodollars
// This should mostly be used for testing, and real usage should be done purely in PicoDollars
func FromDollars(dollars float64) (PicoDollar, error) {
	if math.IsNaN(dollars) || math.IsInf(dollars, 0) {
		return 0, errors.New("invalid dollar amount: must be a finite number")
	}

	picodollars := dollars * PicoDollarsPerDollar
	if (picodollars < 0 && dollars > 0) || (picodollars > 0 && dollars < 0) {
		return 0, errors.New("overflow: dollar amount too large")
	}
	return PicoDollar(picodollars), nil
}

// FromMicrodollars converts an int64 microdollar amount to PicoDollar
func FromMicrodollars(microdollars MicroDollar) PicoDollar {
	return PicoDollar(microdollars * 1e6)
}

// ToMicroDollars converts PicoDollars to MicroDollars (1e6 units per dollar)
func (p PicoDollar) ToMicroDollars() MicroDollar {
	return MicroDollar(p / 1e6)
}

func (m MicroDollar) ToBigInt() *big.Int {
	return big.NewInt(int64(m))
}

// FromWei converts a wei value into a decimal string with the given decimals.
// For ETH, use decimals = 18.
// For an ERC20, use its `decimals()` value.
func FromWei(wei *big.Int, decimals int) string {
	if wei == nil {
		return "0"
	}
	// 10^decimals
	pow10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// Use big.Float with enough precision to avoid rounding issues.
	f := new(big.Float).SetPrec(256).SetInt(wei)
	div := new(big.Float).SetPrec(256).SetInt(pow10)
	val := new(big.Float).Quo(f, div)

	// Format with fixed decimals, then trim.
	s := fmt.Sprintf("%.*f", decimals, val)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".") // <-- fixes the "1." case
	if s == "" {                  // happens only if input was 0
		return "0"
	}
	return s
}
