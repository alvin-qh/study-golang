package gzip

import (
	"os"
	"study/basic/io/archive/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	GZ_ARCHIVE_FILE   = "test.tar.gz"
	GZ_UNARCHIVE_PATH = "unarchive"
)

var (
	// 待压缩文件列表
	fileList = []string{
		"gzip_test.go",
		"gzip.go",
	}
)

// 测试创建 `.gz` 压缩文件
func TestGZip_New(t *testing.T) {
	// 创建压缩文件
	gz, err := New(GZ_ARCHIVE_FILE)
	assert.Nil(t, err)

	defer func() {
		// 关闭压缩文件
		gz.Close()
		os.RemoveAll(GZ_ARCHIVE_FILE)
	}()

	// 确认一个大小为 0 的 `.gz` 文件被创建
	fi, err := os.Stat(GZ_ARCHIVE_FILE)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), fi.Size())
	assert.Equal(t, GZ_ARCHIVE_FILE, fi.Name())
}

// 测试压缩文件
func TestGZip_Archive(t *testing.T) {
	// 创建压缩文件
	gz, err := New(GZ_ARCHIVE_FILE)
	assert.Nil(t, err)

	defer func() {
		// 关闭压缩文件
		gz.Close()
		// 删除创建的压缩文件
		os.RemoveAll(GZ_ARCHIVE_FILE)
	}()

	// 压缩指定文件
	err = gz.Archive(fileList)
	assert.Nil(t, err)

	// 确认 `.gz` 文件的大小不大于 0, 表示有内容写入
	fi, err := os.Stat(GZ_ARCHIVE_FILE)
	assert.Nil(t, err)
	assert.Greater(t, fi.Size(), int64(0))
}

// 测试解压缩文件
func TestGZip_Unarchive(t *testing.T) {
	defer func() {
		os.Remove(GZ_ARCHIVE_FILE)
		os.RemoveAll(GZ_UNARCHIVE_PATH)
	}()

	// 执行压缩
	{
		// 创建压缩实例
		gz, err := New(GZ_ARCHIVE_FILE)
		assert.Nil(t, err)

		defer gz.Close()

		// 压缩指定文件
		err = gz.Archive(fileList)
		assert.Nil(t, err)
	}

	// 执行解压缩
	{
		// 打开压缩文件
		gz, err := New(GZ_ARCHIVE_FILE)
		assert.Nil(t, err)

		defer gz.Close()

		// 解压缩
		err = gz.Unarchive(GZ_UNARCHIVE_PATH)
		assert.Nil(t, err)
	}

	// 判断解压缩前后文件是否一致
	eq, err := common.CheckUnarchiveFiles(GZ_UNARCHIVE_PATH, fileList)
	assert.Nil(t, err)
	assert.True(t, eq)
}
