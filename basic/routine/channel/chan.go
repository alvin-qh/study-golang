package channel

// 返回一个 chan, 可以读取通过 fn 函数不断的产生数据
func Generator[T any](fn func(c chan T)) chan T {
	c := make(chan T)
	go func() {
		fn(c)
		close(c)
	}()
	return c
}
