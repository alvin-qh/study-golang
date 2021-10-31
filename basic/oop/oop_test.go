package oop

import (
	"basic/oop/long"
	"basic/oop/size"
	"basic/oop/size3d"
	"basic/oop/typedef"
	"testing"

	"github.com/stretchr/testify/assert"
)

func eq(left, right typedef.Comparable) bool {
	return left.Compare(right) == 0
}

func ne(left, right typedef.Comparable) bool {
	return left.Compare(right) != 0
}

func gt(left, right typedef.Comparable) bool {
	return left.Compare(right) > 0
}

func lt(left, right typedef.Comparable) bool {
	return left.Compare(right) < 0
}

func ge(left, right typedef.Comparable) bool {
	return left.Compare(right) >= 0
}

func le(left, right typedef.Comparable) bool {
	return left.Compare(right) <= 0
}

// 测试 long.Long 类型
// long.Long 类型相当于为 int64 类型起了个别名，根据 go 语言规范，重新定义的类型即可为其定义相关的类型函数
// go 语言不支持运算符重载，所以定义了 long.Long 类型严格上和 int64 类型并不是同一类型（虽然值是完全一致的），赋值时需要类型转化
func TestLongType(t *testing.T) {
	// 生成 long.Long 类型对象
	var l1 long.Long = long.Long(100)
	assert.Equal(t, 100, int(l1))         // 通过类型转化，可以将 long.Long 转为 int 类型
	assert.Equal(t, "100", l1.ToString()) // 

	// 同类型之间可以直接赋值，对当前对象进行 copy
	l2 := l1
	assert.Equal(t, l1, l2)

	l2 = long.Long(200)
	assert.True(t, eq(l1, l1))
	assert.True(t, ne(l1, l2))

	assert.True(t, gt(l2, l1))
	assert.True(t, lt(l1, l2))

	assert.True(t, ge(l2, l1))
	assert.True(t, ge(l1, l1))

	assert.True(t, le(l1, l2))
	assert.True(t, le(l2, l2))
}

func TestSizeClass(t *testing.T) {
	s1 := size.New(10, 20)
	assert.Equal(t, 200.0, s1.Area())
	assert.Equal(t, "<Size width=10 height=20>", s1.ToString())

	width, height := s1.Value()
	assert.Equal(t, 10.0, width)
	assert.Equal(t, 20.0, height)

	eq := s1.Compare(s1)
	assert.Zero(t, eq)

	s2 := size.New(20, 30)
	eq = s1.Compare(s2)
	assert.Greater(t, eq, 0)

	eq = s2.Compare(s1)
	assert.Less(t, eq, 0)
}

func TestInterface(t *testing.T) {
	s1, s2 := size.New(10, 20), size.New(11, 21)
	s11 := s1.Clone()

	assert.True(t, eq(s1, s1))

	assert.True(t, eq(s1, s11))
	assert.True(t, eq(s11, s1))

	assert.True(t, ne(s1, s2))
	assert.True(t, ne(s2, s1))

	assert.True(t, gt(s2, s1))
	assert.True(t, lt(s1, s2))

	assert.True(t, ge(s2, s1))
	assert.True(t, ge(s1, s1))
	assert.True(t, ge(s1, s11))

	assert.True(t, le(s1, s2))
	assert.True(t, le(s1, s1))
	assert.True(t, le(s1, s11))

	defer func() {
		err := recover().(error)
		assert.ErrorIs(t, err, typedef.ErrType)
	}()

	ss1 := size3d.New(10, 20, 30)
	eq(ss1, s1)
	assert.Fail(t, "Cannot run here")
}
