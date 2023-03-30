package util

import (
	"fmt"
	"reflect"
)

// IsZero returns true if all types are zero value
func IsZero(v ...interface{}) (bool, error) {
	for _, vv := range v {
		t := reflect.TypeOf(vv)
		if t == nil {
			continue
		}
		if !t.Comparable() {
			if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
				return false, fmt.Errorf("type is not comparable: %v", t)
			}
			if t.Kind() == reflect.Array && t.Len() > 0 {
				return false, nil
			}
			if reflect.ValueOf(vv).Len() > 0 {
				return false, nil
			}
		} else {
			if vv != reflect.Zero(t).Interface() {
				return false, nil
			}
		}
	}
	return true, nil
}

// Must panics if err is not nil
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Maybe[T comparable](fallback T) func(v T, err error) T {
	return func(v T, err error) T {
		if err != nil {
			return fallback
		}
		return v
	}
}

func ValueOrZero[T comparable](v *T) T {
	if v == nil {
		return reflect.Zero(reflect.TypeOf(v)).Interface().(T)
	}
	return *v
}

func PtrOption[T any, V any](v *T, f func(T) V) V {
	if v == nil {
		var zv V
		return zv
	}
	return f(*v)
}
