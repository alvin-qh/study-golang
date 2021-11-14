package io

import (
	"archive/tar"
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

	// 创建一个写入 tar 文件的 Writer 对象，归档内容均是通过该 Writer 对象写入
	tw := tar.NewWriter(file)
	defer tw.Close()

	srcList := []string{ // 待归档文件列表
		"archive_test.go",
		"file_test.go",
		"io_test.go",
		"json_test.go",
		"path_test.go",
		"xml_test.go",
	}

	// 将文件列表中的文件逐一进行归档
	// 归档的基本动作为：1. 写入 归档文件头（FileInfoHeader 结构体）；2. 写入归档文件内容（[]byte）；其中，归档文件头可以从待归档文件的 Stat 状态得到
	for _, src := range srcList {
		srcFile := src

		// 归档一个文件
		func() {
			srcFile, err := os.Open(srcFile) // 打开待归档文件
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

	tw.Flush()   // 将 Writer 的内容刷新到归档文件中
	file.Close() // 关闭归档文件，归档完成

	// 从归档文件中恢复原始文件
	// 基本的动作为：1. 打开归档文件；2. 读取一个归档文件头，并根据文件头内容创建恢复文件；3. 将归档文件内容写入恢复文件中
	file, err = os.Open("./test.tar") // 打开归档文件以供恢复
	assert.NoError(t, err)

	defer os.RemoveAll(`./unarchive`)

	err = os.Mkdir(`./unarchive`, 0755) // 创建放置恢复文件的目录
	assert.NoError(t, err)

	tr := tar.NewReader(file)                                        // 从归档文件中创建 Reader 对象，用于读取归档文件
	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() { // 遍历所有的 归档文件头
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
		srcFile := src

		func() {
			srcFile, err := os.Open(srcFile)
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
