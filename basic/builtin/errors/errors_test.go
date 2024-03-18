package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorTypeError(t *testing.T) {
	assert.Equal(t, ErrType, errors.New("invalid type"))
}
