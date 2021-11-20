package archive

import (
	"archive/tar"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

// 归档文件结构体
type Tar struct {
	file *os.File
}

// 创建一个新的 Tar 对象
func NewTar(tarFile string) (*Tar, error) {
	// 创建用于归档的 tar 文件
	file, err := os.OpenFile(tarFile, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	tar := &Tar{file: file}
	runtime.SetFinalizer(tar, func(tar *Tar) { tar.Close() })

	return tar, nil
}

// 关闭一个 Tar 对象
func (t *Tar) Close() error {
	if t.file == nil {
		return nil
	}
	return t.file.Close()
}

// 打包文件
func (t *Tar) Archive(srcFiles []string) error {
	return tarArchiveFiles(t.file, srcFiles)
}

// 归档文件
func tarArchiveFiles(w io.Writer, srcFiles []string) error {
	// 创建一个写入 tar 文件的 Writer 对象，归档内容均是通过该 Writer 对象写入
	tw := tar.NewWriter(w)
	defer func() {
		if err := tw.Flush(); err == nil {
			tw.Close()
		}
	}()

	// 将文件列表中的文件逐一进行归档
	for _, src := range srcFiles {
		if err := tarArchiveEachFile(tw, src); err != nil {
			return err
		}
	}
	return nil
}

// 归档一个文件
// 归档的基本动作为：1. 写入 归档文件头（FileInfoHeader 结构体）；2. 写入归档文件内容（[]byte）；其中，归档文件头可以从待归档文件的 Stat 状态得到
func tarArchiveEachFile(tw *tar.Writer, filename string) error {
	// 打开待归档文件
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取待归档文件状态
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// 从待归档文件状态中生成 归档文件头
	hdr, err := tar.FileInfoHeader(stat, "")
	if err != nil {
		return err
	}

	// 在归档文件中写入 归档文件头
	err = tw.WriteHeader(hdr)
	if err != nil {
		return err
	}

	// 在归档文件中写入待归档文件内容
	_, err = io.Copy(tw, file)
	return err
}

// 恢复归档文件
func (t *Tar) Unarchive(targetPath string) error {
	err := createDirIfNotExists(targetPath)
	if err != nil {
		return err
	}

	return tarUnarchiveFile(t.file, targetPath)
}

// 从归档文件中恢复被归档的文件
func tarUnarchiveFile(r io.Reader, targetPath string) error {
	// 从归档文件中创建 Reader 对象，用于读取归档文件
	tr := tar.NewReader(r)

	// 遍历所有的 归档文件头
	for hdr, err := tr.Next(); ; hdr, err = tr.Next() {
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if hdr == nil {
			break
		}

		// 获取被恢复文件属性
		fi := hdr.FileInfo()

		// 从归档文件头中获得目标文件名称
		if err = tarUnarchiveEachFile(tr, filepath.Join(targetPath, fi.Name()), fi.Mode()); err != nil {
			return err
		}
	}
	return nil
}

// 恢复一个归档文件
func tarUnarchiveEachFile(tr *tar.Reader, filename string, mode fs.FileMode) error {
	// 创建恢复文件
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将数据恢复到文件中
	_, err = io.Copy(file, tr)
	return err
}
