package basic

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	assert.Equal(t, 0, A, "expect A == 0")
	Simple()
	assert.Equal(t, 100, A, "expect A == 100")
}

func TestArguments(t *testing.T) {
	c := Arguments(10, 20)
	assert.Equal(t, c, 30, "expect c == 30")
}

func TestReturnMore(t *testing.T) {
	x, y := ReturnMore(10, 20, "Hello")
	assert.Equal(t, "Hello:10", x, `expect x == "Hello:10"`)
	assert.Equal(t, "Hello:20", y, `expect x == "Hello:20"`)
}
