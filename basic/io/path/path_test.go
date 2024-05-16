// Go 标准库中包含两套路径处理库, 分别为:
//   - 位于 `path` 包下面的函数
//   - 位于 `filepath` 包下面的函数
//
// 这两个包下面的函数功能基本一致, 而 `filepath` 包支持跨平台操作, 在任意平台上均可以使用 `/` 字符作为路径分隔符
package path

import (
	"os"
	"path/filepath"
	"strings"
	"study/basic/expression"
	"study/basic/os/platform"
	"study/basic/testing/assertion"
	"study/basic/testing/testit"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 对于代码的跨平台兼容性方面, go 语言针对不同平台定义了不同的路径分隔符
//
// 路径分隔符定义为: `os.PathSeparator`, `os.PathListSeparator`
//
// Go 语言对路径统一使用 `/` 和 `:` 进行处理, 前者是路径分隔符, 后者是路径列表分隔符, 所以要想正确适应多平台,
// 需要路径输入输出的时候做恰当的转换
func TestFilePath_Separator(t *testing.T) {
	// Windows 系统下, 路径分隔符
	testit.RunIf(t, platform.Windows, func(t *testing.T) {
		assert.Equal(t, '\\', os.PathSeparator)
		assert.Equal(t, ';', os.PathListSeparator)
	})

	// 类 Unix 系统下, 路径分隔符
	testit.RunIf(t, platform.Linux|platform.Darwin, func(t *testing.T) {
		assert.Equal(t, '/', os.PathSeparator)
		assert.Equal(t, ':', os.PathListSeparator)
	})
}

// 测试将路径中的路径分隔符统一替换为 `/` 字符
//
// 即, 在 Windows 上, 会将 Windows 风格路径分隔符转为 Unix 风格路径分隔符, 其它系统不做处理
func TestFilePath_ToSlash(t *testing.T) {
	p := filepath.ToSlash(expression.Must(os.Getwd()))
	assert.True(t, strings.HasSuffix(p, `/study-golang/basic/io/path`))
}

// 测试将路径中的 `/` 字符统一替换操作系统指定的路径分隔符
//
// 即, 在 Windows 上, 会将 Unix 风格路径分隔符转为 Windows 风格路径分隔符, 其它系统不做处理
func TestFilePath_FromSlash(t *testing.T) {
	p := filepath.FromSlash(expression.Must(os.Getwd()))

	testit.RunIf(t, platform.Windows, func(t *testing.T) {
		assert.True(t, strings.HasSuffix(p, `\study-golang\basic\io\path`))
	})

	testit.RunIf(t, platform.Linux|platform.Darwin, func(t *testing.T) {
		assert.True(t, strings.HasSuffix(p, `/study-golang/basic/io/path`))
	})
}

// 将路径中的无效部分进行清理
//
// `filepath.Clean` 函数返回和其参数等效的最短路径, 简化的规则包括:
//
//  1. 将多个连续重复的路径分隔符替换为一个
//  2. 消除路径中表示当前路径的 `.` 部分
//  3. 删除路径中开头的 `.` 字符, 即, 在路径开头用 `/` 替换 `/.` (假设路径为 Unix 风格)
//  4. 删除路径末尾的 `\` 或 `/` 分隔符, 除非路径为根路径 (即 `C:\` 或 `/`)
//  5. 将路径中的 `/` 字符统一替换为操作系统指定的路径分隔符, 注意: `filepath.Clean` 函数
//     会将 `/` 字符统一替换为操作系统指定的路径分隔符, 而不会将 `/` 字符统一替换为操作系统指定的路径分隔符
//  6. 将空字符串替换为 `.`
func TestFilePath_Clean(t *testing.T) {
	p := platform.Choose(platform.Windows, `C:\\Users\\\alvin\\AppData`, `C://Users///alvin//AppData`)
	p = filepath.Clean(p)
	assert.Equal(t, "C:/Users/alvin/AppData", filepath.ToSlash(p))

	p = platform.Choose(platform.Windows, `.\path_test.go`, `./path_test.go`)
	p = filepath.Clean(p)
	assert.Equal(t, "path_test.go", p)

	p = platform.Choose(platform.Windows, `\test\..\path_test.go`, `/test/../path_test.go`)
	p = filepath.Clean(p)
	assert.Equal(t, "/path_test.go", filepath.ToSlash(p))

	p = `/io/path/path_test.go`
	p = filepath.Clean(p)
	assert.Equal(t, platform.Choose(platform.Windows, `\io\path\path_test.go`, `/io/path/path_test.go`), p)

	p = filepath.Clean("")
	assert.Equal(t, ".", p)
}

// 测试获取字符串中表示目录的部分
func TestFilePath_Dir(t *testing.T) {
	src := "a/b/c/d.xml"

	// 获取路径中的 "目录" 部分
	// 即获取路径 "最后一个分隔符之前" 的部分, 例如 "a/b/c" 的结果为 "a/b"
	target := filepath.Dir(src)
	assert.Equal(t, `a/b/c`, filepath.ToSlash(target))
}

// 测试获取字符串中表示文件的部分
func TestFilePath_Base(t *testing.T) {
	src := "a/b/c/d.xml"

	// 获取路径中的 "文件" 部分
	// 即获取路径 "最后一个分隔符之后" 的部分, 例如 "a/b/c" 的结果为 "c"
	target := filepath.Base(src)
	assert.Equal(t, "d.xml", target)
}

// 测试将字符串分为表示目录和文件的两部分
func TestFilePath_Split(t *testing.T) {
	src := "a/b/c/d.xml"

	// 将路径分割为 "路径" 部分和 "文件" 部分, 相当于同时调用 Dir 函数和 Base 函数
	dir, file := filepath.Split(src)
	assert.Equal(t, `a/b/c/`, filepath.ToSlash(dir))
	assert.Equal(t, "d.xml", file)
}

// 测试获取字符串中表示文件扩展名的部分 (如果有)
func TestFilePath_Ext(t *testing.T) {
	src := "a/b/c/d.xml"

	// 获取路径中文件的扩展名 (如果存在)
	ext := filepath.Ext(src)
	assert.Equal(t, ".xml", ext)
}

// 计算相对路径
//
// 即假设当前所在位置为第一个路径, 则计算结果为到达第二个路径的相对路径
// 计算时, 两个路径必须同为绝对路径 (即以 `/` 开始) 或相对路径
func TestFilePath_Rel(t *testing.T) {
	// `/a` 相对于 `/` 的路径为 `a`
	s, err := filepath.Rel("/", "/a")
	assert.Nil(t, err)
	assert.Equal(t, "a", s)

	// `/a` 相对于 `/a` 的路径为 `.`
	s, err = filepath.Rel("/a", "/a")
	assert.Nil(t, err)
	assert.Equal(t, ".", s)

	// `/b` 相对于 `/a` 的路径为 `../b`
	s, err = filepath.Rel("/a", "/b")
	assert.Nil(t, err)
	assert.Equal(t, `../b`, filepath.ToSlash(s))

	// `/a/b` 相对于 `/a` 的路径为 `b`
	s, err = filepath.Rel("/a", "/a/b")
	assert.Nil(t, err)
	assert.Equal(t, "b", s)

	// `/a/b/c` 相对于 `/a` 的路径为 `b/c`
	s, err = filepath.Rel("/a", "/a/b/c")
	assert.Nil(t, err)
	assert.Equal(t, `b/c`, filepath.ToSlash(s))

	// 错误: 第一个参数为绝对路径, 第二个参数为相对路径
	_, err = filepath.Rel("/a", "b")
	assert.Error(t, err)
	assert.Equal(t, "Rel: can't make b relative to /a", err.Error())

	// 错误: 第一个参数为相对路径, 第二个参数为绝对路径
	_, err = filepath.Rel("a", "/d")
	assert.Error(t, err)
	assert.Equal(t, "Rel: can't make /d relative to a", err.Error())

	// `b` 相对于 `a` 的路径为 `../b`
	s, err = filepath.Rel("a", "b")
	assert.Nil(t, err)
	assert.Equal(t, `../b`, filepath.ToSlash(s))

	// `c` 相对于 `a/b` 的路径为 `../../c`
	s, err = filepath.Rel("a/b", "c")
	assert.Nil(t, err)
	assert.Equal(t, `../../c`, filepath.ToSlash(s))

	// `a/b/c` 相对于 `a/b` 的路径为 `c`
	s, err = filepath.Rel("a/b", "a/b/c")
	assert.Nil(t, err)
	assert.Equal(t, "c", s)

	// `a/../c` 相对于 `a/b` 的路径为 `../../c`
	// `a/../c` 即 `c`; `c` 相对于 `a/b` 的路径为 `../../c`
	s, err = filepath.Rel("a/b", "a/../c")
	assert.Nil(t, err)
	assert.Equal(t, `../../c`, filepath.ToSlash(s))
}

// 将多个部分的路径连接为一个完整的路径
func TestFilePath_Join(t *testing.T) {
	// 空字符串将被忽略
	s := filepath.Join("/a", "b", "", "c.txt")
	assert.Equal(t, `/a/b/c.txt`, filepath.ToSlash(s))

	// `..` 会退回上一级目录
	s = filepath.Join("/a", "b", "c", "..", "d.txt")
	assert.Equal(t, `/a/b/d.txt`, filepath.ToSlash(s))

	s = filepath.Join("a", "b", "c.txt")
	assert.Equal(t, `a/b/c.txt`, filepath.ToSlash(s))

	// `/` 或者 `//` 等路径分隔符会被忽略
	s = filepath.Join("a", "/b", "//c", "d.txt")
	assert.Equal(t, `a/b/c/d.txt`, filepath.ToSlash(s))
}

// 测试判断一个路径是否为绝对路径
func TestFilePath_IsAbs(t *testing.T) {
	testit.RunIf(t, platform.Windows, func(t *testing.T) {
		assert.True(t, filepath.IsAbs(`c:\`))
		assert.True(t, filepath.IsAbs(`c:\a\b\c\d.txt`))
		assert.True(t, filepath.IsAbs(`c:\a\b\c`))
		assert.True(t, filepath.IsAbs(`c:\a\b`))
		assert.True(t, filepath.IsAbs(`c:\a`))

		assert.False(t, filepath.IsAbs(`a\b\c\d.txt`))
		assert.False(t, filepath.IsAbs(`a\b\c`))
		assert.False(t, filepath.IsAbs(`a\b`))
		assert.False(t, filepath.IsAbs(`a`))
	})

	testit.RunIf(t, platform.Linux|platform.Darwin, func(t *testing.T) {
		assert.True(t, filepath.IsAbs("/a"))
		assert.True(t, filepath.IsAbs("/a/b"))
		assert.True(t, filepath.IsAbs("/a/b/c"))
		assert.True(t, filepath.IsAbs("/a/b/c/d.txt"))

		assert.False(t, filepath.IsAbs("a"))
		assert.False(t, filepath.IsAbs("a/b"))
		assert.False(t, filepath.IsAbs("a/b/c"))
		assert.False(t, filepath.IsAbs("a/b/c/d.txt"))

	})
}

// 通过相对路径获取绝对路径
func TestFilePath_Abs(t *testing.T) {
	// `/` 的绝对路径为 `/`
	s, err := filepath.Abs("/")
	assert.Nil(t, err)
	assert.Regexp(t, platform.Choose(platform.Windows, `[cCdDeE]:\\`, `/`), s)

	// `.` 的绝对路径为 当前路径
	s, err = filepath.Abs(`.`)
	assert.Nil(t, err)
	assertion.PathEndsWith(t, s, `/study-golang/basic/io/path`)

	// `..` 的绝对路径为 当前路径 的上一级路径
	s, err = filepath.Abs(`..`)
	assert.Nil(t, err)
	assert.True(t, assertion.PathEndsWith(t, s, `/study-golang/basic/io`))

	// `a/b/../c` (即 `./a/c`) 的绝对路径为 当前路径下 的 `a/c` 目录
	s, err = filepath.Abs(`a/b/../c`)
	assert.Nil(t, err)
	assert.True(t, assertion.PathEndsWith(t, s, `/study-golang/basic/io/path/a/c`))
}

// 通过操作系统定义的分隔符将字符串切分为多个路径
//
// 对于 Windows 系统, 分隔符为 `;`, 而 Linux 系统为 `:`
//
// 如果一组字符从是以 `:` (或 `;`) 连接在一起 (例如 $PATH 环境变量), 则可以将它们分割成数组
func TestFilePath_SplitList(t *testing.T) {
	r := filepath.SplitList(
		platform.Choose(
			platform.Windows,
			`a;b/c;d/e;;/f/g`,
			`a:b/c:d/e::/f/g`,
		),
	)
	assert.Equal(t, []string{"a", "b/c", "d/e", "", "/f/g"}, r)
}

// 判断所给的路径是否和指定的 Pattern 匹配
//
// pattern:
//
//	{ term }
//
// term:
//   - '*'                                  匹配0或多个非路径分隔符的字符
//   - '?'                                  匹配1个非路径分隔符的字符
//   - '[' [ '^' ] { character-range } ']'  字符组 (必须非空)
//   - c                                    匹配字符c (c != '*', '?', '\\', '[')
//   - '\\' c                               匹配字符c
//
// character-range:
//
//   - c           匹配字符c (c != '\\', '-', ']')
//   - '\\' c      匹配字符c
//   - lo '-' hi   匹配区间[lo, hi]内的字符
func TestFilePath_Match(t *testing.T) {
	// 判断路径是否和 模式 匹配
	m, err := filepath.Match("*/*.go", "abc/d.go")

	assert.Nil(t, err)
	assert.True(t, m)
}

// 通过给定的模式查找路径
//
// pattern:
//
//	{ term }
//
// term:
//   - '*'                                  匹配0或多个非路径分隔符的字符
//   - '?'                                  匹配1个非路径分隔符的字符
//   - '[' [ '^' ] { character-range } ']'  字符组 (必须非空)
//   - c                                    匹配字符c (c != '*', '?', '\\', '[')
//   - '\\' c                               匹配字符c
//
// character-range:
//   - c           匹配字符c (c != '\\', '-', ']')
//   - '\\' c      匹配字符c
//   - lo '-' hi   匹配区间[lo, hi]内的字符
func TestFilePath_Glob(t *testing.T) {
	files, err := filepath.Glob("./*.go")

	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"path_test.go"}, files)
}

// 遍历指定路径下的所有文件和子目录
//
// 遍历是通过一个 回调函数 进行的, 没到一个目录下, 就会将该目录下的所有文件和子目录作为参数逐一交给回调函数处理
// 一个目录下的所有内容被处理完毕后, 会进入到其中一个子目录下, 重复进行处理直到所有的目录都被处理完毕
func TestFilePath_Walk(t *testing.T) {
	files := make([]string, 0, 100)
	dirs := make([]string, 0, 100)

	// 从当前路径开始遍历
	// 回调文件, 用于处理遍历到的路径或文件
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// 判断遍历过程中是否有错误
		if err != nil {
			return err
		}
		// 判断路径还是文件
		if info.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, []string{"path_test.go"}, files)
	assert.Equal(t, []string{"."}, dirs)
}

// 判断路径是否存在以及其链接的原始路径
func TestFilePath_EvalSymlinks(t *testing.T) {
	// 判断指定的路径是否存在
	p, err := filepath.EvalSymlinks("./path_test.go")
	// 未返回错误即路径存在
	assert.Nil(t, err)
	// 返回和所给路径一致
	assert.Equal(t, "path_test.go", p)

	// 判断路径是否存在
	_, err = filepath.EvalSymlinks("./path_test.go.1")
	// 返回错误, 表示路径不存在
	assert.Error(t, err)

	// 确认错误信息
	if platform.IsOSMatch(platform.Windows) {
		assert.Equal(t, "CreateFile path_test.go.1: The system cannot find the file specified.", err.Error())
	} else {
		assert.Equal(t, "lstat path_test.go.1: no such file or directory", err.Error())
	}

	// 创建软链接
	err = os.Symlink("./path_test.go", "./path_test.go.1")
	assert.Nil(t, err)

	defer os.Remove("./path_test.go.1")

	// 判断软链接是否存在
	p, err = filepath.EvalSymlinks("./path_test.go.1")
	// 未返回错误, 即路径存在
	assert.Nil(t, err)
	// 返回软链接的源路径
	assert.Equal(t, "path_test.go", p)
}
