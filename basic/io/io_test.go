package io

import (
	"basic/io/user"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

// io 包中和读取相关的接口
//  io.Reader：基本的读取操作，从流中读取一系列 bytes
//  io.ReaderAt：随机读取操作，读取流中任意起始位置 bytes
//  io.ByteReader：读取 1 个 byte
//  io.ByteScanner：获取剩余未读取的 bytes
//  io.RuneReader：读取 1 个 rune
//  io.RuneScanner：获取剩余未读取的 runes
//  io.Seeker：随机移动读取指针
//  io.WriterTo：将内容写入另一个 io.Writer 接口对象中
//  io.Closer：关闭当前 Reader 对象
//
// io 包中和写入相关的接口
//  io.Writer：基本的写操作，在流中顺序写入 bytes
//  io.WriterAt：随机写操作，在流的任意位置写入 bytes
//  io.StringWriter：写入字符串操作
//  io.ReadFrom：从另一个 io.Reader 对象中读取内容写入当前对象中

// 测试 strings.Reader
// strings.Reader 接收一个字符串参数，返回一个 Reader 对象，该对象实现了如下接口：
//  io.Reader, io.ReaderAt, io.ByteReader, io.ByteScanner, io.RuneReader, io.RuneScanner, io.Seeker, io.WriterTo
func TestStringIO(t *testing.T) {
	// 实例化一个新的 Reader 对象，用于对字符串内容进行读取操作
	reader := strings.NewReader(`<html>
    <body>
        <div>Hello World</div>
    </body>
</html>`)

	assert.Equal(t, 68, reader.Len()) // io.Reader 接口函数，共 68 字节可以读取

	// 实例化一个 68 字节切片，用于接收读取内容
	data := make([]byte, 68)

	// 读取前 7 个字节（0~6）
	n, err := reader.Read(data[:7]) // io.Reader 接口函数，读取指定长度的 byte 集合
	assert.NoError(t, err)
	assert.Equal(t, 7, n)
	assert.Equal(t, "<html>\n", string(data[:7]))

	// 跳过 7 个字节，读取接下来的 11 个字节（7~17）
	n, err = reader.ReadAt(data[7:18], 7) // io.ReaderAt 接口函数，从流的任意位置开始，读取指定长度的 byte 集合，ReadAt 函数不会移动读取指针
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, "    <body>\n", string(data[7:18]))

	// 移动读取指针
	// 因为 io.ReaderAt 不会移动读取指针，指针停留在 6 的位置
	// 从当前位置移动 11 个字节，指针到达 17 的位置
	reader.Seek(11, io.SeekCurrent) // io.Seeker 接口函数，移动读取指针

	// 读取接下来的 31 个 byte
	for i := 0; i < 31; i++ {
		data[18+i], err = reader.ReadByte() // io.ByteReader 接口函数，读取一个字节，读取指针后移 1
		assert.NoError(t, err)
	}
	assert.Equal(t, "        <div>Hello World</div>\n", string(data[18:49]))

	// rune 转为 []byte 方法
	// 这里认为 rune 为 utf8 编码，一个 rune 可以转为 1~4 个 bytes
	runeToBytes := func(r rune) []byte {
		buf := make([]byte, 4)       // 可接受最长 utf8 编码的 bytes
		n := utf8.EncodeRune(buf, r) // 将字符编码为 utf8 bytes
		return buf[:n]               // 返回有效长度的 []byte 切片
	}

	len := 0
	// 读取 19 个字符，并转为 byte 存入 data 切片
	for i := 0; i < 19; i++ {
		r, n, err := reader.ReadRune() // io.RuneReader 接口函数，读取一个 rune 字符
		assert.NoError(t, err)
		assert.Equal(t, 1, n)

		for _, b := range runeToBytes(r) { // 将 rune 转为 bytes
			data[49+len] = b // 存入 data 切片的后续位置
			len++
		}
	}
	assert.Equal(t, "    </body>\n</html>", string(data[49:]))

	// 查看 Reader 中剩余未读取的数据
	assert.Nil(t, reader.UnreadRune()) // io.RuneScanner 接口函数，获取尚未读取的 rune，该函数必须在 ReadRune() 之后执行，之间不能插入其它 Reader 操作
	assert.Nil(t, reader.UnreadByte()) // io.ByteScanner 接口函数，获取尚未读取的 bytes

	// 对比整个读取结果和源数据
	assert.Equal(t, `<html>
    <body>
        <div>Hello World</div>
    </body>
</html>`, string(data))
}

// 测试 byte 类型数据的读写操作
// 数据读写依赖 bytes.Buffer 类型，其实现了 io.Reader 和 io.Writer 接口，可以同时进行读写操作
// 字符串类型通过编码为 utf8 编码，进行读取和写入
// 对于 int, float, bool, rune, slice 等类型，需要借助 binary 包，转换为 byte 类型后进行读写
func TestBytesBuffer(t *testing.T) {
	// 写入 bytes.Buffer 对象
	buf := bytes.NewBuffer([]byte{}) // 产生一个初始长度为 0 的 Buffer 对象进行写入

	count, err := buf.Write([]byte(`Hello World`)) // io.Writer 接口函数，写入编码过的字符串
	assert.NoError(t, err)
	assert.Equal(t, 11, count)     // 写入 11 字节
	assert.Equal(t, 11, buf.Len()) // io.Writer 接口函数，共写入 11 字节

	binary.Write(buf, binary.BigEndian, int64(100)) // 写入 int64 类型，大端模式，binary.Write 函数用于将任意类型转为 bytes（大端或小端）后写入 io.Write 对象
	assert.Equal(t, 19, buf.Len())                  // 共写入 19 字节，增加 int64 = 8 字节

	binary.Write(buf, binary.BigEndian, true) // 写入 bool 类型，大端模式
	assert.Equal(t, 20, buf.Len())            // 共写入 20 字节，增加 bool = 1 字节

	binary.Write(buf, binary.LittleEndian, []float64{1.1, 2.2, 3.3}) // 写入 float64 切片，小端模式
	assert.Equal(t, 44, buf.Len())                                   // 共写入 44 字节，增加 float64 * 3 = 24 字节

	data := make([]byte, 4)
	count = binary.PutVarint(data, int64(100))        // 将 int64 以“可变长度”（变体）形式存入 byte 数组。变体 （varint）可以根据数值的大小变化编码长度，可以节省存储空间
	assert.Equal(t, 2, count)                         // 变体长度为 2，较 int64 原本长度（长度8）减少 6 个字节
	binary.Write(buf, binary.BigEndian, data[:count]) // 将变体写入 Buffer 对象
	assert.Equal(t, 46, buf.Len())                    // 共写入 46 字节，增加 varint = 2 字节

	count = binary.PutUvarint(data, uint64(123456))   // 将 uint64 以变体形式存入 byte 数组
	assert.Equal(t, 3, count)                         // 变体长度为 2，较 uint64 原本长度（长度8）减少 6 个字节
	binary.Write(buf, binary.BigEndian, data[:count]) // 将变体写入 Buffer 对象
	assert.Equal(t, 49, buf.Len())                    // 共写入 49 字节，增加 uvarint = 3 字节

	assert.Equal(t, 84, buf.Cap())      // Buffer 对象实际容量 84 字节
	buf.Grow(buf.Cap() - buf.Len() + 1) // Grow() 函数用来增加 Buffer 对象的容量，增加到 85 字节

	assert.Equal(t, 49, buf.Len())  // 已写入仍为 49 字节（不变）
	assert.Equal(t, 204, buf.Cap()) // Buffer 实际容量增加到 204 字节

	buf.Truncate(buf.Len()) // 截断内容，将指定长度之后的内容清除

	// 从 bytes.Buffer 进行读取
	buf = bytes.NewBuffer(buf.Bytes()) // 从写入结果产生新的 Buffer 对象。buf.Bytes() 返回一个 []byte 切片，包括 Buffer 对象中所有内容

	data = make([]byte, 11)     // 产生 11 字节的 byte 切片
	count, err = buf.Read(data) // 读取 11 字节到切片中
	assert.NoError(t, err)
	assert.Equal(t, 11, count)                   // 确保读取了 11 字节
	assert.Equal(t, "Hello World", string(data)) // 比较读取的前 11 字节和写入的前 11 字节

	var num int64
	err = binary.Read(buf, binary.BigEndian, &num) // 继续读取 8 字节，写入 int64 变量中，大端模式
	assert.NoError(t, err)
	assert.Equal(t, int64(100), num) // 比较读取的 8 字节和写入的 8 字节

	var b bool
	err = binary.Read(buf, binary.BigEndian, &b) // 继续读取 1 字节，写入 bool 变量中，大端模式
	assert.NoError(t, err)
	assert.Equal(t, true, b) // 比较读取的 1 字节和写入的 1 字节

	arr := make([]float64, 3)
	err = binary.Read(buf, binary.LittleEndian, &arr) // 继续读取 24 字节，写入 3 项的 float64 切片中
	assert.NoError(t, err)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, arr) // 比较读取的 24 字节和写入的 24 字节

	n, err := binary.ReadVarint(buf) // 读取一个 int64 变体
	assert.NoError(t, err)
	assert.Equal(t, int64(100), n) // 比较写入的变体和读取的变体

	un, err := binary.ReadUvarint(buf) // 读取 uint64 变体
	assert.NoError(t, err)
	assert.Equal(t, uint64(123456), un) // 比较写入的变体和读取的变体

	// Buffer 的拷贝
	bufNew := bytes.NewBuffer([]byte{})

	buf.WriteTo(bufNew) // 将 buf 对象内容写入到 newBuf 对象中
	assert.Equal(t, buf.Bytes(), bufNew.Bytes())

	buf.Reset()          // 清除 buf 内容，相当于 buf.Truncate(0)，即在 0 位置截断内容
	buf.ReadFrom(bufNew) // 将 bufNew 内容写入到 buf 对象中
	assert.Equal(t, buf.Bytes(), bufNew.Bytes())
}

// 文本文件的读取和写入
// os.File 对象实现了如下接口
//  io.Writer, io.WriterString, io.ReadFrom, io.Reader, io.ReaderAt, io.WriterTo
// 一般情况下，不直接使用 os.File 对象，而是通过 bufio.Writer, bufio.Reader 以及 gob.Encoder, gob.Decoder 来进行操作
func TestFileIO(t *testing.T) {
	path := "./test.dat"

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	assert.NoError(t, err)
	assert.Equal(t, "./test.dat", file.Name())

	// 函数结束后，关闭文件句柄并删除文件
	defer func() {
		err := file.Close()
		assert.NoError(t, err)

		err = os.Remove(file.Name())
		assert.NoError(t, err)
	}()

	// 获取文件长度
	fileLen := func(file *os.File) int {
		fi, err := os.Stat(file.Name()) // 获取文件属性
		assert.NoError(t, err)
		return int(fi.Size()) // 从文件属性中获取文件实际长度
	}

	count, err := file.WriteString(`Hello World`)
	assert.NoError(t, err)
	assert.Equal(t, 11, count)
	assert.Equal(t, 11, fileLen(file)) // 写入 11 字节

	binary.Write(file, binary.BigEndian, int64(100)) // 写入 int64 类型，大端模式，binary.Write 函数用于将任意类型转为 bytes（大端或小端）后写入 io.Write 对象
	assert.Equal(t, 19, fileLen(file))               // 共写入 19 字节，增加 int64 = 8 字节

	binary.Write(file, binary.BigEndian, true) // 写入 bool 类型，大端模式
	assert.Equal(t, 20, fileLen(file))         // 共写入 20 字节，增加 bool = 1 字节

	binary.Write(file, binary.LittleEndian, []float64{1.1, 2.2, 3.3}) // 写入 float64 切片，小端模式
	assert.Equal(t, 44, fileLen(file))                                // 共写入 44 字节，增加 float64 * 3 = 24 字节

	data := make([]byte, 4)
	count = binary.PutVarint(data, int64(100))         // 将 int64 以“可变长度”（变体）形式存入 byte 数组。变体 （varint）可以根据数值的大小变化编码长度，可以节省存储空间
	assert.Equal(t, 2, count)                          // 变体长度为 2，较 int64 原本长度（长度8）减少 6 个字节
	binary.Write(file, binary.BigEndian, data[:count]) // 将变体写入 Buffer 对象
	assert.Equal(t, 46, fileLen(file))                 // 共写入 46 字节，增加 varint = 2 字节

	count = binary.PutUvarint(data, uint64(123456))    // 将 uint64 以变体形式存入 byte 数组
	assert.Equal(t, 3, count)                          // 变体长度为 2，较 uint64 原本长度（长度8）减少 6 个字节
	binary.Write(file, binary.BigEndian, data[:count]) // 将变体写入 Buffer 对象
	assert.Equal(t, 49, fileLen(file))                 // 共写入 49 字节，增加 uvarint = 3 字节

	file.Truncate(int64(fileLen(file)))

	file.Close()

	// 以只读方式打开文件，文件必须存在
	file, err = os.OpenFile(path, os.O_RDONLY, 0666)
	assert.NoError(t, err)
	assert.Equal(t, "./test.dat", file.Name())

	data = make([]byte, 11)      // 产生 11 个 byte 的 bytes
	count, err = file.Read(data) // 读取 11 个字节
	assert.NoError(t, err)
	assert.Equal(t, 11, count)
	assert.Equal(t, "Hello World", string(data)) // 比较读取的前 11 字节和写入的前 11 字节

	var num int64
	err = binary.Read(file, binary.BigEndian, &num) // 读取 int64（8 字节），大端
	assert.NoError(t, err)
	assert.Equal(t, int64(100), num)

	var b bool
	err = binary.Read(file, binary.BigEndian, &b) // 读取 bool（1 字节），大端
	assert.NoError(t, err)
	assert.Equal(t, true, b)

	arr := make([]float64, 3)
	err = binary.Read(file, binary.LittleEndian, &arr) // 读取长度为 3 的 float64 切片
	assert.NoError(t, err)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, arr)

	data = make([]byte, 8)

	count, err = file.Read(data) // 读取 varint 值，由于 File 对象没有实现 io.ByteReader 接口，所以无法直接读取 varint，需要先读取到内存里，在进行处理
	assert.NoError(t, err)       // 读取 8 字节（varint 的最大长度）
	assert.Equal(t, 5, count)    // 实际读出来 5 字节，即文件到结尾仅剩 5 字节

	n, count := binary.Varint(data) // 从 bytes 中获取 varint 值，返回 varint 的实际长度 2 字节
	assert.Equal(t, int64(100), n)
	assert.Equal(t, 2, count)

	file.Seek(int64(-3), io.SeekCurrent) // 文件指针回退到读取 varint 的实际位置，即 -(5 - 2)

	count, err = file.Read(data) // 读取 uvarint 值
	assert.NoError(t, err)
	assert.Equal(t, 3, count) // 读取 8 字节，实际读取 3 字节（仅剩 3 字节）

	un, count := binary.Uvarint(data) // 从 bytes 中获取 uvarint 值，返回 uvarint 实际长度 3 字节
	assert.Equal(t, uint64(123456), un)
	assert.Equal(t, 3, count)
}

// 使用 GoB（Go Binary）操作 io
// gob.Encoder 对象和 gob.Decoder 对象可以对值（数值、字符串、切片等）进行编解码，编码的结果直接写入 io.Writer 对象，解码则是直接通过 io.Reader 进行
// gob 方式可以极大的简化各类数据写入和读取操作
// 另外，gob.Encoder 的 EncoderValue 以及 gob.Decoder 的 DecoderValue 可以对 reflect.Value 对象进行操作，通过反射处理 io
func TestGobDatabase(t *testing.T) {
	file, err := os.Create("./gob.data") // os.Create(name) 函数是 os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666) 函数的简写，打开一个读写文件
	assert.NoError(t, err)

	// 函数结束后，关闭并删除文件
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	s := "Hello, World!"

	// 通过编码器将各类数据写入 io.Writer 对象
	enc := gob.NewEncoder(file) // 创建 编码器 对象，参数为一个 io.Writer 对象

	err = enc.Encode(len(s)) // 写入 int 类型数值
	assert.NoError(t, err)

	err = enc.Encode(s) // 写入字符串类型数据
	assert.NoError(t, err)

	pu := user.New(1, "Alvin", "alvin@fake.com", []string{"13999912345", "13000056789"}) // 初始化 user.User 对象并返回指针

	err = enc.Encode(pu) // 将结构体变量进行编码
	assert.NoError(t, err)

	file.Close()

	file, err = os.Open("./gob.data") // os.Open(name) 函数是 os.OpenFile(name, os.O_RDONLY, 0) 函数的简写，打开一个只读文件
	assert.NoError(t, err)
	dec := gob.NewDecoder(file) // 创建 解码器 对象，参数为一个 io.Reader 对象

	var n int
	err = dec.Decode(&n) // 解码一个整数
	assert.NoError(t, err)
	assert.Equal(t, len(s), n)

	var rs string
	err = dec.Decode(&rs) // 解码一个字符串
	assert.NoError(t, err)
	assert.Equal(t, s, rs)

	var ru user.User
	err = dec.Decode(&ru) // 解码结构体对象
	assert.NoError(t, err)
	assert.Equal(t, *pu, ru)
}

// 通过缓冲进行 IO 操作
// bufio.Writer 实现了 io.Writer, io.ByteWriter, io.StringWriter 以及 io.RuneWriter 接口，用于包装另一个 io.Writer 对象，为其增加缓存支持
// bufio.Reader 实现了 io.Reader, io.ByteReader, io.StringWriter 以及 io.RuneWriter 接口，用于包装另一个 io.Reader 对象，为其增加缓存支持
// bufio.ReadWriter 同时包装了一个 io.Writer 和一个 io.Reader 对象，实现了同时读写的操作
func TestBufferedIO(t *testing.T) {
	fileName := "buffered.dat"

	// 利用缓冲流写入文件
	file, err := os.Create(fileName) // 打开用于写入的文件
	assert.NoError(t, err)

	// 函数结束后关闭并删除文件
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	bufW := bufio.NewWriterSize(file, 256) // 创建缓存大小 256 字节的 bufio.Writer 对象，用于对文件进行写操作. 另外，bufio.NewWriter(file) 用于创建一个缓冲区大小为 4096 的 bufio.Writer 对象
	assert.Equal(t, 256, bufW.Size())      // 已缓存数据大小
	assert.Equal(t, 256, bufW.Available()) // 缓存总大小
	assert.Equal(t, 0, bufW.Buffered())    // 已缓存数据大小

	bufW.WriteString("Hello\n")            // 写入字符串
	assert.Equal(t, 250, bufW.Available()) // 缓存剩余大小
	assert.Equal(t, 6, bufW.Buffered())    // 写入数据后，已缓存数据大小

	pu := user.New(1, "Alvin", "alvin@fake.com", []string{"13999912345", "13000056789"})

	enc := gob.NewEncoder(bufW)            // 在 bufio.Writer 对象的基础上，包装一个 gob.Encoder 对象
	enc.Encode(pu)                         // 将对象编码后写入 bufio.Writer 对象
	assert.Equal(t, 116, bufW.Available()) // 缓存剩余大小
	assert.Equal(t, 140, bufW.Buffered())  // 写入数据后，已缓存数据大小，共 140 字节

	err = bufW.Flush() // 将缓存的内容写入文件
	assert.NoError(t, err)
	assert.Equal(t, 256, bufW.Available()) // 缓存剩余大小
	assert.Equal(t, 0, bufW.Buffered())    // 缓存数据已经全部写入文件

	file.Close() // 关闭文件

	// 利用缓冲流读文件
	file, err = os.Open(fileName) // 打开文件用于读
	assert.NoError(t, err)

	bufR := bufio.NewReaderSize(file, 256) // 创建缓存大小 256 字节的 bufio.Reader 对象，用于对文件进行读操作. 另外，bufio.NewReader(file) 用于创建一个缓冲区大小为 4096 的 bufio.Reader 对象
	assert.Equal(t, 0, bufR.Buffered())    // 已缓存数据 0 字节（尚未开始读取）

	line, prefix, err := bufR.ReadLine() // 读取一行字符串（以 \n 结尾的一行数据）
	assert.NoError(t, err)
	assert.Equal(t, "Hello", string(line)) // 读取 1 行的内容
	assert.False(t, prefix)                // prefix 返回 false，表示以读完一行，否则表示一行尚未读完，需要再次调用 ReadLine 函数，直到读完一行
	assert.Equal(t, 134, bufR.Buffered())  // 缓冲区使用 134 字节。具体情况为：当第一次读取时，先读入缓存 140 字节，读走 1 行 6 字节后，剩余 134 字节

	u := user.User{}
	dec := gob.NewDecoder(bufR) // 在 bufio.Reader 对象基础上创建 gob.Decoder 对象

	err = dec.Decode(&u) // 从文件中解码一个 user.User 对象
	assert.NoError(t, err)
	assert.Equal(t, *pu, u)
	assert.Equal(t, 0, bufR.Buffered()) // 缓冲区已读完

	// 测试同时读写
	fw, err := os.Create(fileName) // 创建用于写的 os.File 对象
	assert.NoError(t, err)

	fr, err := os.Open(fileName) // 创建用于读的 os.File 对象
	assert.NoError(t, err)

	// 函数结束后关闭并删除文件
	defer func() {
		fw.Close()
		fr.Close()
		os.Remove(fileName)
	}()

	bufRW := bufio.NewReadWriter(bufio.NewReader(fr), bufio.NewWriter(fw)) // 通过一对 bufio.Reader 和 bufio.Writer 对象，创建同时进行读写的 bufio.ReaderWriter 对象

	bufRW.WriteString("Hello\n") // 写入字符串

	pu = user.New(1, "Alvin", "alvin@fake.com", []string{"13999912345", "13000056789"}) // 通过 gob.Encoder 写入结构体
	enc = gob.NewEncoder(bufRW)
	enc.Encode(pu)

	bufRW.Flush() // 将缓冲区内容写入文件

	line, prefix, err = bufRW.ReadLine() // 读取一行字符串
	assert.NoError(t, err)
	assert.Equal(t, "Hello", string(line))
	assert.False(t, prefix)

	u = user.User{}

	dec = gob.NewDecoder(bufRW) // 通过 gob.Decoder 读取结构体
	dec.Decode(&u)
	assert.Equal(t, *pu, u)
}
