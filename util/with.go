package util

func With[T any](v T, fn func(T) error) error {
	return fn(v)
}

func With2[T1, T2 any](v1 T1, v2 T2) func(fn func(T1, T2) error) error {
	return func(fn func(T1, T2) error) error {
		return fn(v1, v2)
	}
}

func Try[T any](v T, err error) func(fn func(T) error) error {
	if err != nil {
		return func(fn func(T) error) error {
			return err
		}
	}
	return func(fn func(T) error) error {
		return fn(v)
	}
}
