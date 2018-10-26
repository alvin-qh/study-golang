package types

func CannotChangeSize(size Size, newWidth int, newHeight int) {
	size.Width = newWidth
	size.Height = newHeight
}

func ChangeSize(size *Size, newWidth int, newHeight int) {
	size.Width = newWidth
	size.Height = newHeight
}

func (i Size) CannotChange(newWidth int, newHeight int) {
	i.Width = newWidth
	i.Height = newHeight
}

func (i *Size) Change(newWidth int, newHeight int) {
	i.Width = newWidth
	i.Height = newHeight
}
