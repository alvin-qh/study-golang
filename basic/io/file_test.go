package io

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 获取当前路径
func TestCurrentPath(t *testing.T) {
	// 获取当前路径
	p1, err := os.Getwd()
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(p1, `/basic/io`))

	p2, err := filepath.Abs(`.`)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(p2, `/basic/io`))

	assert.Equal(t, p2, p1)
}

// 获取可执行文件的路径
// 获取的位置可以是指定位置或者 $PATH 变量定义的路径
func TestLookupExecutableFile(t *testing.T) {
	// 获取 bash 文件的位置，该位置定义在 $PATH 环境变量中
	path, err := exec.LookPath("bash")
	assert.NoError(t, err)
	assert.NotEmpty(t, "/usr/bin/bash", path)

	defer os.Remove(`./temp`)

	// 创建一个可执行文件并查找它，这里需要明确的绝对或相对路径，否则将在 $PATH 变量中查找
	_, err = os.OpenFile(`./temp`, os.O_CREATE|os.O_RDWR, 0777)
	assert.NoError(t, err)

	path, err = exec.LookPath(`./temp`) // 查找可执行文件
	assert.NoError(t, err)
	assert.Equal(t, "./temp", path)
}

// 获取文件（或路径）的属性
// 共有两种方法：1. 通过 os.Stat(`<file or path>`) 函数；2. 打开 文件或路径 的 os.File 对象，通过 os.File::Stat() 函数获取
func TestFileStat(t *testing.T) {
	defer os.RemoveAll(`d`)

	// 创建一个目录
	err := os.Mkdir(`d`, 0755)
	assert.NoError(t, err)

	// 方法 1： 通过 os.Stat 函数获取文件属性
	// 获取目录 `d` 的属性对象
	stat, err := os.Stat(`d`)
	assert.NoError(t, err)
	assert.Equal(t, `d`, stat.Name())                       // 获取路径名
	assert.True(t, stat.IsDir())                            // 获取是一个路径
	assert.Equal(t, os.FileMode(020000000755), stat.Mode()) // 获取路径的访问权限

	// 打开目录 `d` 并获取其 os.File 对象
	file, err := os.Open(`d`)
	assert.NoError(t, err)
	assert.Equal(t, `d`, file.Name()) // 获取路径名称

	defer file.Close()

	stat2, err := file.Stat() // 通过 os.File 对象获取文件属性
	assert.NoError(t, err)
	assert.Equal(t, stat, stat2) // 两种方式获取的文件属性完全一致

	file.Close()

	// 方法 2： 通过 os.File 对象的 Stat 函数获取文件属性
	// 创建一个文件
	file, err = os.Create(`d/e.txt`)
	assert.NoError(t, err)
	assert.Equal(t, `d/e.txt`, file.Name()) // 获取文件名称

	// 获取文件 `d/e.txt` 的属性对象
	stat, err = os.Stat(`d/e.txt`)
	assert.NoError(t, err)
	assert.Equal(t, `e.txt`, stat.Name())           // 获取文件名
	assert.False(t, stat.IsDir())                   // 判断不是一个路径
	assert.Equal(t, os.FileMode(0644), stat.Mode()) // 获取文件的访问权限

	// 通过 os.File 对象获取文件属性对象
	stat2, err = file.Stat()
	assert.NoError(t, err)
	assert.Equal(t, stat, stat2) // 两种方式获取的文件属性对象一致
}

// 当打开的 os.File 对象表示一个 路径 时，可以读取其包括的所有文件（或子目录）的信息
func TestReadDir(t *testing.T) {
	// 打开路径为 os.File 对象
	dir, err := os.Open(`.`)
	assert.NoError(t, err)

	// 读取路径下的信息，返回所有文件（包括路径）的 os.Stat 对象
	// 参数 0 表示不限制返回结果的数量，否则按所给数量返回结果
	infos, err := dir.Readdir(0)
	assert.NoError(t, err)

	assert.Len(t, infos, 6)

	expected := []string{"io_test.go", "file_test.go", "json_test.go", "user", "path_test.go", "xml_test.go"}
	for n, info := range infos {
		assert.Equal(t, expected[n], info.Name())
		if strings.HasSuffix(info.Name(), ".go") {
			assert.False(t, info.IsDir())
		} else {
			assert.True(t, info.IsDir())
		}
	}
}

// 当打开的 os.File 对象表示一个 路径 时，可以读取其包括的所有文件（或子目录）的名称
func TestReadDirnames(t *testing.T) {
	// 打开路径为 os.File 对象
	dir, err := os.Open(`.`)
	assert.NoError(t, err)

	// 读取路径下的信息，返回所有文件（包括路径）的 os.Stat 对象
	// 参数 0 表示不限制返回结果的数量，否则按所给数量返回结果
	names, err := dir.Readdirnames(0)
	assert.NoError(t, err)

	assert.Len(t, names, 6)

	expected := []string{"io_test.go", "file_test.go", "json_test.go", "user", "path_test.go", "xml_test.go"}
	for n, name := range names {
		assert.Equal(t, expected[n], name)
	}
}

// 利用管道进行数据传输
// 创建管道即得到一对 io.Reader 和 io.Writer 对象，对 io.Writer 对象进行写操作，则可随后从 io.Reader 读出所写的内容
func TestPipe(t *testing.T) {
	// 创建一个管道，得到一对 io.Reader 和 io.Writer 对象
	r, w, err := os.Pipe()
	defer func() { // 在函数结束后关闭 io.Reader 和 io.Writer 对象
		r.Close()
		w.Close()
	}()

	assert.NoError(t, err)
	assert.Equal(t, "|0", r.Name()) // 管道创建的文件也具有名称
	assert.Equal(t, "|1", w.Name())

	rs, err := r.Stat() // 获取 io.Reader 对象的 文件属性
	assert.NoError(t, err)
	assert.False(t, rs.IsDir()) // 判断为文件对象

	ws, err := r.Stat() // 获取 io.Writer 对象的 文件属性
	assert.NoError(t, err)
	assert.False(t, ws.IsDir()) // 判断为文件对象

	br := bufio.NewReader(r) // 利用 bufio 给 io.Reader 和 io.Writer 增加缓存
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

// 截取文件，即将文件截断为指定长度，多余的部分将被丢弃
// go 提供了两种截取文件的方法：1. 通过 os.File 对象的 Truncate 函数；2. 通过 os.Truncate 函数。前者需要打开文件，得到一个 os.File 对象
func TestTuncateAttributes(t *testing.T) {
	defer os.Remove(`d.txt`)

	fileLength := func(file *os.File) int {
		if stat, err := file.Stat(); err == nil {
			return int(stat.Size())
		}
		return 0
	}

	// 方法 1: 通过 os.File 对象的 Truncate 函数截断文件
	// 创建文件
	file, err := os.Create(`d.txt`)
	assert.NoError(t, err)

	defer file.Close()

	// 写入10个字符
	count, err := file.WriteString("1234567890")
	assert.NoError(t, err)
	assert.Equal(t, 10, count)
	assert.Equal(t, 10, fileLength(file)) // 此时文件长度为 10 字节

	err = file.Truncate(5) // 截断文件，为原始长度的一半
	assert.NoError(t, err)
	assert.Equal(t, 5, fileLength(file)) // 此时文件长度为 5

	_, err = file.Seek(0, io.SeekStart) // 指针移动到文件开头
	assert.NoError(t, err)

	// 读取全部文件内容，只剩下写入内容的一半
	s, err := io.ReadAll(file)
	assert.NoError(t, err)
	assert.Equal(t, "12345", string(s))

	file.Close()

	// 方法 2: 通过 os.Truncate 函数对文件直接进行截断
	os.Truncate(`d.txt`, 3) // 将文件进一步截取到剩余 3 字节

	file, err = os.Open(`d.txt`) // 打开文件
	assert.NoError(t, err)
	assert.Equal(t, 3, fileLength(file)) // 文件长度只剩余 3 字节

	s, err = io.ReadAll(file) // 读取文件内容，只剩余 3 个字符
	assert.NoError(t, err)
	assert.Equal(t, "123", string(s))
}
