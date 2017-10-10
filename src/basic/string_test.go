package basic

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCharAt(t *testing.T) {
	const s = "Hello, 世界"

	var c = CharAt(s, 0)
	assert.Equal(t, rune(72), c)
	assert.Equal(t, "H", string(c))

	c = CharAt(s, 7)
	assert.Equal(t, "世", string(c))
}
