package file

import (
	"errors"
	"os"
)

// 获取文件长度
func FileLength(file *os.File) int {
	if s, err := file.Stat(); err == nil { // 获取文件属性
		return int(s.Size()) // 从文件属性中获取文件实际长度
	}
	return 0
}

// 获取文件指针位置
func GetFileCursor(file *os.File) int64 {
	cur, err := file.Seek(0, os.SEEK_CUR)
	if err != nil {
		return 0
	}
	return cur
}

// 关闭并删除文件
func CloseAndRemoveFile(f *os.File) error {
	if err := f.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
		return err
	}
	return os.Remove(f.Name())
}

type User struct {
	Id    int64
	Name  string
	Email string
	Phone []string
}

func NewUser(id int64, name, email string, phone []string) *User {
	return &User{Id: id, Name: name, Email: email, Phone: phone}
}
