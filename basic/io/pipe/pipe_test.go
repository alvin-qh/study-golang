package pipe

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试创建管道实例
func TestPipe_New(t *testing.T) {
	// 创建管道实例
	pi, err := New()
	assert.Nil(t, err)

	// 结束时关闭管道实例
	defer pi.Close()

	// 获取管道文件名称
	assert.Equal(t, "|0", pi.Reader().Name())
	assert.Equal(t, "|1", pi.Writer().Name())

	// 获取管道读接口文件信息
	fi, err := pi.Reader().Stat()
	assert.Nil(t, err)
	assert.False(t, fi.IsDir())
	assert.Equal(t, int64(0), fi.Size())

	// 获取管道写接口文件信息
	fi, err = pi.Reader().Stat()
	assert.Nil(t, err)
	assert.False(t, fi.IsDir())
	assert.Equal(t, int64(0), fi.Size())
}

// 测试管道读写
func TestPipe_WriteAndRead(t *testing.T) {
	// 创建管道实例
	pi, err := New()
	assert.Nil(t, err)

	// 结束后关闭管道
	defer pi.Close()

	// 启动并发任务向管道写入
	go func() {
		// 向管道写入数据
		for _, d := range [][]byte{[]byte("one\n"), []byte("two\n"), []byte("three")} {
			// 写入数据
			n, err := pi.Write(d)
			assert.Nil(t, err)
			assert.Equal(t, len(d), n)
		}

		// 写入完成, 关闭写管道
		pi.CloseWriter()
	}()

	// 定义读缓冲区
	data := bytes.NewBuffer(make([]byte, 0, 1024))

	// 将管道内容读取到缓冲区
	err = pi.ReadTo(data)
	assert.Nil(t, err)
	// 确认读取内容正确
	assert.Equal(t, "one\ntwo\nthree", data.String())
}
