package size

import (
	"fmt"
	errs "study/basic/oop/errors"
)

// 定义 Size 结构体
type Size struct {
	width  float64 // 宽度
	height float64 // 高度
}

// 构造函数, 产生 Size 对象
func New(width, height float64) *Size {
	return new(Size).Init(width, height)
}

// 初始化 Size 结构体对象
func (s *Size) Init(width, height float64) *Size {
	s.width = width
	s.height = height
	return s
}

// 将结构体转为 字符串
func (s *Size) String() string { return fmt.Sprintf("<Size width=%v height=%v>", s.width, s.height) }

// 获取面积值
func (s *Size) Area() float64 { return s.width * s.height }

// 实现 typedef.Comparable 接口, 比较两个对象大小
func (s *Size) Compare(other interface{}) int {
	v, ok := other.(*Size)
	if !ok {
		panic(errs.ErrInvalidType)
	}
	if s == v {
		return 0
	}
	return int(s.Area() - v.Area())
}

// 获取 width 属性
func (s *Size) Width() float64 { return s.width }

// 获取 height 属性
func (s *Size) Height() float64 { return s.height }
