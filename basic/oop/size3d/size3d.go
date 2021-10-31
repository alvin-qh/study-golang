package size3d

import (
	"basic/oop/size"
	"basic/oop/typedef"
	"fmt"
)

type Size3D struct {
	size.Size
	depth float64
}

func New(width, height, depth float64) *Size3D { return new(Size3D).Init(width, height, depth) }

func (s *Size3D) Init(width, height, depth float64) *Size3D {
	s.Size.Init(width, height)
	s.depth = depth
	return s
}

func (s *Size3D) Value() (width, height, depth float64) {
	width, height = s.Size.Value()
	depth = s.depth
	return
}

func (s *Size3D) Area() float64 {
	width, height := s.Size.Value()
	return (s.Size.Area() + width*s.depth + height*s.depth) * 2
}

func (s *Size3D) Volume() float64 { return s.Size.Area() * s.depth }

func (s *Size3D) ToString() string {
	width, height, depth := s.Value()
	return fmt.Sprintf("<Size3D width=%v height=%v depth=%v>", width, height, depth)
}

func (s *Size3D) Compare(other interface{}) int {
	v, ok := other.(*Size3D)
	if !ok {
		panic(typedef.ErrType)
	}
	if s == v {
		return 0
	}
	return int(s.Volume() - v.Volume())
}

func (s *Size3D) SideA() *size.Size { return &s.Size }

func (s *Size3D) SideB() *size.Size { return size.New(s.Size.Width(), s.depth) }

func (s *Size3D) SideC() *size.Size { return size.New(s.Size.Height(), s.depth) }
