package builtin

func Len(s string) int {
	return len([]rune(s))
}

func CharAt(s string, index int) rune {
	return []rune(s)[index]
}
