package gzip

import (
	"basic/io/archive/common"
	"basic/io/archive/tar"
	"compress/gzip"
	"os"
	"runtime"
)

// 定义结构体
type GZip struct {
	file *os.File
}

// 创建新实例
func New(gzFile string) (*GZip, error) {
	// 创建 gzip 文件
	file, err := os.OpenFile(gzFile, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	gz := &GZip{file: file}
	runtime.SetFinalizer(gz, func(gz *GZip) { gz.Close() })

	return gz, nil
}

// 关闭实例
func (gz *GZip) Close() error {
	if gz.file == nil {
		return nil
	}
	return gz.file.Close()
}

// 压缩文件
//
// 由于 gzip 算法本身不具备归档结构, 无法压缩多个文件, 所以需要 tar 的基础上进行压缩处理:
//
//   - 需要在 `tar.Writer` 基础上增加一个 `gzip.Writer`
//   - 读取 gzip 文件同理
func (gz *GZip) Archive(srcFiles []string) error {
	// 创建用于压缩的 Writer
	gw := gzip.NewWriter(gz.file)
	defer func() {
		if err := gw.Flush(); err != nil {
			gw.Close()
		}
	}()

	// 调用 tar 包的函数进行归档
	return tar.TarArchiveFiles(gw, srcFiles)
}

// 解压缩文件
func (gz *GZip) Unarchive(targetPath string) error {
	err := common.CreateDirIfNotExists(targetPath)
	if err != nil {
		return err
	}

	// 创建用于解压缩的 Reader
	gr, err := gzip.NewReader(gz.file)
	if err != nil {
		return err
	}

	// 通过 tar 包函数释放归档文件
	return tar.TarUnarchiveFile(gr, targetPath)
}
