package oop

import (
	"basic/oop/long"
	"basic/oop/size"
	"basic/oop/size3d"
	"basic/oop/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	// 通过 typedef.Comparable 接口进行大小判断
	l2 = long.Long(200) // 赋值另一个 long.Long 类型对象

	cmp := l1.Compare(l1)   // 比较同一个对象是否相等
	assert.Equal(t, 0, cmp) // l1 == l1

	cmp = l2.Compare(l1)       // 比较 l2 和 l1 对象是否相等
	assert.NotEqual(t, 0, cmp) // l1 != l2

	cmp = l2.Compare(l1)      // 比较 l2 和 l1 对象大小
	assert.Greater(t, cmp, 0) // l2 > l1

	cmp = l1.Compare(l2)   // 比较 l1 和 l2 对象大小
	assert.Less(t, cmp, 0) // l1 < l2
}

// 测试 size.Size 类
func TestSizeStruct(t *testing.T) {
	// 构造 size.Size 类型的对象
	var s1 *size.Size = size.New(10, 20)
	assert.Equal(t, 200.0, s1.Area())
	assert.Equal(t, "<Size width=10 height=20>", s1.ToString())

	// 获取对象属性值
	width, height := s1.Width(), s1.Height()
	assert.Equal(t, 10.0, width)
	assert.Equal(t, 20.0, height)

	// 通过 typedef.Comparable 接口进行对象比较
	cmp := s1.Compare(s1) // 比较同一个对象
	assert.Zero(t, cmp)   // s1 == s1

	s2 := size.New(20, 30)
	cmp = s1.Compare(s2)   // 比较 s1 和 s2 大小
	assert.Less(t, cmp, 0) // s1 < s2

	cmp = s2.Compare(s1)      // 比较 s2 和 s1 的大小
	assert.Greater(t, cmp, 0) // s2 > s1
}

// 测试 size.Size3D 类
func TestSize3DStruct(t *testing.T) {
	var s1 *size3d.Size3D = size3d.New(10, 20, 30)
	assert.Equal(t, 2200.0, s1.Area())
	assert.Equal(t, "<Size3D width=10 height=20 depth=30>", s1.ToString())

	// 获取 size3d.Size3D 的属性值
	// 其中 Width(), Height() 两个函数是从 size.Size 类继承的，Depth() 函数是由 size3d.Size3D 提供
	width, height, depth := s1.Width(), s1.Height(), s1.Depth()
	assert.Equal(t, 10.0, width)
	assert.Equal(t, 20.0, height)
	assert.Equal(t, 30.0, depth)

	// 通过 typedef.Comparable 接口进行对象比较
	cmp := s1.Compare(s1) // 比较同一个对象
	assert.Zero(t, cmp)   // s1 == s1

	s2 := size3d.New(20, 30, 40)
	cmp = s1.Compare(s2)   // 比较 s1 和 s2 大小
	assert.Less(t, cmp, 0) // s1 < s2

	cmp = s2.Compare(s1)      // 比较 s2 和 s1 的大小
	assert.Greater(t, cmp, 0) // s2 > s1
}

// 定义一组通过 typedef.Comparable 接口对对象进行比较的函数
func eq(left, right types.Comparable) bool {
	return left.Compare(right) == 0
}

func ne(left, right types.Comparable) bool {
	return left.Compare(right) != 0
}

func gt(left, right types.Comparable) bool {
	return left.Compare(right) > 0
}

func lt(left, right types.Comparable) bool {
	return left.Compare(right) < 0
}

func ge(left, right types.Comparable) bool {
	return left.Compare(right) >= 0
}

func le(left, right types.Comparable) bool {
	return left.Compare(right) <= 0
}

// 测试 typedef.Comparable 接口进行对象比较
func TestInterface(t *testing.T) {
	// 比较两个对象的大小
	s1, s2 := size.New(10, 20), size.New(11, 21)
	s11 := s1

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

	// 测试比较不同类型对象时，出现的 panic 异常
	defer func() {
		err := recover().(error)
		assert.ErrorIs(t, err, types.ErrType)
	}()

	ss1 := size3d.New(10, 20, 30)
	eq(ss1, s1)
	assert.Fail(t, "Cannot run here")
}
