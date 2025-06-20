package zip

import (
	"archive/zip"
	"basic/io/archive/common"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// Zip 归档文件结构体
type Zip struct {
	file *os.File
}

// 创建一个新的 Zip 实例
func New(zipFile string) (*Zip, error) {
	// 创建用于归档的 tar 文件
	file, err := os.OpenFile(zipFile, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	z := &Zip{file: file}
	runtime.SetFinalizer(z, func(z *Zip) { z.Close() })

	return z, nil
}

// 关闭一个 Zip 实例
func (z *Zip) Close() error {
	if z.file == nil {
		return nil
	}
	return z.file.Close()
}

// 打包文件
func (z *Zip) Archive(srcFiles []string) error {
	// 创建一个写入 zip 文件的 Writer 实例, 归档内容均是通过该 Writer 实例写入
	zw := zip.NewWriter(z.file)
	defer func() {
		if err := zw.Flush(); err == nil {
			zw.Close()
		}
	}()

	// 将文件列表中的文件逐一进行归档
	for _, src := range srcFiles {
		if err := zipArchiveEachFile(zw, src); err != nil {
			return err
		}
	}
	return nil
}

// 归档一个文件
//
// 归档的基本动作为:
//  1. 写入归档文件头 (`FileInfoHeader` 结构体);
//  2. 写入归档文件内容 (`[]byte`);
//
// 其中, 归档文件头可以从待归档文件的 `Stat` 状态得到
func zipArchiveEachFile(zw *zip.Writer, filename string) error {
	// 打开待归档文件
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取待压缩文件状态
	fi, err := file.Stat()
	if err != nil {
		return err
	}

	// 创建一个压缩文件头
	hdr, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}

	// 根据压缩文件头, 创建一个 Writer,  用于写入压缩内容
	zh, err := zw.CreateHeader(hdr)
	if err != nil {
		return err
	}
	// 也可以不使用 归档文件头, 直接通过一个字符串作为标识写入归档内容
	// zfw, err := zw.Create(srcName)

	// 将源文件压缩并写入压缩文件
	_, err = io.Copy(zh, file)
	return err
}

// 恢复归档文件
func (z *Zip) Unarchive(unarchivePath string) error {
	err := common.CreateDirIfNotExists(unarchivePath)
	if err != nil {
		return err
	}

	// 获取压缩文件状态信息
	fi, err := z.file.Stat()
	if err != nil {
		return err
	}

	// 创建压缩文件的 Reader 实例
	zr, err := zip.NewReader(z.file, fi.Size())
	if err != nil {
		return err
	}

	// 也可以直接打开压缩文件, 得到 Reader 实例
	// zr, err := zip.OpenReader(`./test.zip`)

	// 遍历压缩文件中的归档文件列表, 逐一进行解压缩
	for _, zf := range zr.File {
		if err := zipUnarchiveEachFile(zf, filepath.Join(unarchivePath, zf.Name)); err != nil {
			return err
		}
	}
	return nil
}

// 恢复一个压缩文件
func zipUnarchiveEachFile(zf *zip.File, filename string) error {
	// 创建恢复文件
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, zf.Mode())
	if err != nil {
		return err
	}
	defer file.Close()

	// 打开压缩文件中待解压的那部分
	zfr, err := zf.Open()
	if err != nil {
		return err
	}
	defer zfr.Close()

	// 将数据恢复到文件中
	_, err = io.Copy(file, zfr)
	return err
}
