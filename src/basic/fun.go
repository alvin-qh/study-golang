package basic

var A int = 0

func Simple() {
	A = 100
}

func Arguments(x int, y int) int {
	return x + y
}

func ExchangeByReturn(x interface{}, y interface{}) (interface{}, interface{}) {
	return y, x
}
