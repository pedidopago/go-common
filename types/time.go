package types

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Time time.Time

func (x *Time) FromFormattedTime(s string) error {
	t, err := time.Parse(time.RFC3339Nano, s)
	*x = Time(t)
	return err
}

func UnixMilli(msec int64) Time {
	return Time(time.UnixMilli(msec))
}

func (x *Time) UnmarshalText(text []byte) error {
	var s = string(text)
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		*x = UnixMilli(i)
		return nil
	}
	err = x.FromFormattedTime(s)
	if err == nil {
		return nil
	}
	return fmt.Errorf("invalid format. expected integer or string in format 'RFC3339Nano'")
}

func (x *Time) UnmarshalJSON(data []byte) error {
	if data[0] != '"' {
		i, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return fmt.Errorf("expected a valid integer")
		}
		*x = UnixMilli(i)
		return nil
	}
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	return x.FromFormattedTime(s)
}

func (x Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.Time().Format(time.RFC3339Nano) + `"`), nil
}

func (x Time) Time() time.Time {
	return time.Time(x)
}

func (Time) SchemaTypeKind() reflect.Kind {
	return reflect.String
}

func (x *Time) FromProto(s *timestamppb.Timestamp) *Time {
	if s == nil {
		return nil
	}
	*x = Time(s.AsTime())
	return x
}

func (x *Time) ToProto() *timestamppb.Timestamp {
	if x == nil {
		return nil
	}
	return timestamppb.New(x.Time())
}
