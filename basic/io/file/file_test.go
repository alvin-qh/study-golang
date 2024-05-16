package file

import (
	"io"
	"os"
	"study/basic/io/pathex"
	"study/basic/testing/assertion"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试获取一个目录下所有文件和子文件夹的信息
func TestFile_ReadDir(t *testing.T) {
	// 打开路径为 os.File 对象
	dir, err := os.Open(".")
	assert.Nil(t, err)

	// 读取路径下的信息, 返回所有文件 (包括路径) 的 `os.Stat` 对象
	// 参数 `0` 表示不限制返回结果的数量, 否则按所给数量返回结果
	fis, err := dir.Readdir(0)
	assert.Nil(t, err)

	ns := make([]string, 0, len(fis))

	// 遍历目录下所有文件信息
	for _, fi := range fis {
		ns = append(ns, fi.Name())
	}

	// 确认找到的文件均位于当前目录下
	for _, name := range ns {
		assert.NotContains(t, name, "/")
		assert.NotContains(t, name, "\\")
		assertion.FileOrPathExist(t, name)
	}
}

// 测试获取一个目录下所有文件和子文件夹的名称
func TestFile_ReadDirnames(t *testing.T) {
	// 打开路径为 os.File 对象
	dir, err := os.Open(".")
	assert.Nil(t, err)

	// 读取路径下的信息, 返回所有文件 (包括路径) 的 `os.Stat` 对象
	// 参数 `0` 表示不限制返回结果的数量, 否则按所给数量返回结果
	ns, err := dir.Readdirnames(0)
	assert.Nil(t, err)

	// 确认找到的文件均位于当前目录下
	for _, name := range ns {
		assert.NotContains(t, name, "/")
		assert.NotContains(t, name, "\\")
		assertion.FileOrPathExist(t, name)
	}
}

// 文件读写测试
//
// `os.File` 类型实现了 `io.Reader` 和 `io.Writer` 接口, 可以直接对文件进行读写
//
// 除了可以通过 `os.Create` 函数创建 `os.File` 实例外, `os.OpenFile` 函数也可以创建新文件或打开一个现有文件,
// 返回 `os.File` 实例, 且如果是新建文件, 还可以指定文件的访问权限
func TestFile_WriteRead(t *testing.T) {
	name := pathex.RandomFileName()
	defer os.Remove(name)

	// 测试文件写
	{
		// 打开一个文件
		f, err := os.Create(name)
		assert.Nil(t, err)

		defer f.Close()

		// 在文件 0 位置写入字符串, 共 11 字节. 写指针指向 11 位置
		n, err := f.WriteString("Hello World")
		assert.Nil(t, err)
		assert.Equal(t, 11, n)
		assert.Equal(t, int64(11), GetFileLength(f))

		data := []byte{1, 2, 3, 4, 5}

		// 指定在文件 16 位置写入 5 个字节, 写指针不变, 仍指向 11 位置, 文件长度变为 21 字节
		n, err = f.WriteAt(data, 16)
		assert.Nil(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, int64(21), GetFileLength(f))
		assert.Equal(t, int64(11), GetFilePosition(f)) // 获取当前文件指针, 仍为 11

		// 写指针从文件末尾移动 -10 个偏移, 指向 11 位置
		pos, err := f.Seek(-10, io.SeekEnd)
		assert.Nil(t, err)
		assert.Equal(t, int64(11), pos)

		// 从 11 位置写入 5 字节, 写指针指向 16 位置
		n, err = f.Write([]byte{5, 4, 3, 2, 1})
		assert.Nil(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, int64(21), GetFileLength(f))
	}

	// 测试文件读
	{
		// 打开一个只读文件
		f, err := os.Open(name)
		assert.Nil(t, err)

		defer f.Close()

		// 读取文件的 bytes 缓冲区
		buf := make([]byte, GetFileLength(f))

		// 从 0 位置读取 11 字节, 读指针移动到 11
		n, err := f.Read(buf[0:11])
		assert.Nil(t, err)
		assert.Equal(t, 11, n)
		assert.Equal(t, "Hello World", string(buf[0:11])) // 读出字符串内容

		// 在文件 16 位置读取 5 字节, 读指针不变
		n, err = f.ReadAt(buf[16:21], 16)
		assert.Nil(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, []byte{1, 2, 3, 4, 5}, buf[16:21])
		assert.Equal(t, int64(11), GetFilePosition(f)) // 读指针仍在 11

		// 从文件开始位置, 将读指针移动到 11 位置
		pos, err := f.Seek(11, io.SeekStart)
		assert.Nil(t, err)
		assert.Equal(t, int64(11), pos)

		// 在 11 位置读取 5 字节, 读指针移动到 16
		n, err = f.Read(buf[11:16])
		assert.Nil(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, []byte{5, 4, 3, 2, 1}, buf[11:16])
	}
}

// 测试文件截取
//
// 文件截取即将文件截断为指定长度, 多余的部分将被丢弃
//
// Go 提供了两种截取文件的方法:
//  1. 通过 `os.File` 实例的 `Truncate` 方法, 需要打开文件, 得到一个 `os.File` 对象;
//  2. 通过 `os.Truncate` 函数. 无需打开文件;
func TestFile_Truncate(t *testing.T) {
	name := pathex.RandomFileName()
	defer os.Remove(name)

	// 创建文件
	f, err := os.Create(name)
	assert.Nil(t, err)

	defer f.Close()

	// 写入10个字符
	_, err = f.WriteString("1234567890")
	assert.Nil(t, err)
	assert.Equal(t, int64(10), GetFileLength(f))

	// 截断文件, 为原始长度的一半
	err = f.Truncate(5)
	assert.Nil(t, err)
	assert.Equal(t, int64(5), GetFileLength(f))

	// 指针移动到文件开头
	_, err = f.Seek(0, io.SeekStart)
	assert.Nil(t, err)

	// 读取全部文件内容, 只剩下写入内容的一半
	s, err := io.ReadAll(f)
	assert.Nil(t, err)
	assert.Equal(t, "12345", string(s))
}
