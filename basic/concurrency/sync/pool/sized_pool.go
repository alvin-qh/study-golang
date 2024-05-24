package pool

import (
	"sync"
	"sync/atomic"
)

type SizedPool[T any] struct {
	pool    sync.Pool
	size    atomic.Int64
	maxSize int64
}

func NewPool[T any](maxSize int, create func() T) *SizedPool[T] {
	return &SizedPool[T]{
		pool: sync.Pool{
			New: func() any {
				return create()
			},
		},
		maxSize: int64(maxSize),
	}
}

func (p *SizedPool[T]) Get() (*T, error) {
	size := p.size.Load()

	for size < p.maxSize {
		if p.size.CompareAndSwap(size, size+1) {
            obj := p.pool.Get().(T)
			return &obj, nil
		}
		size = p.size.Load()
	}
	return nil, ErrPoolFull
}
