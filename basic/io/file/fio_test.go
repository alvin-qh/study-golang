package file

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	FILE_NAME = "test.dat"
)

// 文件读写测试
func TestFileIO(t *testing.T) {
	// 测试文件写

	file, err := os.OpenFile(FILE_NAME, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755) // 打开（或截断）一个文件
	assert.NoError(t, err)
	assert.Equal(t, FILE_NAME, file.Name())

	defer CloseAndRemoveFile(file) // 函数结束时关闭并删除文件

	n, err := file.WriteString("Hello World") // 在文件 0 位置写入字符串，共 11 字节。写指针指向 11 位置
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, 11, FileLength(file))

	data := []byte{1, 2, 3, 4, 5}
	n, err = file.WriteAt(data, 16) // 指定在文件 16 位置写入 5 个字节，写指针不变，仍指向 11 位置，文件长度变为 21 字节
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, 21, FileLength(file))
	assert.Equal(t, int64(11), GetFileCursor(file)) // 获取当前文件指针，仍为 11

	cur, err := file.Seek(-10, io.SeekEnd) // 写指针从文件末尾移动 -10 个偏移，指向 11 位置
	assert.NoError(t, err)
	assert.Equal(t, int64(11), cur)

	n, err = file.Write([]byte{5, 4, 3, 2, 1}) // 从 11 位置写入 5 字节，写指针指向 16 位置
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, 21, FileLength(file))

	file.Close()

	// 测试文件读

	file, err = os.OpenFile(FILE_NAME, os.O_RDONLY, 0) // 打开一个只读文件
	assert.NoError(t, err)
	assert.Equal(t, FILE_NAME, file.Name())

	data = make([]byte, FileLength(file)) // 读取文件的 bytes 缓冲区

	n, err = file.Read(data[0:11]) // 从 0 位置读取 11 字节，读指针移动到 11
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, "Hello World", string(data[0:11])) // 读出字符串内容

	n, err = file.ReadAt(data[16:21], 16) // 在文件 16 位置读取 5 字节，读指针不变
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, data[16:21])
	assert.Equal(t, int64(11), GetFileCursor(file)) // 读指针仍在 11

	cur, err = file.Seek(11, io.SeekStart) // 从文件开始位置，将读指针移动到 11 位置
	assert.NoError(t, err)
	assert.Equal(t, int64(11), cur)

	n, err = file.Read(data[11:16]) // 在 11 位置读取 5 字节，读指针移动到 16
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte{5, 4, 3, 2, 1}, data[11:16])
}
