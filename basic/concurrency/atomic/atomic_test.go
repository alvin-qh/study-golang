package atomic

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试对变量进行原子操作
//
// 可以通过 `atomic` 包提供的函数对变量进行原子操作, 以保证在并发环境下对变量的访问和修改是安全的,
// 不会出现竞态条件或数据不一致的问题
func TestAtomic_Operators(t *testing.T) {
	// 定义一个原子变量
	var atomN int32 = 0

	// 实例化一个等待组对象, 并添加 2 个需要等待的 goroutine
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 启动第一个 goroutine, 分别对原子变量进行加 1 操作, 共执行 10000 次
	go func() {
		for range 10000 {
			atomic.AddInt32(&atomN, 1)
		}
		wg.Done()
	}()

	// 启动另一个 goroutine, 分别对原子变量进行减 1 操作, 共执行 10000 次
	go func() {
		for range 10000 {
			atomic.AddInt32(&atomN, -1)
		}
		wg.Done()
	}()

	// 等待两个 goroutine 执行完成
	wg.Wait()

    // 确认两个 goroutine 对原子变量的加 1 和减 1 操作相互抵消, 最终原子变量的值仍然为 0
	assert.Equal(t, int32(0), atomN)
}
