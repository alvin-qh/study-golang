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

func TestExchange(t *testing.T) {
	var x, y interface{} = 1, 2

	x, y = ExchangeByReturn(x, y)

	assert.Equal(t, 2, x, `expect x == 2`)
	assert.Equal(t, 1, y, `expect y == 1`)
}
