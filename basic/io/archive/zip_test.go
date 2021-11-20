package archive

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	Z_ARCHIVE_FILE   = "test.zip"
	Z_UNARCHIVE_PATH = "unarchive"
)

// 通过 zip 算法，对文件进行压缩
// 由于 zip 算法自带归档结构，所以多个文件可以直接压缩到一个归档文件中
func TestArchiveWithZip(t *testing.T) {
	defer func() {
		os.Remove(Z_ARCHIVE_FILE)
		os.RemoveAll(Z_UNARCHIVE_PATH)
	}()

	// 创建 zip 归档文件
	z, err := NewZip(Z_ARCHIVE_FILE)
	assert.NoError(t, err)
	defer z.Close()

	// 压缩文件
	err = z.Archive(FileList)
	assert.NoError(t, err)

	// 解压缩文件
	err = z.Unarchive(Z_UNARCHIVE_PATH)
	assert.NoError(t, err)

	// 检查解压前后文件一致性
	eq, err := CheckUnarchiveFiles(Z_UNARCHIVE_PATH)
	assert.NoError(t, err)
	assert.True(t, eq)
}
