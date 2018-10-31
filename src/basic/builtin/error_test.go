package builtin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithSafeError(t *testing.T) {
	msg, err := WithSafeError("Alvin")
	assert.NoError(t, err)
	assert.Equal(t, msg, "Hello Alvin")

	msg, err = WithSafeError("")
	assert.Error(t, err)
	assert.Equal(t, msg, "")
}

func TestWithPanicError(t *testing.T) {
	msg := WithPanicError("Alvin")
	assert.Equal(t, msg, "Hello Alvin")

	catch := func() {
		if e, ok := recover().(error); e != nil && ok {
			assert.Equal(t, e.Error(), "invalid name")
		}
	}
	defer catch()

	msg = WithPanicError("")
	assert.Fail(t, "Cannot run here")
}
