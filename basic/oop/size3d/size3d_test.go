package size3d

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 Size3D 结构体
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

// 测试继承性
// Size3D 结构体从 Size 结构体继承，可以使用 Size 对象的函数
func TestInherit(t *testing.T) {
	s := New(10, 20, 30)
	assert.Equal(t, 200.0, s.Area()) // Area() 函数从 Size 类型继承
}

// 测试 String 函数
func TestToString(t *testing.T) {
	s := New(10, 20, 30)
	assert.Equal(t, "<Size3D width=10 height=20 depth=30>", s.String())
}

// 测试表面积函数
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

// 测试体积函数
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

// 测试比较两个对象
func TestCompare(t *testing.T) {
	s1 := New(10, 20, 30)
	assert.True(t, s1.Compare(s1) == 0)

	s2 := New(10, 20, 31)
	assert.True(t, s1.Compare(s2) < 0)
	assert.True(t, s2.Compare(s1) > 0)

	s2 = s1
	assert.True(t, s1.Compare(s2) == 0)
}
