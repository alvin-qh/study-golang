package main

import (
	"fmt"

	meta2 "study/module/demo/meta"

	meta1 "gitee.com/go-common-libs/demo-module/meta"
)

func main() {
	fmt.Printf("Module from gitee.com is: \"%v\"\n", meta1.Version())
	fmt.Printf("Module from locale is: \"%v\"\n", meta2.Version())
}
