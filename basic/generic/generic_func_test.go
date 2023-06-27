package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Integer int

func TestGenericIntFloatAdd(t *testing.T) {
	rInt := GenericIntFloatAdd[int](1, 2)
	assert.Equal(t, 3, rInt)

	rInteger := GenericIntFloatAdd[Integer](1, 2)
	assert.Equal(t, Integer(3), rInteger)

	rFloat := GenericIntFloatAdd[float64](1.1, 1.2)
	assert.Equal(t, float64(2.3), rFloat)
}

func TestGenericAdd(t *testing.T) {
	r := GenericAdd(1, 2)
	assert.Equal(t, 3, r)
}
