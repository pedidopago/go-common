package util

import "strings"

// Deprecated: use slice.ToInterfaces instead
func ToInterfaces[T any](s []T) []interface{} {
	if0 := make([]interface{}, len(s))
	for i, v := range s {
		if0[i] = v
	}
	return if0
}

// Deprecated: use slice.ToSliceIfNotZero instead
func ToSliceIfNotZero[T comparable](v T) []T {
	var zv T
	if v == zv {
		return nil
	}
	return []T{v}
}

// Deprecated: use slice.FirstElemOrErr instead
func FirstElemOrErr[T any](emptyErr error) func(v []T, err error) (T, error) {
	return func(v []T, err error) (T, error) {
		var zv T
		if err != nil {
			return zv, err
		}
		if len(v) == 0 {
			return zv, emptyErr
		}
		return v[0], nil
	}
}

// Deprecated: use slice.BytesToStringErr instead
func BytesToStringErr(v []byte, err error) (string, error) {
	if err != nil {
		return "", err
	}
	return string(v), nil
}

// Deprecated: use slice.ContainsAny instead
func ContainsAny(s []string, v string) bool {
	for _, vv := range s {
		if strings.Contains(vv, v) {
			return true
		}
	}
	return false
}

func JoinResults[T any](s []T, err error) func(other []T, err error) ([]T, error) {
	return func(other []T, err2 error) ([]T, error) {
		if err != nil && err2 != nil {
			return nil, err
		}
		return append(s, other...), nil
	}
}

func Deduplicate[T any](s []T, eqfn func(T, T) bool) []T {
	var res []T
	for _, v := range s {
		var found bool
		for _, vv := range res {
			if eqfn(v, vv) {
				found = true
				break
			}
		}
		if !found {
			res = append(res, v)
		}
	}
	return res
}
