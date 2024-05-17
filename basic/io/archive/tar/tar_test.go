package tar

import (
	"os"
	"study/basic/io/archive/common"
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

// 测试创建归档文件
func TestTar_New(t *testing.T) {
	// 创建归档文件
	gz, err := New(TAR_ARCHIVE_FILE)
	assert.Nil(t, err)

	defer func() {
		// 关闭归档文件
		gz.Close()
		// 删除创建的归档文件
		os.RemoveAll(TAR_ARCHIVE_FILE)
	}()

	// 确认一个大小为 0 的 `.gz` 文件被创建
	fi, err := os.Stat(TAR_ARCHIVE_FILE)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), fi.Size())
	assert.Equal(t, TAR_ARCHIVE_FILE, fi.Name())
}

// 测试归档文件
func TestTar_Archive(t *testing.T) {
	// 创建归档文件
	gz, err := New(TAR_ARCHIVE_FILE)
	assert.Nil(t, err)

	defer func() {
		// 关闭归档文件
		gz.Close()
		// 删除创建的归档文件
		os.RemoveAll(TAR_ARCHIVE_FILE)
	}()

	// 归档指定文件
	err = gz.Archive(fileList)
	assert.Nil(t, err)

	// 确认归档文件的大小不大于 0, 表示有文件被归档
	fi, err := os.Stat(TAR_ARCHIVE_FILE)
	assert.Nil(t, err)
	assert.Greater(t, fi.Size(), int64(0))
}

// 测试释放被归档的文件
func TestTar_Unarchive(t *testing.T) {
	defer func() {
		os.Remove(TAR_ARCHIVE_FILE)
		os.RemoveAll(TAR_UNARCHIVE_PATH)
	}()

	// 执行归档
	func() {
		// 创建归档实例
		gz, err := New(TAR_ARCHIVE_FILE)
		assert.Nil(t, err)

		defer gz.Close()

		// 归档指定文件
		err = gz.Archive(fileList)
		assert.Nil(t, err)
	}()

	// 释放归档文件
	func() {
		// 打开归档文件
		gz, err := New(TAR_ARCHIVE_FILE)
		assert.Nil(t, err)

		defer gz.Close()

		// 释放归档文件
		err = gz.Unarchive(TAR_UNARCHIVE_PATH)
		assert.Nil(t, err)
	}()

	// 判断解归档前后文件是否一致
	eq, err := common.CheckUnarchiveFiles(TAR_UNARCHIVE_PATH, fileList)
	assert.Nil(t, err)
	assert.True(t, eq)
}
