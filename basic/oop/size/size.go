package size

import (
	"basic/oop/typedef"
	"fmt"
)

type Size struct {
	width  float64
	height float64
}

func New(width, height float64) *Size { return new(Size).Init(width, height) }

func (s *Size) Init(width, height float64) *Size {
	s.width = width
	s.height = height
	return s
}

func (s *Size) ToString() string { return fmt.Sprintf("<Size width=%v height=%v>", s.width, s.height) }

func (s *Size) Value() (width, height float64) {
	width, height = s.width, s.height
	return
}

func (s *Size) Area() float64 { return s.width * s.height }

func (s *Size) Compare(other interface{}) int {
	v, ok := other.(*Size)
	if !ok {
		panic(typedef.ErrType)
	}
	if s == v {
		return 0
	}
	return int(s.Area() - v.Area())
}

func (s *Size) Clone() *Size {
	ns := *s
	return &ns
}

func (s *Size) Width() float64 { return s.width }

func (s *Size) Height() float64 { return s.height }
