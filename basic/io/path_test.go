package io

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 获取路径中的各个部分
func TestGetPathPart(t *testing.T) {
	srcPath := `a/b/c/d.xml`

	// 获取路径中的 "目录" 部分
	// 即获取路径 "最后一个分隔符之前" 的部分，例如 "a/b/c" 的结果为 "a/b"
	targetPath := filepath.Dir(srcPath)
	assert.Equal(t, `a/b/c`, targetPath)

	// 获取路径中的 "文件" 部分
	// 即获取路径 "最后一个分隔符之后" 的部分，例如 "a/b/c" 的结果为 "c"
	targetPath = filepath.Base(srcPath)
	assert.Equal(t, `d.xml`, targetPath)

	// 将路径分割为 "路径" 部分和 "文件" 部分，相当于同时调用 Dir 函数和 Base 函数
	dir, file := filepath.Split(srcPath)
	assert.Equal(t, `a/b/c/`, dir)
	assert.Equal(t, `d.xml`, file)

	// 获取路径中文件的扩展名（如果存在）
	ext := filepath.Ext(srcPath)
	assert.Equal(t, `.xml`, ext)
}

// 计算相对路径
// 即假设当前所在位置为 第一个路径，则计算结果为 到达第二个路径的 相对路径
// 计算时，两个路径必须同为 绝对路径（即以 `/` 开始）或 相对路径
func TestRelativePaths(t *testing.T) {
	// `/a` 相对于 `/` 的路径为 `a`
	s, err := filepath.Rel(`/`, `/a`)
	assert.NoError(t, err)
	assert.Equal(t, `a`, s)

	// `/a` 相对于 `/a` 的路径为 `.`
	s, err = filepath.Rel(`/a`, `/a`)
	assert.NoError(t, err)
	assert.Equal(t, `.`, s)

	// `/b` 相对于 `/a` 的路径为 `../b`
	s, err = filepath.Rel(`/a`, `/b`)
	assert.NoError(t, err)
	assert.Equal(t, `../b`, s)

	// `/a/b` 相对于 `/a` 的路径为 `b`
	s, err = filepath.Rel(`/a`, `/a/b`)
	assert.NoError(t, err)
	assert.Equal(t, `b`, s)

	// `/a/b/c` 相对于 `/a` 的路径为 `b/c`
	s, err = filepath.Rel(`/a`, `/a/b/c`)
	assert.NoError(t, err)
	assert.Equal(t, `b/c`, s)

	// 错误：第一个参数为绝对路径，第二个参数为相对路径
	_, err = filepath.Rel(`/a`, `b`)
	assert.Error(t, err)
	assert.Equal(t, `Rel: can't make b relative to /a`, err.Error())

	// 错误：第一个参数为相对路径，第二个参数为绝对路径
	_, err = filepath.Rel(`a`, `/d`)
	assert.Error(t, err)
	assert.Equal(t, `Rel: can't make /d relative to a`, err.Error())

	// `b` 相对于 `a` 的路径为 `../b`
	s, err = filepath.Rel(`a`, `b`)
	assert.NoError(t, err)
	assert.Equal(t, `../b`, s)

	// `c` 相对于 `a/b` 的路径为 `../../c`
	s, err = filepath.Rel(`a/b`, `c`)
	assert.NoError(t, err)
	assert.Equal(t, `../../c`, s)

	// `a/b/c` 相对于 `a/b` 的路径为 `c`
	s, err = filepath.Rel(`a/b`, `a/b/c`)
	assert.NoError(t, err)
	assert.Equal(t, `c`, s)

	// `a/../c` 相对于 `a/b` 的路径为 `../../c`
	// `a/../c` 即 `c`；`c` 相对于 `a/b` 的路径为 `../../c`
	s, err = filepath.Rel(`a/b`, `a/../c`)
	assert.NoError(t, err)
	assert.Equal(t, `../../c`, s)
}

func TestJoinPath(t *testing.T) {
	s := filepath.Join(`/a`, `b`, ``, `c.txt`)
	assert.Equal(t, `/a/b/c.txt`, s)

	s = filepath.Join(`/a`, `b`, `c`, `..`, `d.txt`)
	assert.Equal(t, `/a/b/d.txt`, s)

	s = filepath.Join(`a`, `b`, `c`, `d.txt`)
	assert.Equal(t, `a/b/c/d.txt`, s)

	s = filepath.Join(`a`, `b`, `c.txt`)
	assert.Equal(t, `a/b/c.txt`, s)

	s = filepath.Join(`a`, `/b`, `//c`, `d.txt`)
	assert.Equal(t, `a/b/c/d.txt`, s)
}

// 移除路径中的无效路径字符
func TestCleanPath(t *testing.T) {
	// 将路径中的 `..`, `/` 字符移除，避免错误的路径字符影响路径
	s := filepath.Clean(`/a/./b/.//c/..///d.txt///`)
	assert.Equal(t, `/a/b/d.txt`, s)
}

// 获取路径的绝对路径
func TestGetAbsolutePath(t *testing.T) {
	// `/` 的绝对路径为 `/`
	s, err := filepath.Abs(`/`)
	assert.NoError(t, err)
	assert.Equal(t, `/`, s)
	assert.True(t, filepath.IsAbs(s)) // 判断路径是否为绝对路径

	// `.` 的绝对路径为 当前路径
	s, err = filepath.Abs(`.`)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(s, `study-golang/basic/io`))
	assert.True(t, filepath.IsAbs(s))

	// `..` 的绝对路径为 当前路径 的上一级路径
	s, err = filepath.Abs(`..`)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(s, `study-golang/basic`))
	assert.True(t, filepath.IsAbs(s))

	// `./a`（即 `a`）的绝对路径为 当前路径 下的 `a` 目录
	s, err = filepath.Abs(`./a`)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(s, `study-golang/basic/io/a`))
	assert.True(t, filepath.IsAbs(s))

	// `a/b/../c`（即 `./a/c`）的绝对路径为 当前路径下 的 `a/c` 目录
	s, err = filepath.Abs(`a/b/../c`)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(s, `study-golang/basic/io/a/c`))
	assert.True(t, filepath.IsAbs(s))
}

// 通过 `:` 字符将字符串分隔为 切片
// 如果一组字符从是以 `:` 连接在一起（例如 Linux 的 $PATH 环境变量），则可以将它们分隔成数组
func TestSplitToPathList(t *testing.T) {
	s := `a:b/c:d/e::/f/g`
	l := filepath.SplitList(s)
	assert.Equal(t, []string{"a", "b/c", "d/e", "", "/f/g"}, l)
}
