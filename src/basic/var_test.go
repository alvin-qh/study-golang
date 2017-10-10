package basic

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestVar(t *testing.T) {
	assert.Equal(t, 10, NumInt, "expect NumInt == 10")
	assert.IsType(t, int(0), NumInt, "expect NumInt is int")

	assert.IsType(t, int64(0), NumLong, "expect NumInt is int64")

	assert.Equal(t, true, Bool, "expect Bool is true")
}
