package basic

func CharAt(s string, n int) rune {
	as := []rune(s)
	return as[n]
}
