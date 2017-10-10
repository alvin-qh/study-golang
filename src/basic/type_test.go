package basic

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
)

func TestInteger(t *testing.T) {
	a := Integer(0)

	assert.IsType(t, a, int32(0), `expect a is int32`)
}

func TestLong(t *testing.T) {
	a := Long(10)

	assert.Equal(t, "Long", reflect.TypeOf(a).Name())
	assert.Equal(t, "Long", reflect.TypeOf(a).Name())
}


