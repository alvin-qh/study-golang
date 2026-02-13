package function

import "math"

// 定义函数
func Add(a, b int) int {
	return a + b
}

// 定义具备单一返回值的函数
func Pow(n int) float64 {
	return math.Pow(float64(n), 2)
}

// 定义同时返回两个值的函数
func Sqrt(n int) (float64, bool) {
	if n <= 0 {
		return 0, false
	}
	return math.Sqrt(float64(n)), true
}

// 命名返回值
//
// 可以给函数的返回值命名, 在函数体内部, 命名返回值相当于一个变量, 给其赋的值将会作为函数的返回值
func NumAdd(a int, b int) (sum int) {
	// 给明名返回值赋值
	sum = a + b

	if sum == 0 {
		// 可以在 return 语句中显示的书写返回的值
		return sum
	}

	// 可以省略书写返回值, 但 return 语句不能省略
	return
}

// 命名多个返回值
//
// 对于同时返回多个值的情况, 可以为每个返回值命名
func NumAddAndSub(a int, b int) (sum int, sub int) {
	sum = a + b
	sub = a - b

	if sum != sub {
		// 显式返回所有值
		return sum, sub
	}

	// 省略显式返回值, 但命名返回值已经赋值作为返回值
	return
}

// 仅命名部分返回值
//
// 如果某个参数之前的参数类型一致, 则可以将这部分参数用 ',' 分隔, 并仅在最后一个参数上声明类型
//
// 可以混合命名返回值和匿名返回值, 但匿名返回值必须在命名返回值之前
// 此时 return 语句必须显式返回所需的值
func NumAddAndSubForm2(a, b int) (int, sub int) {
	sub = a - b
	return a + b, sub
}

// 使用不定参数
//
// 在参数名后使用 ... 表示该参数为不定参数, 不定参数可以支持 0 或 多个 指定类型的参数
// 不定参数在函数体内部表现为一个切片
func AddForVarargs(a ...int) (sum int) {
	sum = 0
	for _, n := range a {
		sum += n
	}
	return
}

// 将函数类型作为参数传递
//
// 定义具有函数类型参数的函数, 即将参数类型作为参数传递
func Callback(fn func(int, int) int, a, b int) int {
	r := fn(a, b)
	return r
}

// 将函数类型作为返回值返回
//
// 从返回值获取一个函数, 返回值可以为匿名或命名的
func GetExecutor() (exec func(a, b int) int) {
	exec = NumAdd
	return
}
