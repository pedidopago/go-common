package util

import (
	"strconv"
	"strings"
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

// StringDateToTime parses a generic date "2006-01-02" or "02/01/2006" into a time.Time
func StringDateToTime(v string) (*time.Time, error) {
	if v == "" {
		return nil, nil
	}

	if strings.Contains(v, "-") {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return nil, err
		}
		return &t, nil
	}
	t, err := time.Parse("02/01/2006", v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
