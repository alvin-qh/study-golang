package gzip

import (
	"compress/gzip"
	"os"
	"runtime"
	"study-golang/basic/io/archive/common"
	"study-golang/basic/io/archive/tar"
)

// GZip 归档文件结构体
type GZip struct {
	file *os.File
}

// 创建一个新的 GZip 对象
func New(gzFile string) (*GZip, error) {
	// 创建用于归档的 gzip 文件
	file, err := os.OpenFile(gzFile, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	gz := &GZip{file: file}
	runtime.SetFinalizer(gz, func(gz *GZip) { gz.Close() })

	return gz, nil
}

// 关闭一个 Tar 对象
func (gz *GZip) Close() error {
	if gz.file == nil {
		return nil
	}
	return gz.file.Close()
}

// 打包文件
func (gz *GZip) Archive(srcFiles []string) error {
	// 创建用于压缩的 Writer
	gw := gzip.NewWriter(gz.file)
	defer func() {
		if err := gw.Flush(); err != nil {
			gw.Close()
		}
	}()

	// 调用 ta.go 中的函数进行归档
	return tar.TarArchiveFiles(gw, srcFiles)
}

// 恢复归档文件
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
	return tar.TarUnarchiveFile(gr, targetPath)
}
