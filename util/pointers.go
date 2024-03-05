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

// SafeDeref will dereference a generic pointer.
// It will fail at compile time if ptr is not a pointer.
// The difference with Deref is that it will return the zero value of T if ptr is nil.
func SafeDeref[T any](ptr *T) T {
	if ptr == nil {
		var zv T
		return zv
	}
	return *ptr
}

func Ref[T any](v T) *T {
	return &v
}
