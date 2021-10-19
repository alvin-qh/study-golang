package oop

import "basic/types"

// Interface 定义 interface
type Comparable interface {
	Compare(other Comparable) int
}

type Long int64

func (i Long) Compare(other Comparable) int {
	that := other.(Long)
	return int(i - that)
}

type Int int32

func (i Int) Compare(other Comparable) int {
	oi := other.(Int)
	return int(i - oi)
}

type Size struct {
	types.Size
}

func NewSize(width int, height int) Size {
	return Size{types.Size{Width: width, Height: height}}
}

func (i Size) Compare(other Comparable) int {
	that := other.(Size)
	return i.Area() - that.Area()
}

type Size3D struct {
	Size
	Length int
}

func NewSize3D(width int, height int, length int) Size3D {
	return Size3D{Size: NewSize(width, height), Length: length}
}

func (i Size3D) Area() int {
	return i.Size.Area() * i.Length
}

func (i Size3D) Compare(other Comparable) int {
	that, ok := other.(Size3D)
	if ok {
		return i.Area() - that.Area()
	} else {
		that := other.(Size)
		if i.Length == 0 {
			return i.Size.Compare(that)
		}
		return 1
	}
}
