package file

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"study/basic/io/pathex"
	"study/basic/os/platform"
	"study/basic/testing/assertion"
	"study/basic/testing/testit"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试获取文件 (或路径) 的属性
func TestOS_Stat(t *testing.T) {
	// 获取路径的信息
	t.Run("for directory", func(t *testing.T) {
		name := pathex.RandomDirName()
		defer os.RemoveAll(name)

		// 创建一个目录
		err := os.Mkdir(name, 0755)
		assert.Nil(t, err)

		// 查看目录信息
		fi, err := os.Stat(name)
		assert.Nil(t, err)

		// 查看目录信息
		assert.True(t, fi.IsDir())
		assert.Equal(t, name, fi.Name())

		platform.RunIfOSNot(platform.Windows, func() {
			// 获取路径的访问权限
			assert.Equal(t, os.FileMode(020000000755), fi.Mode())
		})
	})

	// 获取文件的信息
	t.Run("for file", func(t *testing.T) {
		name := pathex.RandomFileName()
		defer os.Remove(name)

		// 创建文件
		f, err := os.Create(name)
		assert.Nil(t, err)

		f.Close()

		// 获取文件信息
		fi, err := os.Stat(name)
		assert.Nil(t, err)

		// 查看文件信息
		assert.False(t, fi.IsDir())
		assert.Equal(t, name, fi.Name())

		platform.RunIfOSNot(platform.Windows, func() {
			assert.Equal(t, os.FileMode(0644), fi.Mode())
		})
	})

	// 获取文件所有者
	t.Run("get file owner", func(t *testing.T) {
		testit.SkipTimeOnOS(t, platform.Windows)

		name := pathex.RandomFileName()
		defer os.Remove(name)

		u, err := user.Current()
		assert.Nil(t, err)

		// 创建文件
		f, err := os.Create(name)
		assert.Nil(t, err)

		f.Close()

		// 获取文件所有者
		uid, gid := FileOwner(name)

		// 确认文件所有者为当前用户和用户组
		assert.Equal(t, u.Uid, strconv.FormatUint(uint64(uid), 10))
		assert.Equal(t, u.Gid, strconv.FormatUint(uint64(gid), 10))
	})
}

// 测试写入和读取文件
//
// `os.WriteFile` 相当于快捷写文件, 及将指定的内容一次性写入指定文件名的文件中, 如果文件不存在则创建文件
//
// `os.ReadFile` 相当于快捷读文件, 及一次性读取指定名称文件的所有内容, 如果文件不存在则返回错误
func TestOS_WriteReadFile(t *testing.T) {
	// 文件名
	name := pathex.RandomFileName()
	defer os.Remove(name)

	// 文件内容
	c := []byte("hello world")

	// 写入
	err := os.WriteFile(name, c, 0644)
	assert.Nil(t, err)

	// 读取
	data, err := os.ReadFile(name)
	assert.Nil(t, err)
	assert.Equal(t, c, data)
}

// 测试创建文件
//
// 创建新文件, 文件访问权限为 `666`, 文件已存在时返回错误
//
// 该方法相当于 `os.OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)` 函数调用的简化方式
func TestOS_Create(t *testing.T) {
	name := pathex.RandomFileName()
	defer os.Remove(name)

	// 创建文件
	f, err := os.Create(name)
	assert.Nil(t, err)

	// 关闭文件
	f.Close()

	assert.FileExists(t, name)

	assert.True(t, IsFile(name))
	platform.RunIfOSNot(platform.Windows, func() {
		assert.Equal(t, os.FileMode(0644), FileMode(name))
	})
	assert.Equal(t, int64(0), FileLength(name))
}

// 测试打开文件
//
// 以只读方式打开文件, 文件不存在时返回错误
//
// 该方法相当于 `os.OpenFile(name, O_RDONLY, 0)` 函数调用的简化方式
func TestOS_Open(t *testing.T) {
	name := pathex.RandomFileName()
	defer os.Remove(name)

	// 打开不存在的文件, 返回错误
	_, err := os.Open(name)
	if platform.IsOSMatch(platform.Windows) {
		assert.EqualError(t, err, fmt.Sprintf("open %s: The system cannot find the file specified.", name))
	} else {
		assert.EqualError(t, err, fmt.Sprintf("open %s: no such file or directory", name))
	}

	// 创建文件
	f, err := os.Create(name)
	assert.Nil(t, err)

	f.Close()

	assert.FileExists(t, name)

	// 打开已存在的文件
	f, err = os.Open(name)
	assert.Nil(t, err)

	f.Close()

	assert.True(t, IsFile(name))
	platform.RunIfOSNot(platform.Windows, func() {
		assert.Equal(t, os.FileMode(0644), FileMode(name))
	})
	assert.Equal(t, int64(0), FileLength(name))
}

// 测试打开或创建文件
//
// `os.OpenFile` 函数可用于创建, 打开, 截断等文件操作
//
// `flag` 参数用于定义文件打开方式, 可以通过 `|` 运算符组合多个值, 包括:
//   - `O_RDONLY` 以只读方式打开文件
//   - `O_WRONLY` 以只写方式打开文件
//   - `O_RDWR`   以读写方式打开文件
//   - `O_APPEND` 打开文件并在文件末尾追加
//   - `O_CREATE` 创建新文件
//   - `O_EXCL`   配合 `O_CREATE` 使用, 表示创建文件时文件不能存在
//   - `O_SYNC`   创建异步文件
//   - `O_TRUNC`  打开已存在的文件后清空文件内容
//
// 注意, 不同操作系统支持的 `flag` 参数也略有不同
//
// `perm` 参数表示文件的权限, 可以是 `os.ModeXXX` 值的组合 (通过 `|` 运算符), 也可以是一个 8 进制整数
func TestOS_OpenFile(t *testing.T) {
	name := pathex.RandomFileName()
	defer os.Remove(name)

	// 创建文件
	f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC, 0644)
	assert.Nil(t, err)

	f.Close()

	assert.FileExists(t, name)

	assert.True(t, IsFile(name))
	platform.RunIfOSNot(platform.Windows, func() {
		assert.Equal(t, os.FileMode(0644), FileMode(name))
	})
	assert.Equal(t, int64(0), FileLength(name))

	// 打开文件
	f, err = os.OpenFile(name, os.O_RDONLY, 0)
	assert.Nil(t, err)

	f.Close()

	assert.True(t, IsFile(name))
	platform.RunIfOSNot(platform.Windows, func() {
		assert.Equal(t, os.FileMode(0644), FileMode(name))
	})
	assert.Equal(t, int64(0), FileLength(name))
}

// 测试创建临时文件
//
// 创建临时文件需要指定文件的路径 (`dir` 参数) 以及文件名的模式 (`pattern` 参数)
//
// 文件名的模式包括: 前缀, 占位符和后缀, 后两者可选
//
// `dir` 参数如果为空字符串, 则使用操作系统指定的临时路径
func TestOS_CreateTemp(t *testing.T) {
	f, err := os.CreateTemp("", "temp-*.txt")
	assert.Nil(t, err)

	name := f.Name()
	defer os.Remove(name)

	f.Close()

	assert.True(t, IsFile(name))
	assert.Regexp(t, `temp-\d+\.txt`, name)
	assert.Equal(t, int64(0), FileLength(name))

	fmt.Printf("%s/%s\n", filepath.Dir(name), name)
}

// 测试创建指定目录
func TestOS_Mkdir(t *testing.T) {
	name := pathex.RandomDirName()

	// 创建目录
	err := os.Mkdir(name, 0755)
	assert.Nil(t, err)

	// 结束后删除目录
	defer os.RemoveAll(name)

	assert.True(t, IsDir(name))
	platform.RunIfOSNot(platform.Windows, func() {
		// 文件模式包括目录模式即权限
		assert.Equal(t, os.FileMode(020000000755), FileMode(name))
	})
}

// 测试创建路径
//
// 创建路径会自动创建路径节点中所有的目录
func TestOS_MkdirAll(t *testing.T) {
	names := []string{
		pathex.RandomDirName(),
		pathex.RandomDirName(),
		pathex.RandomDirName(),
		pathex.RandomDirName(),
	}

	// 创建目录, 通过一个多级路径字符串创建全部目录
	err := os.MkdirAll(filepath.Join(names...), 0755)
	assert.Nil(t, err)

	// 结束后删除目录
	defer os.RemoveAll(names[0])

	// 确认各级目录创建成功
	assert.DirExists(t, names[0])
	assert.DirExists(t, filepath.Join(names[:len(names)-3]...))
	assert.DirExists(t, filepath.Join(names[:len(names)-2]...))
	assert.DirExists(t, filepath.Join(names...))
}

// 测试删除文件或目录
func TestOS_Remove(t *testing.T) {
	// 删除文件
	t.Run("delete file", func(t *testing.T) {
		name := pathex.RandomFileName()

		// 创建一个文件
		f, err := os.Create(name)
		assert.Nil(t, err)

		f.Close()

		// 确认文件存在
		assert.FileExists(t, name)

		// 删除文件
		err = os.Remove(name)
		assert.Nil(t, err)

		// 确认文件不存在
		assert.NoFileExists(t, name)
	})

	// 删除目录
	t.Run("delete directory", func(t *testing.T) {
		name := pathex.RandomDirName()

		// 创建一个目录
		err := os.Mkdir(name, os.FileMode(0755))
		assert.Nil(t, err)

		// 确认目录存在
		assert.DirExists(t, name)

		// 删除目录
		err = os.Remove(name)
		assert.Nil(t, err)

		// 确认目录不存在
		assert.NoDirExists(t, name)
	})
}

// 测试删除目录及子目录和目录中的文件
func TestOS_RemoveAll(t *testing.T) {
	names := []string{
		pathex.RandomDirName(),
		pathex.RandomDirName(),
		pathex.RandomDirName(),
		pathex.RandomDirName(),
	}

	// 创建目录
	err := os.MkdirAll(filepath.Join(names...), 0755)
	assert.Nil(t, err)

	// 确认目录创建成功
	assert.DirExists(t, names[0])

	// 删除目录, 因为目录不为空, 则删除失败
	err = os.Remove(names[0])
	if platform.IsOSMatch(platform.Windows) {
		assert.EqualError(t, err, fmt.Sprintf("remove %s: The directory is not empty.", names[0]))
	} else {
		assert.EqualError(t, err, fmt.Sprintf("remove %s: directory not empty", names[0]))
	}

	// 删除目录, 删除第一级目录, 则其中的内容会一并删除
	err = os.RemoveAll(names[0])
	assert.Nil(t, err)

	// 确认目录不存在
	assert.NoDirExists(t, names[0])
}

// 测试移动文件
//
// 移动文件指的是将文件从一个路径移动到另一个路径, 包括所咋目录的变更和文件名的更改
func TestOS_Move(t *testing.T) {
	// 测试文件移动
	t.Run("move file", func(t *testing.T) {
		// 定义两个文件名
		nameBefore := pathex.RandomFileName()
		nameAfter := pathex.RandomFileName()

		// 创建文件
		f, err := os.Create(nameBefore)
		assert.Nil(t, err)

		defer os.Remove(nameBefore)
		defer os.Remove(nameAfter)

		f.Close()

		// 确认文件存在
		assert.FileExists(t, nameBefore)

		// 移动文件
		err = os.Rename(nameBefore, nameAfter)
		assert.Nil(t, err)

		// 确认文件不存在
		assert.NoFileExists(t, nameBefore)

		// 确认文件存在
		assert.FileExists(t, nameAfter)
	})

	// 测试目录移动
	t.Run("move dir", func(t *testing.T) {
		// 定义两个文件名
		nameBefore := pathex.RandomDirName()
		nameAfter := pathex.RandomDirName()

		// 创建文件
		err := os.Mkdir(nameBefore, os.FileMode(0755))
		assert.Nil(t, err)

		defer os.RemoveAll(nameBefore)
		defer os.RemoveAll(nameAfter)

		// 确认文件存在
		assert.DirExists(t, nameBefore)

		// 移动文件
		err = os.Rename(nameBefore, nameAfter)
		assert.Nil(t, err)

		// 确认文件不存在
		assert.NoDirExists(t, nameBefore)

		// 确认文件存在
		assert.DirExists(t, nameAfter)
	})
}

// 测试创建临时文件夹
//
// 创建临时文件需要指定文件的路径 (`dir` 参数) 以及文件名的模式 (`pattern` 参数)
//
// 文件夹名称的模式包括: 前缀, 占位符和后缀, 后两者可选
//
// `dir` 参数如果为空字符串, 则使用操作系统指定的临时路径
func TestOS_CreateTempDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "temp-*")
	assert.Nil(t, err)

	defer os.RemoveAll(dir)

	assert.True(t, IsDir(dir))
	assert.Regexp(t, `temp-\d+`, dir)
}

// 测试获取当前工作路径
func TestOS_Getwd(t *testing.T) {
	// 获取当前工作路径
	p, err := os.Getwd()
	assert.Nil(t, err)

	assertion.PathEndsWith(t, p, "/basic/io/file")
}

// 测试改变当前工作路径
func TestOS_Chdir(t *testing.T) {
	name := pathex.RandomDirName()

	// 在当前路径下创建子目录
	err := os.Mkdir(name, 0755)
	assert.Nil(t, err)

	defer os.RemoveAll(name)

	// 获取当前工作路径
	err = os.Chdir(name)
	assert.Nil(t, err)

	// 改变当前工作目录
	p, err := os.Getwd()
	assert.Nil(t, err)

	// 测试结束后恢复初始工作目录
	defer os.Chdir("..")
	assertion.PathEndsWith(t, p, "/basic/io/file/"+name)
}

// 测试修改文件权限
func TestOS_Chmod(t *testing.T) {
	testit.SkipTimeOnOS(t, platform.Windows)

	name := pathex.RandomFileName()
	defer os.Remove(name)

	// 创建文件
	f, err := os.Create(name)
	assert.Nil(t, err)

	f.Close()

	// 获取文件信息, 确认文件权限
	assert.Equal(t, os.FileMode(0644), FileMode(name))

	// 修改文件权限
	os.Chmod(name, 0755)

	// 重新获取文件信息, 确认修改后的文件权限
	assert.Equal(t, os.FileMode(0755), FileMode(name))
}

// 测试修改文件的所有者和组
func TestOS_Chown(t *testing.T) {
	testit.SkipTimeOnOS(t, platform.Windows)

	name := pathex.RandomFileName()
	defer os.Remove(name)

	f, err := os.Create(name)
	assert.Nil(t, err)
	f.Close()

	// 获取 root 用户
	u, err := user.Lookup("root")
	assert.Nil(t, err)

	// 转换整数类型 uid
	uid, err := strconv.Atoi(u.Uid)
	assert.Nil(t, err)

	// 转换整数类型 gid
	gid, err := strconv.Atoi(u.Gid)
	assert.Nil(t, err)

	err = os.Chown(name, uid, gid)
	assert.EqualError(t, err, fmt.Sprintf("chown %s: operation not permitted", name))
}

// 测试截断文件
//
// 将文件长度截断为所给的长度
func TestOS_Truncate(t *testing.T) {
	name := pathex.RandomFileName()
	defer os.Remove(name)

	// 创建并写入文件内容
	err := os.WriteFile(name, []byte("Hello World"), 0644)
	assert.Nil(t, err)
	assert.Equal(t, int64(11), FileLength(name))

	// 将文件截取到剩余 3 字节
	os.Truncate(name, 3)
	assert.Equal(t, int64(3), FileLength(name))

	bs, err := os.ReadFile(name)
	assert.Nil(t, err)
	assert.Equal(t, "Hel", string(bs))

	// 将文件截断为 0 字节, 相当于清空文件内容
	err = os.Truncate(name, 0)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), FileLength(name))
}

// 测试修改文件的修改时间和最后访问时间
func TestOS_Chtimes(t *testing.T) {
	name := pathex.RandomFileName()
	defer os.Remove(name)

	// 创建文件
	f, err := os.Create(name)
	assert.Nil(t, err)

	f.Close()

	// 指定一个时间表示文件的最后修改时间
	mtime, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	assert.Nil(t, err)

	// 修改文件的最后修改时间
	// 空的 time.Time 实例表示不修改对应的时间值, 本例表示不修改文件的最后访问时间
	err = os.Chtimes(name, time.Time{}, mtime)
	assert.Nil(t, err)

	// 获取文件的最后一次修改时间
	mtime, err = FileModTime(name, true)
	assert.Nil(t, err)

	// 获取文件的最后访问时间
	atime, err := FileAccessTime(name, true)
	assert.Nil(t, err)

	// 确认文件的最后修改时间被改变
	assert.Equal(t, "2020-01-01T00:00:00Z", mtime.Format(time.RFC3339))
	// 确认文件的最后访问时间未被改变
	assert.NotEqual(t, "2021-01-01T00:00:00Z", atime.Format(time.RFC3339))

	// 指定一个时间表示文件的最后访问时间
	atime, err = time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	assert.Nil(t, err)

	// 同时修改文件的的最后修改时间和最后访问时间
	err = os.Chtimes(name, atime, mtime)
	assert.Nil(t, err)

	// 再次获取文件的最后修改时间
	mtime, err = FileModTime(name, true)
	assert.Nil(t, err)

	// 再次获取文件的最后访问时间
	atime, err = FileAccessTime(name, true)
	assert.Nil(t, err)

	// 确认两个时间都被改变
	assert.Equal(t, "2020-01-01T00:00:00Z", mtime.Format(time.RFC3339))
	assert.Equal(t, "2021-01-01T00:00:00Z", atime.Format(time.RFC3339))
}

// 测试创建文件硬链接
//
// 所谓硬链接, 相当于为原文件创建一个同步文件: 链接文件和原文件各自独立, 且链接文件的内容和原文件一致
//
// 当原文件删除后, 硬链接文件不受影响
func TestOS_Link(t *testing.T) {
	// 定义原文件名称
	fname := pathex.RandomFileName()
	defer os.Remove(fname)

	// 定义硬链接文件名称
	lname := fmt.Sprintf("%s.link", fname)
	defer os.Remove(lname)

	// 创建原文件
	f, err := os.Create(fname)
	assert.Nil(t, err)

	f.Close()

	// 确认原文件存在
	assert.FileExists(t, fname)

	// 创建一个文件链接
	err = os.Link(fname, lname)
	assert.Nil(t, err)

	// 确认硬链接文件存在
	assert.FileExists(t, lname)

	// 向原文件写入内容
	err = os.WriteFile(fname, []byte("Hello World"), 0)
	assert.Nil(t, err)

	// 确认硬链接文件和原文件内容相同
	data, err := os.ReadFile(lname)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Hello World"), data)

	// 删除原文件
	os.Remove(fname)

	// 确认原文件删除后, 硬链接文件仍可读取
	data, err = os.ReadFile(lname)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Hello World"), data)
}

// 测试创建文件软链接
//
// 所谓软链接, 相当于为原文件创建一个别名, 和原文件表示同一个文件
//
// 当原文件删除后, 软连接文件也不存在
func TestOS_Symlink(t *testing.T) {
	// 定义原文件名称
	fname := pathex.RandomFileName()
	defer os.Remove(fname)

	// 定义软链接文件名称
	lname := fmt.Sprintf("%s.link", fname)
	defer os.Remove(lname)

	// 创建原文件
	f, err := os.Create(fname)
	assert.Nil(t, err)

	f.Close()

	// 确认原文件存在
	assert.FileExists(t, fname)

	// 创建一个文件软链接
	err = os.Symlink(fname, lname)
	assert.Nil(t, err)

	// 确认软链接文件存在
	assert.FileExists(t, lname)

	c := []byte("Hello World")

	// 向原文件写入内容
	err = os.WriteFile(fname, c, 0)
	assert.Nil(t, err)

	// 确认软链接文件和原文件内容相同
	data, err := os.ReadFile(lname)
	assert.Nil(t, err)
	assert.Equal(t, c, data)

	// 删除原文件
	os.Remove(fname)
	// 确认软链接文件仍存在
	assert.FileExists(t, lname)

	// 确认原文件删除后, 软链接无法读取
	_, err = os.ReadFile(lname)
	if platform.IsOSMatch(platform.Windows) {
		assert.EqualError(t, err, fmt.Sprintf("open %s: The system cannot find the file specified.", lname))
	} else {
		assert.EqualError(t, err, fmt.Sprintf("open %s: no such file or directory", lname))
	}
}
