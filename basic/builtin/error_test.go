package builtin

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorFromReturnValue(t *testing.T) {
	user, err := New(RandomString(20), "Alvin")
	assert.NoError(t, err)
	assert.Equal(t, "Alvin", user.name)

	user, err = New("", "")
	assert.Error(t, err)
	assert.Equal(t, "invalid name", err.Error())
	assert.Nil(t, user)
}

func TestWithPanicError(t *testing.T) {
	user := SafeNew(RandomString(20), "Alvin")
	assert.Equal(t, "Alvin", user.name)

	defer func() {
		if err := recover().(error); err != nil {
			assert.Equal(t, err.Error(), "invalid name")
		}
	}()
	SafeNew("", "")
	assert.Fail(t, "cannot be run here")
}

func TestErrorIsOrAs(t *testing.T) {
	user := SafeNew("", "Alvin")
	err := user.IsValid()
	assert.True(t, errors.Is(err, ErrIdRequired))
	assert.False(t, errors.Is(err, ErrIdLength))

	user = SafeNew("12345", "Alvin")
	err = user.IsValid()
	assert.False(t, errors.Is(err, ErrIdRequired))
	assert.True(t, errors.Is(err, ErrIdLength))
}
