package generic

// 定义泛型方法
//
// 该方法接受 `int` 和 `float64` 类型参数, 并返回相同类型结果, 其中:
//   - `T ~int | float64` 表示定义一个泛型参数 `T`, 类型可为 `int` 或 `float64`
//   - `~int` 表示接受 `int` 及其衍生类型
func GenericIntFloatAdd[T ~int | float64](a T, b T) T {
	return a + b
}

// 定义泛型接口类
// 该接口类型可以可以表示 `~int`, `~int8`, `~int32`, `~int64`, `~float32`, `~float64`, `~complex64` 这些类型中的任意一个
type Number interface {
	~int | ~int8 | ~int32 | ~int64 | ~float32 | ~float64 | ~complex64
}

// 通过 `Number` 接口定义泛型类型
// 由此, 所有能被 `Number` 接口表示的类型都可以作为该方法的参数和返回值
func GenericAdd[T Number](a T, b T) T {
	return a + b
}

// 定义一个泛型切片
type GenericSlice[T ~string | ~int | ~float32 | ~float64] []T
