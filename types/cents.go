package types

import "github.com/shopspring/decimal"

// Cents represents BRL cents as an integer. i.e. R$ 1,99 = 199
type Cents int32

func (x *Cents) FromString(s string) *Cents {
	if s == "" {
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}
	return x.FromDecimal(&d)
}

func (x *Cents) FromDecimal(d *decimal.Decimal) *Cents {
	if d == nil {
		return nil
	}
	*x = Cents(d.Shift(2).IntPart())
	return x
}

func NewCentsFromString(s string) *Cents {
	if s == "" {
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}
	return NewCentsFromDecimal(&d)
}

func NewCentsFromDecimal(d *decimal.Decimal) *Cents {
	if d == nil {
		return nil
	}
	cents := Cents(d.Shift(2).IntPart())
	return &cents
}
