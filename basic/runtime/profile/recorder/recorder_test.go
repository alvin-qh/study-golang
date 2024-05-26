package record

import (
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	CPU_PROFILE_FILENAME  = "cpu.profile"
	HEAP_PROFILE_FILENAME = "heap.profile"
)

// 测试记录 Profile 数据
func TestProfile_Record(t *testing.T) {
	defer os.Remove(CPU_PROFILE_FILENAME)
	defer os.Remove(HEAP_PROFILE_FILENAME)

	// 创建文件用于记录 CPU Profile
	cpuf, err := os.Create(CPU_PROFILE_FILENAME)
	assert.Nil(t, err)
	defer cpuf.Close()

	// 创建文件用于记录 Heap Profile
	heapf, err := os.Create(HEAP_PROFILE_FILENAME)
	assert.Nil(t, err)
	defer heapf.Close()

	// 创建 Profile 记录器实例
	r := New(cpuf, heapf)

	// 开始 CPU 信息记录
	err = r.Start()
	assert.Nil(t, err)

	// 进行高负载运算
	x := int64(0)
	for i := 0; i < 10000; i++ {
		x += int64(math.Pow(float64(x), 2))
	}

	// 停止 CPU 信息记录, 并记录堆内存信息
	err = r.Stop()
	assert.Nil(t, err)
}
