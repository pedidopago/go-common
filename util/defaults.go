package util

import "golang.org/x/exp/constraints"

// Default will return the first non zero value in the given arguments.
func Default[T comparable](vs ...T) T {
	var zv T
	for _, v := range vs {
		if v != zv {
			return v
		}
	}
	return zv
}

// DefaultFn will return the first non zero value in the given arguments.
// If all arguments are zero, the fallback function will be called.
func DefaultFn[T comparable](fallback func() T, vs ...T) T {
	var zv T
	for _, v := range vs {
		if v != zv {
			return v
		}
	}
	return fallback()
}

func Max[T constraints.Ordered](vs ...T) T {
	var zv T
	for i, v := range vs {
		if v > zv || i == 0 {
			zv = v
		}
	}
	return zv
}

func Min[T constraints.Ordered](vs ...T) T {
	var zv T
	for i, v := range vs {
		if v < zv || i == 0 {
			zv = v
		}
	}
	return zv
}

func DidRun[T any](fn func() T, target *bool) func() T {
	return func() T {
		*target = true
		return fn()
	}
}

func DidRunFn[T any](fn func() T, targetfn func(v T)) func() T {
	return func() T {
		v := fn()
		defer targetfn(v)
		return v
	}
}
