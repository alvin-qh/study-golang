package types

const A = 100

var B int64

func Add(n int) int64 {
	B = B + A * int64(n)
	return B
}
