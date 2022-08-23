package util

import (
	"strconv"
	"time"
)

func UnixTSStringToTime(v string) (time.Time, error) {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0), nil
}

func UnixTSStringToTimePtr(v string) *time.Time {
	t, err := UnixTSStringToTime(v)
	if err != nil {
		return nil
	}
	return &t
}
