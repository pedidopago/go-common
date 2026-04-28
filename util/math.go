package util

type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func Abs[T signed](v T) T {
	if v < 0 {
		return -v
	}
	return v
}
