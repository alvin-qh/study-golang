package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	assert.Equal(t, 100, Default(100, 0))
	assert.Equal(t, 0, Default(0, 0))
	assert.Equal(t, "OK", Default("", "OK"))
	assert.Equal(t, "", Default("", ""))
}
