package file

import (
	"bufio"
	"encoding/gob"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 获取当前路径
//
// Go 语言提供了两套文件路径处理包: `path` 和 `filepath`, 两者提供了类似的函数库,
// 但后者具备跨平台路径处理能力
func TestCurrentPath(t *testing.T) {
	curPath := filepath.Join("basic", "io", "file")

	// 获取当前路径
	p1, err := os.Getwd()
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(p1, curPath))

	p2, err := filepath.Abs(".")
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(p2, curPath))

	assert.Equal(t, p2, p1)
}

// 获取可执行文件的路径
//
// 获取的位置可以是指定位置或者 $PATH 变量定义的路径
func TestLookupExecutableFile(t *testing.T) {
	var exeFile, target string
	if runtime.GOOS == "windows" {
		exeFile = "cmd.exe"
		target = "c:\\windows\\system32\\" + exeFile
	} else {
		exeFile = "bash"
		target = "/bin/" + exeFile
	}

	// 获取 bash 文件的位置, 该位置定义在 $PATH 环境变量中
	path, err := exec.LookPath(exeFile)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(strings.ToLower(path), target))

	// 创建一个可执行文件并查找它, 这里需要明确的绝对或相对路径, 否则将在 $PATH 变量中查找
	f, err := os.OpenFile("./temp.exe", os.O_CREATE|os.O_RDWR, 0777)
	assert.NoError(t, err)

	defer func() {
		f.Close()
		os.Remove("./temp.exe")
	}()

	path, err = exec.LookPath("./temp.exe") // 查找可执行文件
	assert.NoError(t, err)
	assert.Equal(t, "./temp.exe", path)
}

// 获取文件 (或路径) 的属性
//
// 共有两种方法:
//  1. 通过 `os.Stat` 函数以及文件路径获取;
//  2. 获取 `os.File` 实例, 通过其 `Stat` 方法获取;
func TestFileStat(t *testing.T) {
	// 创建一个目录
	err := os.Mkdir("d", 0755)
	assert.NoError(t, err)

	defer os.RemoveAll("d")

	// 方法 1:  通过 os.Stat 函数获取文件属性
	// 获取目录 `d` 的属性对象
	stat1, err := os.Stat("d")
	assert.NoError(t, err)
	assert.Equal(t, "d", stat1.Name()) // 获取路径名
	assert.True(t, stat1.IsDir())      // 获取是一个路径
	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(020000000755), stat1.Mode()) // 获取路径的访问权限
	}

	// 打开目录 `d` 并获取其 os.File 对象
	file, err := os.Open("d")
	assert.NoError(t, err)
	assert.Equal(t, "d", file.Name()) // 获取路径名称

	stat2, err := file.Stat() // 通过 os.File 对象获取文件属性
	assert.NoError(t, err)
	assert.True(t, CompareFileInfo(stat1, stat2)) // 两种方式获取的文件属性完全一致

	file.Close()

	// 方法 2:  通过 os.File 对象的 Stat 函数获取文件属性
	// 创建一个文件
	file, err = os.Create("d/e.txt")
	assert.NoError(t, err)
	assert.Equal(t, "d/e.txt", file.Name()) // 获取文件名称

	// 获取文件 "d/e.txt" 的属性对象
	stat1, err = os.Stat("d/e.txt")
	assert.NoError(t, err)
	assert.Equal(t, "e.txt", stat1.Name()) // 获取文件名
	assert.False(t, stat1.IsDir())         // 判断不是一个路径
	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(0644), stat1.Mode()) // 获取文件的访问权限
	}

	// 通过 os.File 对象获取文件属性对象
	stat2, err = file.Stat()
	assert.NoError(t, err)
	assert.True(t, CompareFileInfo(stat1, stat2)) // 两种方式获取的文件属性完全一致

	file.Close()
}

// 获取一个目录下所有内容的信息
//
// 当 `os.File` 实例表示路径时, 可以通过其 `Readdir` 方法读取其包括的所有文件 (或子目录) 的信息
func TestReadDir(t *testing.T) {
	// 打开路径为 os.File 对象
	dir, err := os.Open(`.`)
	assert.NoError(t, err)

	// 读取路径下的信息, 返回所有文件 (包括路径) 的 `os.Stat` 对象
	// 参数 `0` 表示不限制返回结果的数量, 否则按所给数量返回结果
	infos, err := dir.Readdir(0)
	assert.NoError(t, err)

	expected := map[string]struct{}{
		"file_test.go": {},
		"file.go":      {},
	}

	// 遍历目录下所有文件信息
	for _, info := range infos {
		_, exist := expected[info.Name()]
		assert.True(t, exist)

		delete(expected, info.Name())

		if strings.HasSuffix(info.Name(), ".go") {
			assert.False(t, info.IsDir())
		} else {
			assert.True(t, info.IsDir())
		}
	}
	assert.Len(t, expected, 0)
}

// 获取一个目录下所有内容的名称
//
// 当打开的 `os.File` 对象表示一个路径时, 可以读取其包括的所有文件 (或子目录) 的名称
func TestReadDirnames(t *testing.T) {
	// 打开路径为 os.File 对象
	dir, err := os.Open(`.`)
	assert.NoError(t, err)

	// 读取路径下的信息, 返回所有文件 (包括路径) 的 os.Stat 对象
	// 参数 0 表示不限制返回结果的数量, 否则按所给数量返回结果
	names, err := dir.Readdirnames(0)
	assert.NoError(t, err)
	sort.Strings(names)

	expected := []string{"file_test.go", "file.go"}
	sort.Strings(expected)

	assert.Equal(t, expected, names)
}

// 测试文件截取
//
// 文件截取即将文件截断为指定长度, 多余的部分将被丢弃
//
// Go 提供了两种截取文件的方法:
//  1. 通过 `os.File` 实例的 `Truncate` 方法, 需要打开文件, 得到一个 `os.File` 对象;
//  2. 通过 `os.Truncate` 函数. 无需打开文件;
func TestFileTruncate(t *testing.T) {
	// 方法 1: 通过 `os.File` 对象的 `Truncate` 函数截断文件

	// 创建文件
	file, err := os.Create("d.txt")
	assert.NoError(t, err)

	defer os.Remove("d.txt")

	// 写入10个字符
	count, err := file.WriteString("1234567890")
	assert.NoError(t, err)
	assert.Equal(t, 10, count)
	assert.Equal(t, 10, FileLength(file)) // 此时文件长度为 10 字节

	// 截断文件, 为原始长度的一半
	err = file.Truncate(5)
	assert.NoError(t, err)
	assert.Equal(t, 5, FileLength(file)) // 此时文件长度为 5

	// 指针移动到文件开头
	_, err = file.Seek(0, io.SeekStart)
	assert.NoError(t, err)

	// 读取全部文件内容, 只剩下写入内容的一半
	s, err := io.ReadAll(file)
	assert.NoError(t, err)
	assert.Equal(t, "12345", string(s))

	file.Close()

	// 方法 2: 通过 `os.Truncate` 函数对文件直接进行截断

	// 将文件进一步截取到剩余 3 字节
	os.Truncate("d.txt", 3)

	// 打开文件
	file, err = os.Open("d.txt")
	assert.NoError(t, err)
	assert.Equal(t, 3, FileLength(file)) // 文件长度只剩余 3 字节

	// 读取文件内容, 只剩余 3 个字符
	s, err = io.ReadAll(file)
	assert.NoError(t, err)
	assert.Equal(t, "123", string(s))

	file.Close()
}

const (
	FILE_NAME = "test.dat"
)

// 测试 IO 的结构体
type user struct {
	Id    int64
	Name  string
	Email string
	Phone []string
}

// 文件读写测试
//
// `os.File` 类型实现了 `io.Reader` 和 `io.Writer` 接口, 可以直接对文件进行读写
//
// 除了可以通过 `os.Create` 函数创建 `os.File` 实例外, `os.OpenFile` 函数也可以创建新文件或打开一个现有文件,
// 返回 `os.File` 实例, 且如果是新建文件, 还可以指定文件的访问权限
func TestFileIO(t *testing.T) {
	// 测试文件写

	// 打开 (或截断) 一个文件
	file, err := os.OpenFile(FILE_NAME, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	assert.NoError(t, err)
	assert.Equal(t, FILE_NAME, file.Name())

	defer os.Remove(FILE_NAME)

	// 在文件 0 位置写入字符串, 共 11 字节. 写指针指向 11 位置
	n, err := file.WriteString("Hello World")
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, 11, FileLength(file))

	data := []byte{1, 2, 3, 4, 5}

	// 指定在文件 16 位置写入 5 个字节, 写指针不变, 仍指向 11 位置, 文件长度变为 21 字节
	n, err = file.WriteAt(data, 16)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, 21, FileLength(file))
	assert.Equal(t, int64(11), GetFileCursor(file)) // 获取当前文件指针, 仍为 11

	// 写指针从文件末尾移动 -10 个偏移, 指向 11 位置
	cur, err := file.Seek(-10, io.SeekEnd)
	assert.NoError(t, err)
	assert.Equal(t, int64(11), cur)

	// 从 11 位置写入 5 字节, 写指针指向 16 位置
	n, err = file.Write([]byte{5, 4, 3, 2, 1})
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, 21, FileLength(file))

	file.Close()

	// 测试文件读

	// 打开一个只读文件
	file, err = os.OpenFile(FILE_NAME, os.O_RDONLY, 0)
	assert.NoError(t, err)
	assert.Equal(t, FILE_NAME, file.Name())

	// 读取文件的 bytes 缓冲区
	data = make([]byte, FileLength(file))

	n, err = file.Read(data[0:11]) // 从 0 位置读取 11 字节, 读指针移动到 11
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, "Hello World", string(data[0:11])) // 读出字符串内容

	// 在文件 16 位置读取 5 字节, 读指针不变
	n, err = file.ReadAt(data[16:21], 16)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, data[16:21])
	assert.Equal(t, int64(11), GetFileCursor(file)) // 读指针仍在 11

	// 从文件开始位置, 将读指针移动到 11 位置
	cur, err = file.Seek(11, io.SeekStart)
	assert.NoError(t, err)
	assert.Equal(t, int64(11), cur)

    // 在 11 位置读取 5 字节, 读指针移动到 16
	n, err = file.Read(data[11:16])
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte{5, 4, 3, 2, 1}, data[11:16])

	file.Close()
}

// 通过缓冲进行 IO 操作
//
// `bufio.Writer` 类型用于包装一个 `io.Writer` 实例, 为其增加缓冲支持;
// `bufio.Writer` 类型同时实现了 `io.Writer`, `io.ByteWriter`, `io.StringWriter` 以及 `io.RuneWriter` 接口, 可以将各类数据写入目标;
//
// `bufio.Reader` 类型用于包装一个 `io.Reader` 实例, 为其增加缓冲支持;
// `bufio.Reader` 类型同时实现了 `io.Reader`, `io.ByteReader`, `io.StringReader` 以及 `io.RuneReader` 接口, 可从数据源读取各类数据;
//
// `bufio.ReadWriter` 类型用于同时包装一个 `io.Writer` 和一个 `io.Reader` 实例, 具备读写能力并增加了缓冲支持
func TestBufferedFileIO(t *testing.T) {
	// 利用缓冲流写入文件
	file, err := os.Create(FILE_NAME) // 打开用于写入的文件
	assert.NoError(t, err)

	defer os.Remove(FILE_NAME)

	// 创建缓存大小 256 字节的 bufio.Writer 对象, 用于对文件进行写操作. 另外, bufio.NewWriter(file) 用于创建一个缓冲区大小为 4096 的 bufio.Writer 实例
	bufW := bufio.NewWriterSize(file, 256)
	assert.Equal(t, 256, bufW.Size())      // 已缓存数据大小
	assert.Equal(t, 256, bufW.Available()) // 缓存总大小
	assert.Equal(t, 0, bufW.Buffered())    // 已缓存数据大小

	bufW.WriteString("Hello\n")            // 写入字符串
	assert.Equal(t, 250, bufW.Available()) // 缓存剩余大小
	assert.Equal(t, 6, bufW.Buffered())    // 写入数据后, 已缓存数据大小

	pu := &user{
		1, "Alvin", "alvin@fake.com", []string{"13999912345", "13000056789"},
	}
	// 在 bufio.Writer 对象的基础上, 包装一个 gob.Encoder 对象
	enc := gob.NewEncoder(bufW)
	// 将对象编码后写入 bufio.Writer 对象
	enc.Encode(pu)
	assert.Equal(t, 117, bufW.Available()) // 缓存剩余大小
	assert.Equal(t, 139, bufW.Buffered())  // 写入数据后, 已缓存数据大小, 共 140 字节

	err = bufW.Flush() // 将缓存的内容写入文件
	assert.NoError(t, err)
	assert.Equal(t, 256, bufW.Available()) // 缓存剩余大小
	assert.Equal(t, 0, bufW.Buffered())    // 缓存数据已经全部写入文件

	file.Close()

	// 利用缓冲流读文件
	// 打开文件用于读
	file, err = os.Open(FILE_NAME)
	assert.NoError(t, err)

	// 创建缓存大小 256 字节的 bufio.Reader 对象, 用于对文件进行读操作. 另外, bufio.NewReader(file) 用于创建一个缓冲区大小为 4096 的 bufio.Reader 对象
	bufR := bufio.NewReaderSize(file, 256)
	assert.Equal(t, 0, bufR.Buffered()) // 已缓存数据 0 字节 (尚未开始读取)

	// 读取一行字符串 (以 \n 结尾的一行数据)
	line, prefix, err := bufR.ReadLine()
	assert.NoError(t, err)
	assert.Equal(t, "Hello", string(line)) // 读取 1 行的内容
	assert.False(t, prefix)                // prefix 返回 false, 表示以读完一行, 否则表示一行尚未读完, 需要再次调用 ReadLine 函数, 直到读完一行
	assert.Equal(t, 133, bufR.Buffered())  // 缓冲区使用 134 字节. 具体情况为: 当第一次读取时, 先读入缓存 140 字节, 读走 1 行 6 字节后, 剩余 134 字节

	u := user{}
	dec := gob.NewDecoder(bufR) // 在 bufio.Reader 对象基础上创建 gob.Decoder 对象

	// 从文件中解码一个 user.User 对象
	err = dec.Decode(&u)
	assert.NoError(t, err)
	assert.Equal(t, *pu, u)
	assert.Equal(t, 0, bufR.Buffered()) // 缓冲区已读完

	file.Close()

	// 测试同时读写
	// 创建用于写的 os.File 对象
	fw, err := os.Create(FILE_NAME)
	assert.NoError(t, err)

	// 创建用于读的 os.File 对象
	fr, err := os.Open(FILE_NAME)
	assert.NoError(t, err)

	// 通过一对 bufio.Reader 和 bufio.Writer 对象, 创建同时进行读写的 bufio.ReaderWriter 对象
	bufRW := bufio.NewReadWriter(
		bufio.NewReader(fr),
		bufio.NewWriter(fw),
	)
	bufRW.WriteString("Hello\n") // 写入字符串

	// 通过 gob.Encoder 写入结构体
	pu = &user{
		1, "Alvin", "alvin@fake.com", []string{"13999912345", "13000056789"},
	}
	enc = gob.NewEncoder(bufRW)
	enc.Encode(pu)

	// 将缓冲区内容写入文件
	bufRW.Flush()

	// 读取一行字符串
	line, prefix, err = bufRW.ReadLine()
	assert.NoError(t, err)
	assert.Equal(t, "Hello", string(line))
	assert.False(t, prefix)

	u = user{}

	// 通过 gob.Decoder 读取结构体
	dec = gob.NewDecoder(bufRW)
	dec.Decode(&u)
	assert.Equal(t, *pu, u)

	fw.Close()
	fr.Close()
}

// 利用管道进行数据传输
//
// 创建管道即得到一对 `io.Reader` 和 `io.Writer` 对象, 对 `io.Writer` 对象进行写操作, 则可随后从 `io.Reader` 读出所写的内容
func TestPipe(t *testing.T) {
	// 创建一个管道, 得到一对 `io.Reader` 和 `io.Writer` 对象
	r, w, err := os.Pipe()

	assert.NoError(t, err)
	assert.Equal(t, "|0", r.Name()) // 管道创建的文件也具有名称
	assert.Equal(t, "|1", w.Name())

	// 在函数结束后关闭 `io.Reader` 和 `io.Writer` 对象
	defer func() {
		r.Close()
		w.Close()
	}()

	// 获取 io.Reader 对象的 文件属性
	rs, err := r.Stat()
	assert.NoError(t, err)
	assert.False(t, rs.IsDir()) // 判断为文件对象

	// 获取 io.Writer 对象的 文件属性
	ws, err := r.Stat()
	assert.NoError(t, err)
	assert.False(t, ws.IsDir()) // 判断为文件对象

	// 利用 bufio 给 io.Reader 和 io.Writer 增加缓存
	br := bufio.NewReader(r)
	bw := bufio.NewWriter(w)

	// 通过 io.Writer 对象对管道进行写操作
	bw.WriteString("Hello world!\n")
	bw.Flush()

	// 通过 io.Reader 对象从管道进行读操作
	s, prefix, err := br.ReadLine()
	assert.NoError(t, err)
	assert.False(t, prefix)
	assert.Equal(t, "Hello world!", string(s))
}
