// Package currency implements the type PicoDollar and MicroDollar, and its conversions.
package currency

import (
	"errors"
	"fmt"
	"math"
	"math/big"
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

// toDollarsTestOnly converts PicoDollars to a dollar amount (as a float)
func (p PicoDollar) toDollarsTestOnly() float64 {
	return float64(p) / PicoDollarsPerDollar
}

// ToMicroDollars converts PicoDollars to MicroDollars (1e6 units per dollar)
func (p PicoDollar) ToMicroDollars() MicroDollar {
	return MicroDollar(p / 1e6)
}

func (p PicoDollar) String() string {
	return fmt.Sprintf("%.12f", p.toDollarsTestOnly())
}

func (m MicroDollar) ToBigInt() *big.Int {
	return big.NewInt(int64(m))
}
