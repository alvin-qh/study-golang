package memory

import (
	"bytes"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试内存 Profile 信息的记录
func TestMemoryProfile(t *testing.T) {
	// 定义缓冲区, 2 MB大小
	buf := bytes.NewBuffer(make([]byte, 0, 1024*1024*2))

	// 创建 Profile 记录实例
	pf := New(buf, 500)

	// 开始记录内存信息
	err := pf.Start()
	assert.NoError(t, err)

	// 分配 100 MB 的内存空间, 用于测试内存 Profile 信息的记录
	mem := make([]byte, 1024*1024*100)

	// 等待 1 秒钟以产生内存 Profile 记录
	time.Sleep(time.Millisecond * 500)

	mem[len(mem)-1] = 0

	mem = nil
	runtime.GC()

	time.Sleep(time.Millisecond * 500)

	pf.Stop()
	assert.Equal(t, buf.String(), "")
}
