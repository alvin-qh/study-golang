package file

import (
	"bufio"
	"encoding/gob"
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

	file, err := os.OpenFile(FILE_NAME, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755) // 打开 (或截断) 一个文件
	assert.NoError(t, err)
	assert.Equal(t, FILE_NAME, file.Name())

	defer CloseAndRemoveFile(file) // 函数结束时关闭并删除文件

	n, err := file.WriteString("Hello World") // 在文件 0 位置写入字符串, 共 11 字节. 写指针指向 11 位置
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, 11, FileLength(file))

	data := []byte{1, 2, 3, 4, 5}
	n, err = file.WriteAt(data, 16) // 指定在文件 16 位置写入 5 个字节, 写指针不变, 仍指向 11 位置, 文件长度变为 21 字节
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, 21, FileLength(file))
	assert.Equal(t, int64(11), GetFileCursor(file)) // 获取当前文件指针, 仍为 11

	cur, err := file.Seek(-10, io.SeekEnd) // 写指针从文件末尾移动 -10 个偏移, 指向 11 位置
	assert.NoError(t, err)
	assert.Equal(t, int64(11), cur)

	n, err = file.Write([]byte{5, 4, 3, 2, 1}) // 从 11 位置写入 5 字节, 写指针指向 16 位置
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, 21, FileLength(file))

	file.Close()

	// 测试文件读

	file, err = os.OpenFile(FILE_NAME, os.O_RDONLY, 0) // 打开一个只读文件
	assert.NoError(t, err)
	assert.Equal(t, FILE_NAME, file.Name())

	data = make([]byte, FileLength(file)) // 读取文件的 bytes 缓冲区

	n, err = file.Read(data[0:11]) // 从 0 位置读取 11 字节, 读指针移动到 11
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, "Hello World", string(data[0:11])) // 读出字符串内容

	n, err = file.ReadAt(data[16:21], 16) // 在文件 16 位置读取 5 字节, 读指针不变
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, data[16:21])
	assert.Equal(t, int64(11), GetFileCursor(file)) // 读指针仍在 11

	cur, err = file.Seek(11, io.SeekStart) // 从文件开始位置, 将读指针移动到 11 位置
	assert.NoError(t, err)
	assert.Equal(t, int64(11), cur)

	n, err = file.Read(data[11:16]) // 在 11 位置读取 5 字节, 读指针移动到 16
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte{5, 4, 3, 2, 1}, data[11:16])
}

// 通过缓冲进行 IO 操作
// bufio.Writer 实现了 io.Writer, io.ByteWriter, io.StringWriter 以及 io.RuneWriter 接口, 用于包装另一个 io.Writer 对象, 为其增加缓存支持
// bufio.Reader 实现了 io.Reader, io.ByteReader, io.StringWriter 以及 io.RuneWriter 接口, 用于包装另一个 io.Reader 对象, 为其增加缓存支持
// bufio.ReadWriter 同时包装了一个 io.Writer 和一个 io.Reader 对象, 实现了同时读写的操作
func TestBufferedFileIO(t *testing.T) {
	// 利用缓冲流写入文件
	file, err := os.Create(FILE_NAME) // 打开用于写入的文件
	assert.NoError(t, err)

	defer CloseAndRemoveFile(file) // 函数结束后关闭并删除文件

	bufW := bufio.NewWriterSize(file, 256) // 创建缓存大小 256 字节的 bufio.Writer 对象, 用于对文件进行写操作. 另外, bufio.NewWriter(file) 用于创建一个缓冲区大小为 4096 的 bufio.Writer 对象
	assert.Equal(t, 256, bufW.Size())      // 已缓存数据大小
	assert.Equal(t, 256, bufW.Available()) // 缓存总大小
	assert.Equal(t, 0, bufW.Buffered())    // 已缓存数据大小

	bufW.WriteString("Hello\n")            // 写入字符串
	assert.Equal(t, 250, bufW.Available()) // 缓存剩余大小
	assert.Equal(t, 6, bufW.Buffered())    // 写入数据后, 已缓存数据大小

	pu := NewUser(1, "Alvin", "alvin@fake.com", []string{"13999912345", "13000056789"})
	enc := gob.NewEncoder(bufW)            // 在 bufio.Writer 对象的基础上, 包装一个 gob.Encoder 对象
	enc.Encode(pu)                         // 将对象编码后写入 bufio.Writer 对象
	assert.Equal(t, 116, bufW.Available()) // 缓存剩余大小
	assert.Equal(t, 140, bufW.Buffered())  // 写入数据后, 已缓存数据大小, 共 140 字节

	err = bufW.Flush() // 将缓存的内容写入文件
	assert.NoError(t, err)
	assert.Equal(t, 256, bufW.Available()) // 缓存剩余大小
	assert.Equal(t, 0, bufW.Buffered())    // 缓存数据已经全部写入文件

	file.Close() // 关闭文件

	// 利用缓冲流读文件
	file, err = os.Open(FILE_NAME) // 打开文件用于读
	assert.NoError(t, err)

	bufR := bufio.NewReaderSize(file, 256) // 创建缓存大小 256 字节的 bufio.Reader 对象, 用于对文件进行读操作. 另外, bufio.NewReader(file) 用于创建一个缓冲区大小为 4096 的 bufio.Reader 对象
	assert.Equal(t, 0, bufR.Buffered())    // 已缓存数据 0 字节 (尚未开始读取)

	line, prefix, err := bufR.ReadLine() // 读取一行字符串 (以 \n 结尾的一行数据)
	assert.NoError(t, err)
	assert.Equal(t, "Hello", string(line)) // 读取 1 行的内容
	assert.False(t, prefix)                // prefix 返回 false, 表示以读完一行, 否则表示一行尚未读完, 需要再次调用 ReadLine 函数, 直到读完一行
	assert.Equal(t, 134, bufR.Buffered())  // 缓冲区使用 134 字节. 具体情况为: 当第一次读取时, 先读入缓存 140 字节, 读走 1 行 6 字节后, 剩余 134 字节

	u := User{}
	dec := gob.NewDecoder(bufR) // 在 bufio.Reader 对象基础上创建 gob.Decoder 对象

	err = dec.Decode(&u) // 从文件中解码一个 user.User 对象
	assert.NoError(t, err)
	assert.Equal(t, *pu, u)
	assert.Equal(t, 0, bufR.Buffered()) // 缓冲区已读完

	// 测试同时读写
	fw, err := os.Create(FILE_NAME) // 创建用于写的 os.File 对象
	assert.NoError(t, err)
	defer fw.Close()

	fr, err := os.Open(FILE_NAME) // 创建用于读的 os.File 对象
	assert.NoError(t, err)
	defer fr.Close()

	bufRW := bufio.NewReadWriter(bufio.NewReader(fr), bufio.NewWriter(fw)) // 通过一对 bufio.Reader 和 bufio.Writer 对象, 创建同时进行读写的 bufio.ReaderWriter 对象
	bufRW.WriteString("Hello\n")                                           // 写入字符串

	pu = NewUser(1, "Alvin", "alvin@fake.com", []string{"13999912345", "13000056789"}) // 通过 gob.Encoder 写入结构体
	enc = gob.NewEncoder(bufRW)
	enc.Encode(pu)

	bufRW.Flush() // 将缓冲区内容写入文件

	line, prefix, err = bufRW.ReadLine() // 读取一行字符串
	assert.NoError(t, err)
	assert.Equal(t, "Hello", string(line))
	assert.False(t, prefix)

	u = User{}

	dec = gob.NewDecoder(bufRW) // 通过 gob.Decoder 读取结构体
	dec.Decode(&u)
	assert.Equal(t, *pu, u)
}
