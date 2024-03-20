package file

import (
	"io"
	"os"
)

// 获取文件长度
func FileLength(file *os.File) int {
	// 获取文件属性
	if s, err := file.Stat(); err == nil {
		// 从文件属性中获取文件实际长度
		return int(s.Size())
	}
	return 0
}

// 获取文件指针位置
func GetFileCursor(file *os.File) int64 {
	cur, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0
	}
	return cur
}

// 比较两个文件信息是否相同
func CompareFileInfo(fi1 os.FileInfo, fi2 os.FileInfo) bool {
	return fi1.Name() == fi2.Name() &&
		fi1.Size() == fi2.Size() &&
		fi1.Mode() == fi2.Mode() &&
		fi1.ModTime() == fi2.ModTime() &&
		fi1.IsDir() == fi2.IsDir()
}
