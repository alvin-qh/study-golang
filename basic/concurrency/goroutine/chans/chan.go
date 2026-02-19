package chans

// 返回一个 chan, 可以读取通过 fn 函数不断的产生数据
//
// 参数 fn 是一个函数, 以 chan 作为参数, 可以通过该 chan 不断的向外发送数据, 直到 fn 函数执行完成, chan 被关闭
//
// 返回值是一个 chan, 可以通过该 chan 读取 fn 函数发送的数据, 直到 fn 函数执行完成, chan 被关闭
func Generator[T any](fn func(c chan T)) chan T {
	c := make(chan T)
	go func() {
		defer close(c)

		fn(c)
	}()
	return c
}
