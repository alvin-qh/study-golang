package builtin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	s := "Hello, 大家好"
	assert.Equal(t, rune(s[1]), 'e')
	assert.Equal(t, string(s[1]), "e")

	as := []rune(s)
	assert.Equal(t, as[1], int32(s[1]))
	assert.NotEqual(t, as[8], int32(s[8]))
	assert.Equal(t, as[8], '家')
}

func TestLen(t *testing.T) {
	s := "Hello, 大家好"
	assert.Equal(t, Len(s), 10)
}

func TestCharAt(t *testing.T) {
	s := "Hello, 大家好"
	assert.Equal(t, CharAt(s, 1), 'e')
	assert.Equal(t, CharAt(s, 8), '家')
}
