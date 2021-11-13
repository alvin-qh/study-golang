package io

import (
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
// 共有两种方法：1. 通过 os.Stat(`<file or path>`) 函数，或者打开 文件或路径 的 os.File 对象，通过 os.File::Stat() 函数获取
func TestFileStat(t *testing.T) {
	defer os.RemoveAll(`d`)

	// 创建一个目录
	err := os.Mkdir(`d`, 0755)
	assert.NoError(t, err)

	// 获取目录 `d` 的属性对象
	stat, err := os.Stat(`d`)
	assert.NoError(t, err)
	assert.Equal(t, `d`, stat.Name())                       // 获取路径名
	assert.True(t, stat.IsDir())                            // 获取是一个路径
	assert.Equal(t, os.FileMode(020000000755), stat.Mode()) // 获取路径的访问权限

	// 打开目录 `d` 并获取其 os.File 对象
	file, err := os.Open(`d`)
	assert.NoError(t, err)

	stat2, err := file.Stat() // 通过 os.File 对象获取文件属性
	assert.NoError(t, err)
	assert.Equal(t, stat, stat2) // 两种方式获取的文件属性完全一致

	// 创建一个文件
	file, err = os.Create(`d/e.txt`)
	assert.NoError(t, err)

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

func TestReadDir(t *testing.T) {
	// 打开路径为 os.File 对象
	dir, err := os.Open(`.`)
	assert.NoError(t, err)

	// 读取路径下的信息，返回所有文件（包括路径）的 os.Stat 对象
	infos, err := dir.Readdir(0)
	assert.NoError(t, err)

	assert.Len(t, infos, 5)

	expected := []string{"io_test.go", "file_test.go", "json_test.go", "user", "path_test.go"}
	for n, info := range infos {
		assert.Equal(t, expected[n], info.Name())
		if strings.HasSuffix(info.Name(), ".go") {
			assert.False(t, info.IsDir())
		} else {
			assert.True(t, info.IsDir())
		}
	}
}
