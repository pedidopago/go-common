package types

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Timestamp int64

func (x *Timestamp) FromProto(s *timestamppb.Timestamp) *Timestamp {
	if s == nil {
		return nil
	}
	*x = Timestamp(s.AsTime().UnixMilli())
	return x
}

func (x *Timestamp) ToProto() *timestamppb.Timestamp {
	if x == nil {
		return nil
	}
	return timestamppb.New(time.UnixMilli(int64(*x)))
}
