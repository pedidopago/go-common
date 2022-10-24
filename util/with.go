package util

func With[T any](v T, fn func(T) error) error {
	return fn(v)
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
