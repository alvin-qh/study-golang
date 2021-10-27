package builtin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeferFunc(t *testing.T) {
	s := DeferFunc(func(s string) string {
		return s + " World"
	})
	assert.Equal(t, s, "Hello World")
}

func TestPanicFunc(t *testing.T) {
	s := PanicFunc("Hello")
	assert.Equal(t, s, "Hello")

	defer func() { assert.Equal(t, recover().(string), "Empty") }()
	PanicFunc("")
	assert.Fail(t, "Cannot run here")
}
