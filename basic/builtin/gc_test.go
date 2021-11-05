package builtin

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type memory struct {
	buf []int
}

func (m *memory) alloc(n int) {
	m.buf = make([]int, n, n*2)
}

func (m *memory) size() int {
	return len(m.buf)
}

func (m *memory) clear() {
	m.buf = nil
}

// 测试 GC
func TestGc(t *testing.T) {
	var ms runtime.MemStats // 内存状态结构体

	runtime.ReadMemStats(&ms) // 获取内存状态
	fmt.Printf("Alloc: %d(bytes) HeapIdle: %d(bytes) HeapReleased: %d(bytes)\n", ms.Alloc, ms.HeapIdle, ms.HeapReleased)

	m := memory{}

	finalized := false
	runtime.SetFinalizer(&m, func(m *memory) { // 调用 GC 时，会同时调用设置在对象上的 finalizer 函数，释放资源
		finalized = true
	})

	m.alloc(10000) // 分配内存
	assert.Equal(t, 10000, m.size())

	m.clear()    // 释放内存
	runtime.GC() // 手动唤起 GC 释放堆

	runtime.ReadMemStats(&ms) // 获取内存状态
	fmt.Printf("Alloc: %d(bytes) HeapIdle: %d(bytes) HeapReleased: %d(bytes)\n", ms.Alloc, ms.HeapIdle, ms.HeapReleased)

	assert.True(t, finalized)
}
