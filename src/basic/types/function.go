package types

func ReturnOne() string {
	return "First"
}

func ReturnTwo() (string, string) {
	return "First", "Second"
}

func ReturnFunc() func(a int, b int) int {
	return func(a int, b int) int {
		return a + b
	}
}

type FuncType func(a interface{}, b interface{}) interface{}

func ArgumentAsFunction(a interface{}, b interface{}, callback FuncType) interface{} {
	return callback(a, b)
}
