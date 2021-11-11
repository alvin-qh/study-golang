package size3d

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSize3DStruct(t *testing.T) {
	s := New(10, 20, 30)
	assert.Equal(t, 10.0, s.Width())
	assert.Equal(t, 20.0, s.Height())
	assert.Equal(t, 30.0, s.Depth())

	s.Init(100, 200, 300)
	assert.Equal(t, 100.0, s.Width())
	assert.Equal(t, 200.0, s.Height())
	assert.Equal(t, 300.0, s.Depth())
}

func TestInherit(t *testing.T) {
	s := New(10, 20, 30)
	assert.Equal(t, 200.0, s.Area()) // Area() 函数从 Size 类型继承
}

func TestToString(t *testing.T) {
	s := New(10, 20, 30)
	assert.Equal(t, "<Size3D width=10 height=20 depth=30>", s.String())
}

func TestSurfaceArea(t *testing.T) {
	s := New(10, 20, 30)
	assert.Equal(t, 2200.0, s.SurfaceArea())

	s.Init(100, 200, 300)
	assert.Equal(t, 220000.0, s.SurfaceArea())

	s.Init(100, 200, 0)
	assert.Equal(t, 40000.0, s.SurfaceArea())

	s.Init(100, 0, 300)
	assert.Equal(t, 60000.0, s.SurfaceArea())
}

func TestVolume(t *testing.T) {
	s := New(10, 20, 30)
	assert.Equal(t, 6000.0, s.Volume())

	s.Init(100, 200, 300)
	assert.Equal(t, 6e+06, s.Volume())

	s.Init(100, 200, 0)
	assert.Equal(t, 0.0, s.Volume())

	s.Init(100, 0, 300)
	assert.Equal(t, 0.0, s.Volume())
}

func TestCompare(t *testing.T) {
	s1 := New(10, 20, 30)
	assert.True(t, s1.Compare(s1) == 0)

	s2 := New(10, 20, 31)
	assert.True(t, s1.Compare(s2) < 0)
	assert.True(t, s2.Compare(s1) > 0)

    s2 = s1
    assert.True(t, s1.Compare(s2) == 0)
}
