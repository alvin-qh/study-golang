package main

import (
	"fmt"

	"gitee.com/go-common-libs/demo-module/meta"
)

func main() {
	fmt.Printf("Module from gitee.com is: \"%v\"\n", meta.Version())
}
