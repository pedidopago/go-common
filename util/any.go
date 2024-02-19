package util

func As[T any](v any, defaultv T) T {
	if v == nil {
		return defaultv
	}
	x, ok := v.(T)

	if !ok {
		return defaultv
	}

	return x
}
