package fs

import (
	"embed"
	"fmt"
)

var (
	// 将 `./asset` 路径下的内容嵌入为文件系统
	//go:embed asset/*
	STATIC_ASSETS embed.FS
)

// 定义表示文件类型的类型
type FileType int

// 定义文件类型枚举
const (
	FILE FileType = iota
	DIR
)

// 定义文件类型名称切片
var fileTypeStr []string = []string{"file", "dir"}

// 保存文件 (或路径) 信息的结构体
type FileItem struct {
	// 文件 (或路径) 的名称
	Name string
	// 表示是文件或路径的类型
	Type FileType
}

// 获取文件或路径类型的名称字符串
func (f *FileItem) TypeName() string {
	return fileTypeStr[f.Type]
}

// 获取文件 (或路径) 信息的字符串表达
func (f *FileItem) String() string {
	return fmt.Sprintf("%v<%v>", f.Name, f.TypeName())
}

// 列举 `embed.FS` 实例中包含的所有文件 (或路径)
func ListFiles(fs *embed.FS) ([]FileItem, error) {
	// 保存结果的切片实例
	items := make([]FileItem, 0)

	var readDir func(string) error

	// 定义递归函数, 用于遍历 `embed.FS` 对象中的所有文件或路径
	readDir = func(path string) error {
		// 读取所给路径下的文件 (或子路径)
		entries, err := STATIC_ASSETS.ReadDir(path)
		if err != nil {
			return err
		}

		// 遍历结果
		for _, entry := range entries {
			// 处理文件 (或路径) 名称, 将遍历中的当前文件 (或路径名) 和上一级路径名拼合为完整路径文件名
			// 如果上一级路径名为 `"."`, 则需要从路径中排除, `embed.FS.ReadDir` 方法不接受以 `./` 开头的路径名
			var name string
			if path == "." {
				name = entry.Name()
			} else {
				name = fmt.Sprintf("%v/%v", path, entry.Name())
			}

			if entry.IsDir() {
				// 保存当前遍历的路径项
				items = append(items, FileItem{
					Name: name,
					Type: DIR,
				})

				// 如果当前遍历的项表示路径, 则进一步读取下一级路径的内容
				err = readDir(name)
				if err != nil {
					return err
				}
			} else {
				// 保存当前遍历的文件项
				items = append(items, FileItem{
					Name: name,
					Type: FILE,
				})
			}
		}
		return nil
	}

	// 从根路径开始读取 `embed.FS` 下的所有路径和文件
	if err := readDir("."); err != nil {
		return nil, err
	}
	return items, nil
}
