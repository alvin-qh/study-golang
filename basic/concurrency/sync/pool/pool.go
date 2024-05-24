package pool

import "sync"

// 定义一个池中实例的类型
type PoolElem[T any] struct {
	elem T          // 内容实例
	pool *sync.Pool // 池实例指针
}

// 将元素返回池
func (pe *PoolElem[T]) Release() {
	pe.pool.Put(pe)
}

// 获取池元素中存储的实例
func (pe *PoolElem[T]) Get() T {
	return pe.elem
}

// 基于 `sync.Pool` 类型设置新类型
type Pool[T any] sync.Pool

// 创建实例, 设置创建内容实例
func NewPool[T any](creator func() T) *Pool[T] {
	var pool sync.Pool

	// 设置池的实例创建函数
	pool.New = func() any {
		return &PoolElem[T]{
			elem: creator(),
			pool: &pool,
		}
	}

	return (*Pool[T])(&pool)
}

// 从池中获取池元素实例
func (p *Pool[T]) Get() *PoolElem[T] {
	return (*sync.Pool)(p).Get().(*PoolElem[T])
}
