package zip

import (
	"os"
	"study-golang/basic/io/archive/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	Z_ARCHIVE_FILE   = "test.zip"
	Z_UNARCHIVE_PATH = "unarchive"
)

var (
	fileList = []string{ // 待归档文件列表
		"zip_test.go",
		"zip.go",
	}
)

// 测试通过 zip 算法对文件进行压缩
//
// 由于 zip 算法自带归档结构, 所以多个文件可以直接压缩到一个归档文件中
func TestArchiveWithZip(t *testing.T) {
	defer func() {
		os.Remove(Z_ARCHIVE_FILE)
		os.RemoveAll(Z_UNARCHIVE_PATH)
	}()

	// 创建 zip 归档文件
	z, err := New(Z_ARCHIVE_FILE)
	assert.NoError(t, err)

	// 压缩文件
	err = z.Archive(fileList)
	assert.NoError(t, err)

	z.Close()

	z, err = New(Z_ARCHIVE_FILE)
	assert.NoError(t, err)

	// 解压缩文件
	err = z.Unarchive(Z_UNARCHIVE_PATH)
	assert.NoError(t, err)

	z.Close()

	// 检查解压前后文件一致性
	eq, err := common.CheckUnarchiveFiles(Z_UNARCHIVE_PATH, fileList)
	assert.NoError(t, err)
	assert.True(t, eq)
}
