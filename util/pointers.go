package util

func StringPtr(v string) *string {
	v2 := v
	return &v2
}
