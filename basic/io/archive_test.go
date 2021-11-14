package io

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试一系列文件归档为 tar 文件以及从归档中恢复文件
func TestCreateTarFile(t *testing.T) {
	defer os.Remove("./test.tar")

	// 创建一个用于归档的 tar 文件
	file, err := os.Create("./test.tar")
	assert.NoError(t, err)
	defer file.Close()

	srcList := []string{ // 待归档文件列表
		"archive_test.go",
		"file_test.go",
		"io_test.go",
		"json_test.go",
		"path_test.go",
		"xml_test.go",
	}

	// 创建一个写入 tar 文件的 Writer 对象，归档内容均是通过该 Writer 对象写入
	tw := tar.NewWriter(file)
	defer tw.Close()

	// 将文件列表中的文件逐一进行归档
	// 归档的基本动作为：1. 写入 归档文件头（FileInfoHeader 结构体）；2. 写入归档文件内容（[]byte）；其中，归档文件头可以从待归档文件的 Stat 状态得到
	for _, src := range srcList {
		fn := src

		// 归档一个文件
		func() {
			srcFile, err := os.Open(fn) // 打开待归档文件
			assert.NoError(t, err)

			defer srcFile.Close()

			stat, err := srcFile.Stat() // 获取待归档文件状态
			assert.NoError(t, err)

			hdr, err := tar.FileInfoHeader(stat, "") // 从待归档文件状态中生成 归档文件头
			assert.NoError(t, err)

			err = tw.WriteHeader(hdr) // 在归档文件中写入 归档文件头
			assert.NoError(t, err)

			_, err = io.Copy(tw, srcFile) // 在归档文件中写入待归档文件内容
			assert.NoError(t, err)
		}()
	}

	tw.Flush() // 将 Writer 的内容刷新到归档文件中
	tw.Close()

	file.Close() // 关闭归档文件，归档完成

	// 从归档文件中恢复原始文件
	// 基本的动作为：1. 打开归档文件；2. 读取一个归档文件头，并根据文件头内容创建恢复文件；3. 将归档文件内容写入恢复文件中
	defer os.RemoveAll(`./unarchive`)

	err = os.Mkdir(`./unarchive`, 0755) // 创建放置恢复文件的目录
	assert.NoError(t, err)

	file, err = os.Open("./test.tar") // 打开归档文件以供恢复
	assert.NoError(t, err)

	tr := tar.NewReader(file) // 从归档文件中创建 Reader 对象，用于读取归档文件

	for hdr, err := tr.Next(); err != io.EOF && hdr != nil; hdr, err = tr.Next() { // 遍历所有的 归档文件头
		fi := hdr.FileInfo()
		fn := filepath.Join(`./unarchive`, fi.Name()) // 从归档文件头中获得目标文件名称

		// 恢复一个归档文件
		func() {
			dstFile, err := os.Create(fn) // 创建恢复文件
			assert.NoError(t, err)

			defer dstFile.Close()

			io.Copy(dstFile, tr) // 将数据恢复到文件中
		}()
	}

	// 确认归档前的文件和归档恢复后的文件数量和内容一致
	for _, src := range srcList {
		fn := src

		func() {
			srcFile, err := os.Open(fn)
			assert.NoError(t, err)
			defer srcFile.Close()

			arcFile, err := os.Open(filepath.Join(`./unarchive`, src))
			assert.NoError(t, err)
			defer arcFile.Close()

			srcData, err := io.ReadAll(srcFile)
			assert.NoError(t, err)

			arcData, err := io.ReadAll(arcFile)
			assert.NoError(t, err)

			assert.Equal(t, arcData, srcData)
		}()
	}
}

// 在 tar 归档的基础上，通过 gzip 压缩算法对 tar 归档文件进行压缩，产生 tar.gz 归档文件
// 由于 gzip 算法本身不具备归档结构，无法压缩多个文件，所以在 tar 的基础上进行压缩处理，所以只需要在 tar.Writer 基础上，增加一个 gzip.Writer 即可
// 读取 gzip 压缩的归档文件同理
func TestArchiveWithGZip(t *testing.T) {
	defer os.Remove("./test.tar.gz")

	file, err := os.Create("./test.tar.gz")
	assert.NoError(t, err)
	defer file.Close()

	srcList := []string{ // 待归档文件列表
		"archive_test.go",
		"file_test.go",
		"io_test.go",
		"json_test.go",
		"path_test.go",
		"xml_test.go",
	}

	// 创建一个写入 gzip 文件的 Writer 对象，对写入内容进行压缩
	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw) // 创建写入 tar 文件的 Writer，写入到 gzip.Writer 中
	defer tw.Close()

	for _, src := range srcList { // 对文件进行归档
		fn := src

		// 归档一个文件
		func() {
			srcFile, err := os.Open(fn) // 打开待归档文件
			assert.NoError(t, err)

			defer srcFile.Close()

			stat, err := srcFile.Stat() // 获取待归档文件状态
			assert.NoError(t, err)

			hdr, err := tar.FileInfoHeader(stat, "") // 从待归档文件状态中生成 归档文件头
			assert.NoError(t, err)

			err = tw.WriteHeader(hdr) // 在归档文件中写入 归档文件头
			assert.NoError(t, err)

			_, err = io.Copy(tw, srcFile) // 在归档文件中写入待归档文件内容
			assert.NoError(t, err)
		}()
	}

	tw.Flush() // 结束归档，将 gzip 和 tar 的 Writer 进行刷新，并关闭归档文件
	tw.Close()

	gw.Flush()
	gw.Close()

	file.Close()

	// 从压缩后的归档文件中恢复源文件
	defer os.RemoveAll(`./unarchive`)

	err = os.Mkdir(`./unarchive`, 0755) // 创建放置恢复文件的目录
	assert.NoError(t, err)

	file, err = os.Open("./test.tar.gz") // 打开归档文件以供恢复
	assert.NoError(t, err)

	gr, err := gzip.NewReader(file) // 对归档文件先进行一个 gzip.Reader 的包装，对其进行解压缩
	assert.NoError(t, err)

	tr := tar.NewReader(gr) // 在 gzip.Reader 基础上，通过 tar.Reader 进行归档内容读取

	for hdr, err := tr.Next(); err != io.EOF && hdr != nil; hdr, err = tr.Next() { // 遍历所有的 归档文件头
		fi := hdr.FileInfo()
		fn := filepath.Join(`./unarchive`, fi.Name()) // 从归档文件头中获得目标文件名称

		// 恢复一个归档文件
		func() {
			dstFile, err := os.Create(fn) // 创建恢复文件
			assert.NoError(t, err)

			defer dstFile.Close()

			io.Copy(dstFile, tr) // 将数据恢复到文件中
		}()
	}

	// 确认归档前的文件和归档恢复后的文件数量和内容一致
	for _, src := range srcList {
		srcName := src

		func() {
			srcFile, err := os.Open(srcName)
			assert.NoError(t, err)
			defer srcFile.Close()

			arcFile, err := os.Open(filepath.Join(`./unarchive`, src))
			assert.NoError(t, err)
			defer arcFile.Close()

			srcData, err := io.ReadAll(srcFile)
			assert.NoError(t, err)

			arcData, err := io.ReadAll(arcFile)
			assert.NoError(t, err)

			assert.Equal(t, arcData, srcData)
		}()
	}
}

// 通过 zip 算法，对文件进行压缩
// 由于 zip 算法自带归档结构，所以多个文件可以直接压缩到一个归档文件中
func TestArchiveWithZip(t *testing.T) {
	defer os.Remove("./test.zip")

	// 创建 zip 归档文件
	file, err := os.Create("./test.zip")
	assert.NoError(t, err)
	defer file.Close()

	srcList := []string{ // 待归档文件列表
		"archive_test.go",
		"file_test.go",
		"io_test.go",
		"json_test.go",
		"path_test.go",
		"xml_test.go",
	}

	zw := zip.NewWriter(file) // 通过 zip.Writer 包装 file 对象，进行压缩写入
	defer zw.Close()

	for _, src := range srcList { // 对文件进行归档
		srcName := src

		// 归档一个文件
		func() {
			srcFile, err := os.Open(srcName) // 打开待压缩文件
			assert.NoError(t, err)

			defer srcFile.Close()

			stat, err := srcFile.Stat() // 获取待压缩文件状态
			assert.NoError(t, err)

			hdr, err := zip.FileInfoHeader(stat) // 创建一个压缩文件头
			assert.NoError(t, err)

			zfw, err := zw.CreateHeader(hdr) // 根据压缩文件头，创建一个 Writer， 用于写入压缩内容
			assert.NoError(t, err)

			// zfw, err := zw.Create(srcName) // 也可以不使用 归档文件头，直接通过一个字符串作为标识写入归档内容

			_, err = io.Copy(zfw, srcFile) // 将源文件压缩并写入压缩文件
			assert.NoError(t, err)
		}()
	}

	zw.Flush() // 将 zip Writer 的内容刷新到文件中，并关闭 zip 压缩文件
	zw.Close()
	file.Close()

	// 解压并恢复归档文件
	defer os.RemoveAll(`./unarchive`)

	err = os.Mkdir(`./unarchive`, 0755)
	assert.NoError(t, err)

	file, err = os.Open(`./test.zip`) // 打开 zip 压缩文件
	assert.NoError(t, err)

	stat, err := file.Stat() // 获取压缩文件状态信息
	assert.NoError(t, err)

	zr, err := zip.NewReader(file, stat.Size()) // 创建压缩文件的 Reader 对象
	assert.NoError(t, err)

	// zr, err := zip.OpenReader(`./test.zip`) // 也可以直接打开压缩文件，得到 Reader 对象

	// 遍历压缩文件中的归档文件列表，逐一进行解压缩
	for _, zf := range zr.File {
		fn := filepath.Join(`./unarchive`, zf.Name) // 生成解压缩目标文件名

		func() {
			df, err := os.Create(fn) // 创建解压缩文件
			assert.NoError(t, err)
			defer df.Close()

			zfr, err := zf.Open() // 打开压缩文件中待解压的那部分
			assert.NoError(t, err)
			defer zfr.Close()

			io.Copy(df, zfr) // 将压缩内容解压后写入解压缩文件
		}()
	}

	// 确认归档前的文件和归档恢复后的文件数量和内容一致
	for _, src := range srcList {
		srcName := src

		func() {
			srcFile, err := os.Open(srcName)
			assert.NoError(t, err)
			defer srcFile.Close()

			arcFile, err := os.Open(filepath.Join(`./unarchive`, src))
			assert.NoError(t, err)
			defer arcFile.Close()

			srcData, err := io.ReadAll(srcFile)
			assert.NoError(t, err)

			arcData, err := io.ReadAll(arcFile)
			assert.NoError(t, err)

			assert.Equal(t, arcData, srcData)
		}()
	}
}
