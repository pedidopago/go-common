package util

func NewFromValue[T any](v T) *T {
	ptr := new(T)
	*ptr = v
	return ptr
}
