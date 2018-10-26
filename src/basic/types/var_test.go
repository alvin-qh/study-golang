package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	b := Add(100)
	assert.Equal(t, b, int64(10000))
	assert.Equal(t, B, b)

	b = Add(30)
	assert.Equal(t, b, int64(13000))
	assert.Equal(t, B, b)
}
