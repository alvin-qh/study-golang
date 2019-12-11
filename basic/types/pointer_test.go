package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestPointer(t *testing.T) {
	i := 0x12345678 // var i int = 0x12345678
	pi := &i        // var pi *int = &i

	assert.Equal(t, *pi, 0x12345678)

	i = 0x56781234
	assert.Equal(t, *pi, 0x56781234)

	*pi = 0x12345678
	assert.Equal(t, i, 0x12345678)

	psi := (*int16)(unsafe.Pointer(&i)) // var psi *int16 = (*int16)(unsafe.Pointer(&i))
	assert.Equal(t, *psi, int16(0x5678))
}

func TestChangeSize(t *testing.T) {
	size := Size{Width: 100, Height: 200}

	CannotChangeSize(size, 10, 20)
	assert.Equal(t, size.Width, 100)
	assert.Equal(t, size.Height, 200)

	ChangeSize(&size, 10, 20)
	assert.Equal(t, size.Width, 10)
	assert.Equal(t, size.Height, 20)

	size.CannotChange(100, 200)
	assert.Equal(t, size.Width, 10)
	assert.Equal(t, size.Height, 20)

	size.Change(100, 200)
	assert.Equal(t, size.Width, 100)
	assert.Equal(t, size.Height, 200)
}

func TestCreatePoint(t *testing.T) {
	var pa *int = nil
	CreatePoint(&pa)

	assert.NotEqual(t, pa, nil)
	assert.Equal(t, *pa, 100)
}
