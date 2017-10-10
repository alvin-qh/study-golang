package basic

import (
	"testing"
	"fmt"
)

func Test(t *testing.T) {
	var s = "Hello, 你好"


	fmt.Println(s[7])

	//var c = []rune(s)
	//fmt.Print(string(c[7]))


	fmt.Println(s[7:9])
}
