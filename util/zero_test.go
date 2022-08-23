package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsZero(t *testing.T) {
	assert.True(t, Must(IsZero("", 0, nil)))
	assert.True(t, Must(IsZero([]string{})))
	assert.False(t, Must(IsZero([]string{""})))
}
