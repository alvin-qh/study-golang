package atomic

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试对变量进行原子操作
//
// 可以通过
func TestAtomic_ForVariable(t *testing.T) {
	cond := sync.NewCond(&sync.Mutex{})
	wait := func() {
		cond.L.Lock()
		defer cond.L.Unlock()

		cond.Wait()
	}

	var n int32 = 0

	go func() {
		wait()

		for i := 0; i < 10000; i++ {
			atomic.AddInt32(&n, 1)
		}
	}()

	go func() {
		wait()

		for i := 0; i < 10000; i++ {
			atomic.AddInt32(&n, -1)
		}
	}()

	cond.Broadcast()
	assert.Equal(t, int32(0), n)
}
