package util

func StringPtr(v string) *string {
	v2 := v
	return &v2
}

// Deref will dereference a generic pointer.
// It will fail at compile time if ptr is not a pointer.
// This is useful for functions that require pointers to structs and not structs directly.
func Deref[T any](ptr *T) T {
	return *ptr
}
