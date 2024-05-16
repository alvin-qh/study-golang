package generator

import (
	"errors"
	"runtime"
)

// 生成器错误集合
var (
	ErrChanClosed = errors.New("channel was closed")
)

// 定义生成器结构体
type Generator[T any] struct {
	ch chan T
}

// 定义数据生成函数
type GeneratorFunc[T any] func(ch chan T) T

// 产生一个新的生成器对象
func New[T any](gen GeneratorFunc[T]) *Generator[T] {
	g := &Generator[T]{ch: make(chan T)}

	// 析构函数, 关闭生成器对象
	runtime.SetFinalizer(g, func(g *Generator[T]) { g.Close() })

	go func() {
		defer func() { recover() }()

		// 在协程中异步生成数据
		gen(g.ch)
	}()

	return g
}

// 关闭生成器对象
func (g *Generator[T]) Close() {
	if g.ch != nil {
		close(g.ch)
		g.ch = nil
	}
}

// 获取生成器下一个数据
func (g *Generator[T]) Next(defVal T) (T, error) {
	if g.ch == nil {
		return defVal, ErrChanClosed
	}

	// 从 channel 中接收下一个数据
	v, ok := <-g.ch
	if !ok {
		return defVal, ErrChanClosed
	}
	return v, nil
}


func (g *Generator[T]) Range() <-chan T {
	if g.ch == nil {
		panic(ErrChanClosed)
	}
	return g.ch
}
