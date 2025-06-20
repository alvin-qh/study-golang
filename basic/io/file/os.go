package file

import (
	"basic/io/file/cross"
	"os"
	"time"
)

// 比较两个文件信息是否相同
func FileInfoCompare(fi1 os.FileInfo, fi2 os.FileInfo) bool {
	if fi1 == fi2 {
		return true
	}

	return fi1.Name() == fi2.Name() &&
		fi1.Size() == fi2.Size() &&
		fi1.Mode() == fi2.Mode() &&
		fi1.ModTime() == fi2.ModTime() &&
		fi1.IsDir() == fi2.IsDir()
}

// 获取所给名称文件的长度
func FileLength(name string) int64 {
	fi, err := os.Stat(name)
	if err != nil {
		return 0
	}
	return fi.Size()
}

// 判断所给名称表示一个路径
func IsDir(name string) bool {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// 判断所给名称表示一个文件
func IsFile(name string) bool {
	return !IsDir(name)
}

// 获取文件模式
func FileMode(name string) os.FileMode {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return 0
	}
	return fileInfo.Mode()
}

// 获取文件的所有者
//
// 第一个返回值为 `uid`, 第二个返回值为 `gid`
func FileOwner(name string) (uint32, uint32) {
	fi, err := os.Stat(name)
	if err != nil {
		return 0, 0
	}

	uid, gid, ok := cross.FileOwner(fi)
	if !ok {
		return 0, 0
	}
	return uid, gid
}

// 获取指定文件的最后修改时间
func FileModTime(name string, utc bool) (time.Time, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return time.UnixMilli(0), err
	}
	if utc {
		return fi.ModTime().UTC(), nil
	}
	return fi.ModTime(), nil
}

// 获取文件的最后访问时间
func FileAccessTime(name string, utc bool) (time.Time, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return time.UnixMilli(0), err
	}

	// 调用跨平台函数获取文件最后一次访问时间
	atime, ok := cross.FileAtime(fi)
	if !ok {
		return time.UnixMilli(0), nil
	}

	if utc {
		return atime.UTC(), nil
	}
	return atime, nil
}
