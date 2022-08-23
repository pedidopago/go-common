package util

import (
	"os"
	"strconv"
)

func BoolEnv(name string) bool {
	e := os.Getenv(name)
	if e == "" {
		return false
	}
	b, _ := strconv.ParseBool(e)
	return b
}

func StringEnv(name, defaultValue string) string {
	e := os.Getenv(name)
	if e == "" {
		return defaultValue
	}
	return e
}
