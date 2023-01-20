package types

import "github.com/pedidopago/go-common/util"

// TimeInterval is our version of googleapis google.type.Interval.
// It includes additional boolean fields indicating if the parameters
// are inclusive or exclusive.
// It is designed so that the new fields default value are interpreted the same
// as the google.type.Interval.
type TimeInterval struct {
	StartTime          *Time
	StartTimeExclusive bool
	EndTime            *Time
	EndTimeInclusive   bool
}

func (x *TimeInterval) UnmarshalText(text []byte) error {
	if x.StartTime != nil && x.EndTime != nil {
		// already filled
		return nil
	}
	if len(text) == 0 {
		return nil
	}

	var timePart []byte
	switch text[0] {
	case '>':
		if x.StartTime != nil {
			// already filled
			return nil
		}
		if len(text) < 2 {
			return nil
		}
		switch text[1] {
		case '=':
			if len(text) < 3 {
				return nil
			}
			timePart = text[2:]
		default:
			timePart = text[1:]
			x.StartTimeExclusive = true
		}
		x.StartTime = new(Time)
		return x.StartTime.UnmarshalText(timePart)

	case '<':
		if x.EndTime != nil {
			// already filled
			return nil
		}
		if len(text) < 2 {
			return nil
		}
		switch text[1] {
		case '=':
			if len(text) < 3 {
				return nil
			}
			timePart = text[2:]
			x.EndTimeInclusive = true
		default:
			timePart = text[1:]
		}
		x.EndTime = new(Time)
		return x.EndTime.UnmarshalText(timePart)
	}

	x.StartTime = new(Time)
	err := x.StartTime.UnmarshalText(text)
	x.EndTime = util.NewFromValue(*x.StartTime)

	return err
}
