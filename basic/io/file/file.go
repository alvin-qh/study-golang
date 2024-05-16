package file

import (
	"io"
	"os"
)

// 获取文件指针的当前位置
func GetFilePosition(f *os.File) int64 {
	cur, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0
	}
	return cur
}

// 获取文件长度
func GetFileLength(f *os.File) int64 {
	fi, err := f.Stat()
	if err != nil {
		return 0
	}
	return fi.Size()
}
