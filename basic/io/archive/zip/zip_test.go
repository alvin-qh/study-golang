package zip

import (
	"os"
	"study/basic/io/archive/common"
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

// 测试创建 `.zip` 压缩文件
func TestZip_New(t *testing.T) {
	// 创建压缩文件
	gz, err := New(Z_ARCHIVE_FILE)
	assert.Nil(t, err)

	defer func() {
		// 关闭压缩文件
		gz.Close()
		// 删除创建的压缩文件
		os.RemoveAll(Z_ARCHIVE_FILE)
	}()

	// 确认一个大小为 0 的 `.gz` 文件被创建
	fi, err := os.Stat(Z_ARCHIVE_FILE)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), fi.Size())
	assert.Equal(t, Z_ARCHIVE_FILE, fi.Name())
}

// 测试压缩文件
func TestZip_Archive(t *testing.T) {
	// 创建压缩文件
	gz, err := New(Z_ARCHIVE_FILE)
	assert.Nil(t, err)

	defer func() {
		// 关闭压缩文件
		gz.Close()
		// 删除创建的压缩文件
		os.RemoveAll(Z_ARCHIVE_FILE)
	}()

	// 压缩指定文件
	err = gz.Archive(fileList)
	assert.Nil(t, err)

	// 确认 `.gz` 文件的大小不大于 0, 表示有内容写入
	fi, err := os.Stat(Z_ARCHIVE_FILE)
	assert.Nil(t, err)
	assert.Greater(t, fi.Size(), int64(0))
}

// 测试解压缩文件
func TestZip_Unarchive(t *testing.T) {
	defer func() {
		os.Remove(Z_ARCHIVE_FILE)
		os.RemoveAll(Z_UNARCHIVE_PATH)
	}()

	// 执行压缩
	func() {
		// 创建压缩实例
		gz, err := New(Z_ARCHIVE_FILE)
		assert.Nil(t, err)

		defer gz.Close()

		// 压缩指定文件
		err = gz.Archive(fileList)
		assert.Nil(t, err)
	}()

	// 执行解压缩
	func() {
		// 打开压缩文件
		gz, err := New(Z_ARCHIVE_FILE)
		assert.Nil(t, err)

		defer gz.Close()

		// 解压缩
		err = gz.Unarchive(Z_UNARCHIVE_PATH)
		assert.Nil(t, err)
	}()

	// 判断解压缩前后文件是否一致
	eq, err := common.CheckUnarchiveFiles(Z_UNARCHIVE_PATH, fileList)
	assert.Nil(t, err)
	assert.True(t, eq)
}
