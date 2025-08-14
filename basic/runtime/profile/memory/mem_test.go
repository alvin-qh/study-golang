package memory

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试内存 Profile 信息的记录
func TestProfile_Memory(t *testing.T) {
	t.Skip("skip")

	// 定义缓冲区, 2MB 大小
	buf := bytes.NewBuffer(make([]byte, 0, 1024*1024*10))
	// w, _ := os.Create("d.txt")

	// 创建 Profile 记录实例
	pf := New(buf, 500)

	// 开始记录内存信息
	err := pf.Start()
	assert.Nil(t, err)

	// 分配 100 MB 的内存空间, 用于测试内存 Profile 信息的记录
	mem := make([]byte, 1024*1024*100)

	// 等待 1 秒钟以产生内存 Profile 记录
	time.Sleep(500 * time.Millisecond)
	mem[len(mem)-1] = 0

	// 停止记录内存信息
	pf.Stop()
	time.Sleep(500 * time.Millisecond)

	fmt.Println(buf.String())

	// 代码第 25 行进行了一次内存分配, 确认记录中包含该位置
	hit := bytes.Index(buf.Bytes(), []byte("study-golang/basic/runtime/profile/memory/mem_test.go:26"))
	assert.True(t, hit >= 0)
}
