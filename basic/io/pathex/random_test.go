package pathex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试生成随机文件名
func TestFileEx_RandomFileName(t *testing.T) {
	// 生成随机文件名, 具备默认的前缀和扩展名
	n := RandomFileName()
	assert.Regexp(t, `f-\d{17}\.txt`, n)

	// 确认以 10ms 为间隔生成不重复的文件名
	for i := 0; i < 10; i++ {
		nn := RandomFileName()
		assert.NotEqual(t, n, nn)

		n = nn
		time.Sleep(10 * time.Millisecond)
	}

	// 为随机文件吗指定前缀和扩展名
	n = RandomFileName(WithPrefix("x-"), WithExt(".png"))
	assert.Regexp(t, `x-\d{17}\.png`, n)
}

// 生成随机目录名
func TestFileEx_RandomDirName(t *testing.T) {
	// 生成随机目录名, 具备默认的前缀
	n := RandomDirName()
	assert.Regexp(t, `d-\d{17}`, n)

	// 确认以 10ms 为间隔生成不重复的目录名
	for i := 0; i < 10; i++ {
		nn := RandomDirName()
		assert.NotEqual(t, n, nn)

		n = nn
		time.Sleep(10 * time.Millisecond)
	}

	// 为随机目录名设置前缀
	n = RandomDirName(WithPrefix("x-"))
	assert.Regexp(t, `x-\d{17}`, n)
}
