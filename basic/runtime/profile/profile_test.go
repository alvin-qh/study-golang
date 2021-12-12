package profile

import (
	"os"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/**
 * 在测试中使用 profile 记录性能数据
 *
 * 记录 CPU 使用情况：$ go test ./runtime/profile -test.cpuprofile cpu.profile  # 记录测试的 CPU 使用情况，保存在 cpu.profile 文件中
 * 记录 MEM 使用情况：$ go test ./runtime/profile -test.memprofile mem.profile  # 记录测试的 内存 使用情况，保存在 mem.profile 文件中
 */

const (
	MEM_PROFILE_FILENAME  = "mem.profile"
	CPU_PROFILE_FILENAME  = "cpu.profile"
	HEAP_PROFILE_FILENAME = "heap.profile"

	frequency = 500
)

// 测试记录 Profile 数据
// cspell: ignore memf cpuf heapf
func TestRecordProfile(t *testing.T) {
	defer os.Remove(MEM_PROFILE_FILENAME)
	defer os.Remove(CPU_PROFILE_FILENAME)
	defer os.Remove(HEAP_PROFILE_FILENAME)

	p := NewProfile() // 创建 Profile 对象

	memf, err := os.Create(MEM_PROFILE_FILENAME) // 创建记录 profile 信息的文件
	assert.NoError(t, err)
	defer memf.Close()

	p.Use(NewMemProfileRecorder(memf, frequency)) // 使用内存信息记录对象

	cpuf, err := os.Create(CPU_PROFILE_FILENAME)
	assert.NoError(t, err)
	defer cpuf.Close()

	p.Use(NewCpuProfileRecorder(cpuf, frequency)) // 使用CPU信息记录对象

	p.Start() // 开始记录

	data := make([]int64, 0)
	for i := 0; i < 1e8; i++ {
		data = append(data, int64(i))
	}
	assert.Len(t, data, 1e8)

	heapf, err := os.Create(HEAP_PROFILE_FILENAME)
	assert.NoError(t, err)

	err = pprof.WriteHeapProfile(heapf) // 记录堆内存使用情况
	assert.NoError(t, err)

	time.Sleep(time.Second) // 留有一段记录时间

	p.Stop() // 结束记录
}
