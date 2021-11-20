package archive

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	GZ_ARCHIVE_FILE   = "test.tar.gz"
	GZ_UNARCHIVE_PATH = "unarchive"
)

// 在 tar 归档的基础上，通过 gzip 压缩算法对 tar 归档文件进行压缩，产生 tar.gz 归档文件
// 由于 gzip 算法本身不具备归档结构，无法压缩多个文件，所以在 tar 的基础上进行压缩处理，所以只需要在 tar.Writer 基础上，增加一个 gzip.Writer 即可
// 读取 gzip 压缩的归档文件同理
func TestArchiveWithGZip(t *testing.T) {
	defer func() {
		os.Remove(GZ_ARCHIVE_FILE)
		os.RemoveAll(GZ_UNARCHIVE_PATH)
	}()

	// 创建一个用于归档的 gz 对象
	gz, err := NewGZip(GZ_ARCHIVE_FILE)
	assert.NoError(t, err)

	// 归档指定文件
	err = gz.Archive(FileList)
	assert.NoError(t, err)

	gz.Close()

	// 创建一个用于恢复归档的 tar 对象
	gz, err = NewGZip(GZ_ARCHIVE_FILE)
	assert.NoError(t, err)

	// 恢复归档中的文件
	err = gz.Unarchive(GZ_UNARCHIVE_PATH)
	assert.NoError(t, err)

	// 判断归档前后文件是否一致
	eq, err := CheckUnarchiveFiles(GZ_UNARCHIVE_PATH)
	assert.NoError(t, err)
	assert.True(t, eq)
}
