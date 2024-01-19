package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimLines(t *testing.T) {
	testString := "hello\nworld\nhow\nare\nyou\ndoing"

	assert.Equal(t, "hello\nworld\nhow\nare\nyou\ndoing", TrimLines(testString, 0))
	assert.Equal(t, "hello\nworld\nhow\nare\nyou", TrimLines(testString, -1))
	assert.Equal(t, "how\nare\nyou\ndoing", TrimLines(testString, 2))
}
