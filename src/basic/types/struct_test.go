package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSize(t *testing.T) {
	s := Size{Width: 10, Height: 20}
	assert.Equal(t, s.Area(), 200)
}
