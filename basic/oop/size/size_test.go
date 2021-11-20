package size

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 Size 结构体
func TestSizeStruct(t *testing.T) {
	s := New(10, 20)
	assert.Equal(t, 10.0, s.Width())
	assert.Equal(t, 20.0, s.Height())

	s.Init(100, 200)
	assert.Equal(t, 100.0, s.Width())
	assert.Equal(t, 200.0, s.Height())
}

// 测试 Size 结构体计算面积
func TestArea(t *testing.T) {
	s := New(10, 20)
	assert.Equal(t, 200.0, s.Area())

	s.Init(20, 30)
	assert.Equal(t, 600.0, s.Area())

	s.Init(20, 0)
	assert.Equal(t, 0.0, s.Area())
}

// 测试 String 函数
func TestToString(t *testing.T) {
	var s Size
	assert.Equal(t, "<Size width=0 height=0>", s.String())

	ps := New(10, 20)
	assert.Equal(t, "<Size width=10 height=20>", ps.String())
}

// 测试比较两个 Size 结构体对象
func TestCompare(t *testing.T) {
	s1 := New(10, 20)
	assert.True(t, s1.Compare(s1) == 0)

	s2 := New(10, 21)
	assert.True(t, s1.Compare(s2) < 0)
	assert.True(t, s2.Compare(s1) > 0)

	s2 = s1
	assert.True(t, s1.Compare(s2) == 0)
}
