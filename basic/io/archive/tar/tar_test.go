package tar

import (
	"os"
	"study-golang/basic/io/archive/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TAR_ARCHIVE_FILE   = "test.tar"
	TAR_UNARCHIVE_PATH = "unarchive"
)

var (
	fileList = []string{ // 待归档文件列表
		"tar_test.go",
		"tar.go",
	}
)

// 测试一系列文件归档为 tar 文件以及从归档中恢复文件
func TestCreateTarFile(t *testing.T) {
	defer func() {
		os.Remove(TAR_ARCHIVE_FILE)
		os.RemoveAll(TAR_UNARCHIVE_PATH)
	}()

	// 创建一个用于归档的 tar 对象
	tar, err := New(TAR_ARCHIVE_FILE)
	assert.NoError(t, err)

	// 归档指定文件
	err = tar.Archive(fileList)
	assert.NoError(t, err)

	tar.Close()

	// 创建一个用于恢复归档的 tar 对象
	tar, err = New(TAR_ARCHIVE_FILE)
	assert.NoError(t, err)

	// 恢复归档中的文件
	err = tar.Unarchive(TAR_UNARCHIVE_PATH)
	assert.NoError(t, err)

	tar.Close()

	// 判断归档前后文件是否一致
	eq, err := common.CheckUnarchiveFiles(TAR_UNARCHIVE_PATH, fileList)
	assert.NoError(t, err)
	assert.True(t, eq)
}
