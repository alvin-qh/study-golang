package oop

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompare(t *testing.T) {
	a := Long(100)
	b := Long(200)

	assert.True(t, a.Compare(b) < 0)
	assert.True(t, b.Compare(a) > 0)

	c := Int(100)
	assert.True(t, Int(a).Compare(c) == 0)

	p1 := NewSize(10, 20)
	p2 := p1
	assert.True(t, p1.Compare(p2) == 0)

	p2.Height = 10
	assert.True(t, p1.Compare(p2) > 0)

	p3 := NewSize3D(10, 20, 0)
	assert.True(t, p3.Compare(p1) == 0)

	p3 = NewSize3D(10, 20, 30)
	assert.True(t, p3.Compare(p1) > 0)
}
