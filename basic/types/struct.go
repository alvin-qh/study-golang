package types

type Size struct {
	Width  int
	Height int
}

func (i Size) Area() int {
	return i.Width * i.Height
}
