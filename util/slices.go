package util

func ToInterfaces[T any](s []T) []interface{} {
	if0 := make([]interface{}, len(s))
	for i, v := range s {
		if0[i] = v
	}
	return if0
}

func ToSliceIfNotZero[T comparable](v T) []T {
	var zv T
	if v == zv {
		return nil
	}
	return []T{v}
}

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
