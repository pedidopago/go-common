package util

import "strings"

func TrimLines(s string, n int) string {
	lines := strings.Split(s, "\n")

	reverse := n < 0
	n = Abs(n)

	if reverse {
		lines = lines[:len(lines)-n]
	} else {
		lines = lines[n:]
	}

	return strings.Join(lines, "\n")
}
