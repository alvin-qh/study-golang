package generic

func GenericIntFloatAdd[T ~int | float64](a T, b T) T {
	return a + b
}

type Number interface {
	~int | ~int8 | ~int32 | ~int64 | ~float32 | ~float64 | ~complex64
}

func GenericAdd[T Number](a T, b T) T {
	return a + b
}
