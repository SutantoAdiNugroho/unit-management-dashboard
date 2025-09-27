package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmptyString(t *testing.T) {
	t.Run("Positive Case", func(t *testing.T) {
		assert.True(t, IsEmptyString(""))
		assert.True(t, IsEmptyString(" "))
		assert.True(t, IsEmptyString("   "))
		assert.True(t, IsEmptyString("\t"))
		assert.True(t, IsEmptyString("\n"))
		assert.True(t, IsEmptyString(" \t\n "))
	})

	t.Run("Negative Case", func(t *testing.T) {
		assert.False(t, IsEmptyString("hello"))
		assert.False(t, IsEmptyString(" hello "))
		assert.False(t, IsEmptyString("h"))
		assert.False(t, IsEmptyString("123"))
		assert.False(t, IsEmptyString("."))
		assert.False(t, IsEmptyString("!@#$"))
	})
}
