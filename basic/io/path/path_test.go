package path

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试获取路径中的各个部分
func TestGetPathPart(t *testing.T) {
	srcPath := "a/b/c/d.xml"

	// 获取路径中的 "目录" 部分
	// 即获取路径 "最后一个分隔符之前" 的部分, 例如 "a/b/c" 的结果为 "a/b"
	targetPath := filepath.Dir(srcPath)
	assert.Equal(t, filepath.Clean("a/b/c"), filepath.Clean(targetPath))

	// 获取路径中的 "文件" 部分
	// 即获取路径 "最后一个分隔符之后" 的部分, 例如 "a/b/c" 的结果为 "c"
	targetPath = filepath.Base(srcPath)
	assert.Equal(t, "d.xml", targetPath)

	// 将路径分割为 "路径" 部分和 "文件" 部分, 相当于同时调用 Dir 函数和 Base 函数
	dir, file := filepath.Split(srcPath)
	assert.Equal(t, filepath.Clean("a/b/c"), filepath.Clean(dir))
	assert.Equal(t, "d.xml", file)

	// 获取路径中文件的扩展名 (如果存在)
	ext := filepath.Ext(srcPath)
	assert.Equal(t, ".xml", ext)
}

// 计算相对路径
//
// 即假设当前所在位置为 第一个路径, 则计算结果为 到达第二个路径的 相对路径
// 计算时, 两个路径必须同为 绝对路径 (即以 `/` 开始) 或 相对路径
func TestRelativePaths(t *testing.T) {
	// `/a` 相对于 `/` 的路径为 `a`
	s, err := filepath.Rel("/", "/a")
	assert.NoError(t, err)
	assert.Equal(t, "a", s)

	// `/a` 相对于 `/a` 的路径为 `.`
	s, err = filepath.Rel("/a", "/a")
	assert.NoError(t, err)
	assert.Equal(t, ".", s)

	// `/b` 相对于 `/a` 的路径为 `../b`
	s, err = filepath.Rel("/a", "/b")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Clean("../b"), filepath.Clean(s))

	// `/a/b` 相对于 `/a` 的路径为 `b`
	s, err = filepath.Rel("/a", "/a/b")
	assert.NoError(t, err)
	assert.Equal(t, "b", s)

	// `/a/b/c` 相对于 `/a` 的路径为 `b/c`
	s, err = filepath.Rel("/a", "/a/b/c")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Clean("b/c"), filepath.Clean(s))

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
	assert.NoError(t, err)
	assert.Equal(t, filepath.Clean("../b"), filepath.Clean(s))

	// `c` 相对于 `a/b` 的路径为 `../../c`
	s, err = filepath.Rel("a/b", "c")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Clean("../../c"), filepath.Clean(s))

	// `a/b/c` 相对于 `a/b` 的路径为 `c`
	s, err = filepath.Rel("a/b", "a/b/c")
	assert.NoError(t, err)
	assert.Equal(t, "c", s)

	// `a/../c` 相对于 `a/b` 的路径为 `../../c`
	// `a/../c` 即 `c`; `c` 相对于 `a/b` 的路径为 `../../c`
	s, err = filepath.Rel("a/b", "a/../c")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Clean("../../c"), filepath.Clean(s))
}

// 将多个部分的路径连接为一个完整的路径
func TestJoinPath(t *testing.T) {
	// 空字符串将被忽略
	s := filepath.Join("/a", "b", "", "c.txt")
	assert.Equal(t, filepath.Clean("/a/b/c.txt"), filepath.Clean(s))

	// `..` 会退回上一级目录
	s = filepath.Join("/a", "b", "c", "..", "d.txt")
	assert.Equal(t, filepath.Clean("/a/b/d.txt"), filepath.Clean(s))

	s = filepath.Join("a", "b", "c.txt")
	assert.Equal(t, filepath.Clean("a/b/c.txt"), s)

	// `/` 或者 `//` 等路径分隔符会被忽略
	s = filepath.Join("a", "/b", "//c", "d.txt")
	assert.Equal(t, filepath.Clean("a/b/c/d.txt"), filepath.Clean(s))
}

// 删除多余的路径分隔符, 遵循的规则为:
//  1. 将连续的多个路径分隔符替换为单个路径分隔符
//  2. 剔除每一个 `.` 路径名元素 (代表当前目录)
//  3. 剔除每一个路径内的 `..` 路径名元素 (代表父目录) 和它前面的非 `..` 路径名元素
//  4. 剔除开始一个根路径的 `..` 路径名元素, 即将路径开始处的 `/..` 替换为 `/` (假设路径分隔符是 `/`)
func TestCleanPath(t *testing.T) {
	// 将路径中的 `..`, `/` 字符移除, 避免错误的路径字符影响路径
	s := filepath.Clean("/a/./b/.//c/..///d.txt///")
	assert.Equal(t, filepath.Clean("/a/b/d.txt"), s)
}

// 获取路径的绝对路径
func TestGetAbsolutePath(t *testing.T) {
	// `/` 的绝对路径为 `/`
	s, err := filepath.Abs("/")
	assert.NoError(t, err)
	if runtime.GOOS == "windows" {
		assert.Regexp(t, `[cCdDeE]:\\`, s)
	} else {
		assert.Regexp(t, "/", s)
	}
	assert.True(t, filepath.IsAbs(s)) // 判断路径是否为绝对路径

	// `.` 的绝对路径为 当前路径
	s, err = filepath.Abs(`.`)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(s, filepath.Clean("study-golang/basic/io/path")))
	assert.True(t, filepath.IsAbs(s))

	// `..` 的绝对路径为 当前路径 的上一级路径
	s, err = filepath.Abs(`..`)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(s, filepath.Clean("study-golang/basic/io")))
	assert.True(t, filepath.IsAbs(s))

	// `./a` (即 `a`) 的绝对路径为 当前路径 下的 `a` 目录
	s, err = filepath.Abs(`./a`)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(s, filepath.Clean("study-golang/basic/io/path/a")))
	assert.True(t, filepath.IsAbs(s))

	// `a/b/../c` (即 `./a/c`) 的绝对路径为 当前路径下 的 `a/c` 目录
	s, err = filepath.Abs(`a/b/../c`)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(s, filepath.Clean("study-golang/basic/io/path/a/c")))
	assert.True(t, filepath.IsAbs(s))
}

// 通过操作系统定义的分隔符将字符串切分为多个路径
//
// 对于 Windows 系统, 分隔符为 `;`, 而 Linux 系统为 `:`
//
// 如果一组字符从是以 `:` (或 `;`) 连接在一起 (例如 $PATH 环境变量), 则可以将它们分割成数组
func TestSplitToPathList(t *testing.T) {
	var s string
	if runtime.GOOS == "windows" {
		s = "a;b/c;d/e;;/f/g"
	} else {
		s = "a:b/c:d/e::/f/g"
	}

	l := filepath.SplitList(s)
	assert.Equal(t, []string{"a", "b/c", "d/e", "", "/f/g"}, l)
}

// 判断所给的路径是否和指定的 Pattern 匹配
//
// pattern:
//
//	{ term }
//
// term:
//
//	'*'                                  匹配0或多个非路径分隔符的字符
//	'?'                                  匹配1个非路径分隔符的字符
//	'[' [ '^' ] { character-range } ']'  字符组 (必须非空)
//	c                                    匹配字符c (c != '*', '?', '\\', '[')
//	'\\' c                               匹配字符c
//
// character-range:
//
//	c           匹配字符c (c != '\\', '-', ']')
//	'\\' c      匹配字符c
//	lo '-' hi   匹配区间[lo, hi]内的字符
func TestPathMatched(t *testing.T) {
	// 判断路径是否和 模式 匹配
	m, err := filepath.Match("*/*.go", "abc/d.go")
	assert.NoError(t, err)
	assert.True(t, m)
}

// 通过给定的模式查找路径
//
// pattern:
//
//	{ term }
//
// term:
//
//	'*'                                  匹配0或多个非路径分隔符的字符
//	'?'                                  匹配1个非路径分隔符的字符
//	'[' [ '^' ] { character-range } ']'  字符组 (必须非空)
//	c                                    匹配字符c (c != '*', '?', '\\', '[')
//	'\\' c                               匹配字符c
//
// character-range:
//
//	c           匹配字符c (c != '\\', '-', ']')
//	'\\' c      匹配字符c
//	lo '-' hi   匹配区间[lo, hi]内的字符
func TestFindPathFile(t *testing.T) {
	files, err := filepath.Glob("./*.go")
	assert.NoError(t, err)

	expected := []string{"path_test.go"}
	sort.Strings(expected)
	sort.Strings(files)

	assert.Equal(t, expected, files)
}

// 遍历指定路径下的所有文件和子目录
// 遍历是通过一个 回调函数 进行的, 没到一个目录下, 就会将该目录下的所有文件和子目录作为参数逐一交给回调函数处理
// 一个目录下的所有内容被处理完毕后, 会进入到其中一个子目录下, 重复进行处理
// 直到所有的目录都被处理完毕
func TestWalk(t *testing.T) {
	files := make([]string, 0, 100)
	dirs := make([]string, 0, 100)

	// 回调文件, 用于处理遍历到的路径或文件
	walkFun := func(path string, info os.FileInfo, err error) error {
		if err != nil { // 判断遍历过程中是否有错误
			return err
		}
		if info.IsDir() { // 判断路径还是文件
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
		return nil
	}

	// 从当前路径开始遍历
	err := filepath.Walk(".", walkFun)
	assert.NoError(t, err)

	filesExpected := []string{"path_test.go"}
	dirsExpected := []string{"."}

	assert.Equal(t, filesExpected, files)
	assert.Equal(t, dirsExpected, dirs)
}

// 对于代码的跨平台兼容性方面, go 语言针对不同平台定义了不同的路径分隔符
//
//	os.PathSeparator, os.PathListSeparator
//
// go 语言对路径统一使用 `/` 和 `:` 进行处理, 前者是路径分隔符, 后者是路径列表分隔符, 所以要想正确适应多平台, 需要路径输入输出的时候做恰当的转换
func TestSlashOperate(t *testing.T) {
	sys := runtime.GOOS

	if sys == "linux" || sys == "darwin" {

		// 类 Unix 系统下, 路径分隔符
		assert.Equal(t, '/', os.PathSeparator)
		assert.Equal(t, ':', os.PathListSeparator)

		// 类 Unix 系统下, 不做转换
		rp := filepath.FromSlash("a/b/c.txt")
		assert.Equal(t, "a/b/c.txt", rp)

		// 类 Unix 系统下, 不做转换
		rp = filepath.ToSlash("a\\b\\c.txt")
		assert.Equal(t, "a\\b\\c.txt", rp)

	} else if sys == "windows" {

		// Windows 系统下, 路径分隔符
		assert.Equal(t, '\\', os.PathSeparator)
		assert.Equal(t, ';', os.PathListSeparator)

		// Windows 系统下, 路径中的 `/` 会被转为 `\\`
		rp := filepath.FromSlash("a/b/c.txt")
		assert.Equal(t, "a\\b\\c.txt", rp)

		// Windows 系统下, 路径中的 `\\` 会被转为 `/`
		rp = filepath.ToSlash("a\\b\\c.txt")
		assert.Equal(t, "a/b/c.txt", rp)

	} else {
		assert.Fail(t, "Not supported")
	}
}

// 判断路径是否存在以及其链接的原始路径
func TestFileSymlinks(t *testing.T) {
	// 判断指定的路径是否存在
	p, err := filepath.EvalSymlinks("./path_test.go")
	assert.NoError(t, err)             // 未返回错误即路径存在
	assert.Equal(t, "path_test.go", p) // 返回和所给路径一致

	// 判断路径是否存在
	_, err = filepath.EvalSymlinks("./path_test.go.1")
	assert.Error(t, err) // 返回错误, 表示路径不存在
	if runtime.GOOS == "windows" {
		assert.Equal(t, "CreateFile path_test.go.1: The system cannot find the file specified.", err.Error())
	} else {
		assert.Equal(t, "lstat path_test.go.1: no such file or directory", err.Error())
	}

	// 创建软链接
	err = os.Symlink("./path_test.go", "./path_test.go.1")
	assert.NoError(t, err)

	defer os.Remove("./path_test.go.1")

	// 判断软链接是否存在
	p, err = filepath.EvalSymlinks("./path_test.go.1")
	assert.NoError(t, err)             // 未返回错误, 即路径存在
	assert.Equal(t, "path_test.go", p) // 返回软链接的源路径
}

// 创建和删除路径
//
//	os.Mkdir, os.MkdirAll 分别可以创建子目录或多级子目录
//	os.Remove, os.RemoveAll 分别可以删除一个空目录或删除路径并同时删除其中的所有内容
func TestMakeAndRemoveDir(t *testing.T) {
	// 创建路径
	err := os.Mkdir("./d", 0755)
	assert.NoError(t, err)

	p, err := filepath.EvalSymlinks("./d") // 判断路径是否存在
	assert.NoError(t, err)
	assert.Equal(t, "d", p)

	// 删除当前路径, 要求路径必须为空, 即内部不能有文件或子目录
	os.Remove("./d")

	// 创建多级路径
	err = os.MkdirAll("./d/e/f", 0755)
	assert.NoError(t, err)

	// 删除路径以及路径下的所有内容
	defer os.RemoveAll("./d")

	p, err = filepath.EvalSymlinks("./d/e/f") // 判断路径是否存在
	assert.NoError(t, err)
	assert.Equal(t, filepath.Clean("d/e/f"), filepath.Clean(p))
}

// 修改当前工作路径
func TestChangeCurrentPath(t *testing.T) {
	// 获取当前工作路径
	dir, err := os.Getwd()
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(dir, filepath.Clean("study-golang/basic/io/path")))

	// 修改当前工作路径 (向上一级)
	err = os.Chdir(filepath.Join(dir, "../.."))
	assert.NoError(t, err)

	// 获取修改工作路径后, 当前工作路径
	dir, err = os.Getwd()
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(dir, filepath.Clean("study-golang/basic")))
}
