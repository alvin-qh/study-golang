package runtime

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	MEM_PROFILE_FILENAME = "mem_profile.pr"
)

// 测试记录 Profile 数据
func TestRecordProfile(t *testing.T) {
	defer os.Remove(MEM_PROFILE_FILENAME)

	p := NewProfile() // 创建 Profile 对象

	f, err := os.Create(MEM_PROFILE_FILENAME) // 创建记录 profile 信息的文件
	assert.NoError(t, err)
	defer f.Close()

	p.Use(NewMemProfileRecorder(f, 500)) // 使用内存信息记录对象

	p.Start() // 开始记录

	time.Sleep(time.Second) // 留有一段记录时间

	p.Stop() // 结束记录

	stat, err := f.Stat()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, stat.Size(), int64(0))
}
