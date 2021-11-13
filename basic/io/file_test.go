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

func TestFileStat(t *testing.T) {
	defer os.RemoveAll(`d`)

	err := os.Mkdir(`d`, 0755)
	assert.NoError(t, err)

	stat, err := os.Stat(`d`)
	assert.NoError(t, err)
	assert.Equal(t, `d`, stat.Name())
	assert.True(t, stat.IsDir())
	assert.Equal(t, os.FileMode(020000000755), stat.Mode())

	file, err := os.Create(`d/e.txt`)
	assert.NoError(t, err)

	stat, err = os.Stat(`d/e.txt`)
	assert.NoError(t, err)
	assert.Equal(t, `e.txt`, stat.Name())
	assert.False(t, stat.IsDir())
	assert.Equal(t, os.FileMode(0644), stat.Mode())

	stat2, err := file.Stat()
	assert.NoError(t, err)
	assert.Equal(t, stat, stat2)
}
