package size3d

import (
	"fmt"
	errs "study/basic/oop/errors"
	"study/basic/oop/size"
)

// 定义 Size3D 结构体, 从 Size 结构体继承
//
// 可以为类型定义方法, 格式为:
//
//	func (receiver_type) func_name([parameter_list]) [return_types] {...}
//
// 其中类型的方法会传入类型实例或类型实例的指针, 称为 Receiver
//
// 结构体也是一个类型, 所以可以为其定义方法, 对于结构体类型来说
//   - 类型 Size3D 包含全部 Receiver 为 Size3D 的方法
//   - 类型 *Size3D 包含全部 Receiver 为 Size3D + *Size3D 的方法
//   - Size3D 包含 Size 类型匿名字段, 则 Size3D 和 *Size3D 类型包含 Size 类型的方法
//   - Size3D 包含 *Size 类型匿名字段, 则 Size3D 和 *Size3D 类型包含 Size + *Size 类型的方法
//   - *Size3D 总会包含 Size + *Size 类型的方法
type Size3D struct {
	size.Size         // 继承 Size 结构体
	depth     float64 // 定义深度
}

// 构造函数, 生成 Size3D 对象
func New(width, height, depth float64) *Size3D { return new(Size3D).Init(width, height, depth) }

// 初始化 Size3D 结构体
func (s *Size3D) Init(width, height, depth float64) *Size3D {
	s.Size.Init(width, height)
	s.depth = depth
	return s
}

// 获取深度
func (s *Size3D) Depth() float64 { return s.depth }

// 求 `Size3D` 表面积
func (s *Size3D) SurfaceArea() float64 {
	width, height := s.Size.Width(), s.Size.Height()
	return (s.Size.Area() + width*s.depth + height*s.depth) * 2
}

// 求 `Size3D` 体积
func (s *Size3D) Volume() float64 { return s.Size.Area() * s.depth }

// 结构体转为字符串
//
// 如果省略了这个函数定义, 则 `Size3D` 类型会使用 `size.Size` 类型的 `String()` 方法
//
// 所以, 如果父类实现了某个接口, 则子类也会实现这个接口, 子类也可以重写接口的实现方法
func (s *Size3D) String() string {
	width, height, depth := s.Size.Width(), s.Size.Height(), s.depth
	return fmt.Sprintf("<Size3D width=%v height=%v depth=%v>", width, height, depth)
}

// 实现 `typedef.Comparable` 接口, 比较两个对象大小
//
// 接口实现原理同 `String` 方法
func (s *Size3D) Compare(other interface{}) int {
	v, ok := other.(*Size3D)
	if !ok {
		panic(errs.ErrInvalidType)
	}
	if s == v {
		return 0
	}
	return int(s.Volume() - v.Volume())
}

// 获取 A 面尺寸, 为一个 `Size` 对象
func (s *Size3D) SideA() *size.Size { return &s.Size }

// 获取 B 面尺寸, 为一个 `Size` 对象
func (s *Size3D) SideB() *size.Size { return size.New(s.Size.Width(), s.depth) }

// 获取 C 面尺寸, 为一个 `Size` 对象
func (s *Size3D) SideC() *size.Size { return size.New(s.Size.Height(), s.depth) }
