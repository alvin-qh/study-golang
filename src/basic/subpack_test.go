package basic

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"basic/subpack"
)

func TestSubpack(t *testing.T)  {
	s := subpack.Name()
	assert.Equal(t, "subpack", s, `expect s == "subpack"`)
}

func ExampleImportAnonymous() {
	ImportAnonymous()
}