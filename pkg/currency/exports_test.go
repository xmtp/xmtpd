package currency

import "fmt"

// toDollarsTestOnly converts PicoDollars to a dollar amount (as a float)
func (p PicoDollar) ToDollarsTestOnly() float64 {
	return float64(p) / PicoDollarsPerDollar
}

func (p PicoDollar) String() string {
	return fmt.Sprintf("%.12f", p.ToDollarsTestOnly())
}
