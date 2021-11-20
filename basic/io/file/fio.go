package file

import "os"

// 获取文件长度
func FileLength(file *os.File) int {
	if s, err := file.Stat(); err == nil { // 获取文件属性
		return int(s.Size()) // 从文件属性中获取文件实际长度
	}
	return 0
}
