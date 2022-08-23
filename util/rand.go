package util

import "math/rand"

func RandomStringFn(length int) func() string {
	return func() string {
		return RandomString(length, nil)
	}
}

var (
	defaultCharset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func RandomString(length int, charset []rune) string {
	if charset == nil {
		charset = defaultCharset
	}
	b := make([]rune, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
