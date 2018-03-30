package subpack

import "fmt"

// init method will run after current package was imported
func init() {
	fmt.Println("Hello, package")
}

func Name() string {
	return "subpack"
}
