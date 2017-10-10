package subpack

import "fmt"

func init() {
	fmt.Println("Hello, package")
}

func Name() string {
	return "subpack"
}
