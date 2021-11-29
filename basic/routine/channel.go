package routine

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrLockClosed = errors.New("lock was closed")
)

type Lock struct {
	ch     chan struct{}
	mut    sync.Mutex
	locked atomic.Value
}

func NewLock() *Lock {
	l := &Lock{ch: make(chan struct{}, 1)}
	l.locked.Store(false)
	return l
}

func (l *Lock) Lock(ctx context.Context) bool {
	for !l.locked.CompareAndSwap(false, true) {
		ch := l.channel()

		select {
		case <-ch:
			continue
		case <-ctx.Done():
			return false
		}
	}
	return true
}

func (l *Lock) channel() <-chan struct{} {
	l.mut.Lock()
	defer l.mut.Unlock()

	return l.ch
}

func (l *Lock) Unlock() {
	if l.locked.CompareAndSwap(true, false) {
		defer func() {
			if recover() != nil {
				l.ch = nil
			}
		}()

		ch := l.swapChannel(make(chan struct{}, 1))
		close(ch)
	}
}

func (l *Lock) swapChannel(newOne chan struct{}) chan struct{} {
	l.mut.Lock()
	defer l.mut.Unlock()

	oldOne := l.ch
	l.ch = newOne
	return oldOne
}
