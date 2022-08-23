package util

type ContextGetter interface {
	Get(name string) interface{}
}

// EchoGetBool returns false if the boool doesnt exist, or the value is false.
func EchoGetBool(c ContextGetter, name string) bool {
	vi := c.Get(name)
	if vi == nil {
		return false
	}
	v, _ := vi.(bool)
	return v
}
