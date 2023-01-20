package types

import "github.com/shopspring/decimal"

// Percentage represents a percentage value. e.g. 99.99% = 99.99
type Percentage float64

func (x *Percentage) FromString(s string) *Percentage {
	if s == "" {
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}
	return x.FromDecimal(&d)
}

func (x *Percentage) FromDecimal(d *decimal.Decimal) *Percentage {
	if d == nil {
		return nil
	}
	f, _ := d.Shift(2).Float64()
	*x = Percentage(f)
	return x
}
