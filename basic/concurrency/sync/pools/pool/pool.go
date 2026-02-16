package pool

import "sync"

// 定义一个对象池元素类型
// 该类型对象存储在对象池中, 起到管理对象的作用
type Elem[T any] struct {
	elem T          // 内容实例
	pool *sync.Pool // 池实例指针
}

// 释放对象, 令对象返回对象池
func (pe *Elem[T]) Release() { pe.pool.Put(pe) }

// 获取池元素中存储的实例
func (pe *Elem[T]) Get() T { return pe.elem }

// 基于 `sync.Pool` 类型设置新类型
type Pool[T any] sync.Pool

// 创建实例, 设置创建内容实例
//
// `new` 参数: 用于创建池元素中存储的内容实例
//
// 该函数会被池在需要创建新实例时调用
//
// 该函数返回一个 `Pool` 对象, 该对象是基于 `sync.Pool` 类型设置的新类型
func New[T any](new func() T) *Pool[T] {
	var pool sync.Pool

	// 设置池的实例创建函数
	pool.New = func() any {
		return &Elem[T]{
			elem: new(),
			pool: &pool,
		}
	}
	return (*Pool[T])(&pool)
}

// 从池中获取池元素实例
//
// 该函数返回一个池元素实例, 该实例中存储了内容实例
func (p *Pool[T]) Get() *Elem[T] {
	return (*sync.Pool)(p).Get().(*Elem[T])
}
