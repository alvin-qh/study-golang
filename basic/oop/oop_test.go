package oop

import (
	"testing"

	errs "study/basic/oop/errors"
	ifs "study/basic/oop/interfaces"
	"study/basic/oop/long"
	"study/basic/oop/size"
	"study/basic/oop/size3d"

	"github.com/stretchr/testify/assert"
)

// 测试 `long.Long` 类型
//
// `long.Long` 类型相当于为 `int64` 类型起了个别名, 根据 Go 语言规范, 重新定义的类型即可为其定义相关的类型函数
// Go 语言不支持运算符重载, 所以定义了 `long.Long` 类型严格上和 `int64` 类型并不是同一类型 (虽然值是完全一致的),
// 赋值时需要类型转化
//
// Go 语言不允许为其它包中定义的类型定义方法, 但通过给类型定义别名后, 即可为该别名类型定义方法, 例如 `long.Long`
// 类型的 `Compare` 方法和 `String` 方法
func TestOOP_TypeAlias(t *testing.T) {
	l1 := long.Long(100)
	assert.Equal(t, "100", l1.String())

	l2 := long.Long(200)
	assert.Equal(t, "200", l2.String())

	assert.False(t, l1.Compare(l2) == 0)
	assert.False(t, l1.Compare(l2) > 0)
	assert.True(t, l1.Compare(l2) < 0)
}

// 测试 `Size` 结构体
//
// `Size` 结构体提供了 `Width`, `Height` 以及 `Area` 方法, 用于获取宽度, 高度和面积值
//
// `Size` 结构体通过 `Compare` 方法实现了 `Comparable` 接口; 通过 `String` 方法实现了 `ToString` 接口
func TestOOP_StructType(t *testing.T) {
	// 初始化 `*size.Size` 结构体指针变量
	s1 := size.New(10, 20)
	assert.IsType(t, &size.Size{}, s1) // 确认 s1 的类型为 `*size.Size`
	assert.Equal(t, 10.0, s1.Width())  // 获取宽度
	assert.Equal(t, 20.0, s1.Height()) // 获取高度
	assert.Equal(t, 200.0, s1.Area())  // 获取面积
	assert.Equal(t, "<Size width=10 height=20>", s1.String())

	// 将 `*size.Size` 类型转化为 `ToString` 接口类型
	si := ifs.ToString(s1)
	assert.Equal(t, "<Size width=10 height=20>", si.String())

	// 将 `ToString` 接口类型恢复为 `*size.Size` 类型
	s1, ok := si.(*size.Size)
	assert.True(t, ok)

	// 定义未初始化的 `size.Size` 结构体变量
	var s2 size.Size
	// 调用初始化方法
	s2.Init(100, 200)
	// 测试 `Comparable` 接口方法
	assert.False(t, s1.Compare(&s2) == 0)
	assert.False(t, s1.Compare(&s2) > 0)
	assert.True(t, s1.Compare(&s2) < 0)

	// 因为 `size.Size::ToString` 方法的签名为 `func (s *Size) String() string`,
	// 相当于实现 `ToString` 接口的是 `*size.Size` 类型, 所以 `size.Size` 类型变量
	// 不能直接转为 `ToString` 接口类型;
	// 但如果方法签名为 `func (s Size) String() string`, 则 `*size.Size` 和 `size.Size`
	// 类型就都可以转为 `ToString` 接口类型
	// si = ToString(s2)
}

// 测试 `Size` 结构体
//
// `Size` 结构体提供了 `Width`, `Height` 以及 `Area` 方法, 用于获取宽度, 高度和面积值
//
// `Size` 结构体通过 `Compare` 方法实现了 `Comparable` 接口; 通过 `String` 方法实现了 `ToString` 接口
func TestOOP_StructInherit(t *testing.T) {
	s1 := size3d.New(10, 20, 30)
	assert.Equal(t, 10.0, s1.Width())
	assert.Equal(t, 20.0, s1.Height())
	assert.Equal(t, 30.0, s1.Depth())
	assert.Equal(t, 200.0, s1.Area())         // 上面测试的方法均从 `size.Size` 类型继承
	assert.Equal(t, 2200.0, s1.SurfaceArea()) // 下面的方法为 `size3d.Size3D` 类型定义
	assert.Equal(t, 6000.0, s1.Volume())

	// `String` 方法为 `size3d.Size3D` 类型实现 `ToString` 接口的方法
	assert.Equal(t, "<Size3D width=10 height=20 depth=30>", s1.String())

	// 将 `*size.Size` 类型转化为 `ToString` 接口类型
	si := ifs.ToString(s1)
	assert.Equal(t, "<Size3D width=10 height=20 depth=30>", si.String())

	// 将 `ToString` 接口类型恢复为 `*size.Size` 类型
	s1, ok := si.(*size3d.Size3D)
	assert.True(t, ok)

	var s2 size3d.Size3D
	s2.Init(100, 200, 300)
	// 测试 `Comparable` 接口方法
	assert.False(t, s1.Compare(&s2) == 0)
	assert.False(t, s1.Compare(&s2) > 0)
	assert.True(t, s1.Compare(&s2) < 0)

	// 因为 `size3d.Size3D` 类型的 `ToString` 方法的签名为 `func (s *Size3D) String() string`,
	// 相当于实现 `ToString` 接口的是 `*size.Size` 类型, 所以 `size.Size` 类型变量
	// 不能直接转为 `ToString` 接口类型;
	// 但如果方法签名为 `func (s Size3D) String() string`, 则 `*size3d.Size3D` 和 `size3d.Size3D`
	// 类型就都可以转为 `ToString` 接口类型
	// si = ToString(s2)
}

// 测试 `Comparable` 接口
//
// 对于 `Comparable` 接口的实现类型, 可以作为 `Eq`, `Ne`, `Gt`, `Ge` 等一系列函数的参数,
func TestOOP_Interface(t *testing.T) {
	// 定义两个 `Comparable` 接口的实例, 这里为 `*size.Size` 类型实例
	s1, s2 := size.New(10, 20), size.New(11, 21)

	// 调用函数, 传入 `Comparable` 接口实例作为参数
	assert.True(t, ifs.Eq(s1, s1))
	assert.True(t, ifs.Ne(s1, s2))
	assert.True(t, ifs.Ne(s2, s1))
	assert.True(t, ifs.Gt(s2, s1))
	assert.True(t, ifs.Lt(s1, s2))
	assert.True(t, ifs.Ge(s2, s1))
	assert.True(t, ifs.Ge(s1, s1))
	assert.True(t, ifs.Le(s1, s2))
	assert.True(t, ifs.Le(s1, s1))

	// 测试比较不同类型对象时, 出现的 panic 异常
	defer func() {
		e := recover().(error)
		assert.ErrorIs(t, e, errs.ErrInvalidType)
	}()

	// 虽然 `s3` 变量也为 `Comparable` 接口类型, 但其实际类型为 `*size3d.Size3D` 类型,
	// 所以当其和 `s1` 变量 (`*size.Size` 类型) 比较时会引发 Panic
	s3 := size3d.New(10, 20, 30)
	ifs.Eq(s3, s1)
	assert.Fail(t, "Cannot run here")
}
