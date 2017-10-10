package basic

import "fmt"

var A int = 0

func Simple()  {
	A = 100
}

func Arguments(x int, y int) int  {
	return x + y
}

func ReturnMore(x int, y int, m string) (string, string) {
	return fmt.Sprintf("%s:%d", m, x), fmt.Sprintf("%s:%d", m, y)
}