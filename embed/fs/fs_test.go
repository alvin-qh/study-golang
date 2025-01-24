package fs

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试列举 `embed.FS` 中嵌入的所有路径和文件
func TestListFiles(t *testing.T) {
	// 列举 `STATIC_ASSETS` 实例下的所有文件或路径信息
	files, err := ListFiles(&STATIC_ASSETS)

	assert.Nil(t, err)
	assert.Equal(t, []FileItem{
		{Name: "asset", Type: DIR},
		{Name: "asset/01", Type: DIR},
		{Name: "asset/01/static1.txt", Type: FILE},
		{Name: "asset/02", Type: DIR},
		{Name: "asset/02/static2.txt", Type: FILE},
		{Name: "asset/static.txt", Type: FILE},
	}, files)
}

// 测试从 `embed.FS` 实例中读取指定文件的内容
func TestReadEmbedFileFromFS(t *testing.T) {
	fileNames := []string{
		"asset/static.txt",
		"asset/01/static1.txt",
		"asset/02/static2.txt",
	}

	expectContent := []string{
		"asset/static.txt file\n",
		"asset/01/static1.txt file\n",
		"asset/02/static2.txt file\n",
	}

	for n, name := range fileNames {
		// 根据嵌入式文件的路径读取文件内容
		data, err := STATIC_ASSETS.ReadFile(name)
		assert.Nil(t, err)
		assert.Equal(t, expectContent[n], string(data))
	}
}

// 测试从 `embed.FS` 实例中读取指定文件的内容
func TestOperateEmbedFileFromFS(t *testing.T) {
	fileNames := []string{
		"asset/static.txt",
		"asset/01/static1.txt",
		"asset/02/static2.txt",
	}

	expectContent := []string{
		"asset/static.txt file\n",
		"asset/01/static1.txt file\n",
		"asset/02/static2.txt file\n",
	}

	// 打开 `embed.FS` 中的指定路径的文件, 并确认文件内容
	assertFile := func(fileName string, index int) {
		// 根据嵌入式文件的路径读取文件内容
		f, err := STATIC_ASSETS.Open(fileName)
		assert.Nil(t, err)

		defer f.Close()

		// 获取嵌入式文件的信息
		stat, err := f.Stat()
		assert.Nil(t, err)

		assert.Greater(t, stat.Size(), int64(1))
		assert.False(t, stat.IsDir())
		assert.Equal(t, filepath.Base(fileName), stat.Name())

		buf := make([]byte, stat.Size())

		// 读取文件内容
		n, err := f.Read(buf)
		assert.Nil(t, err)

		assert.Equal(t, n, len(buf))
		assert.Equal(t, expectContent[index], string(buf))
	}

	// 读取集合中列出的所有嵌入式文件
	for n, name := range fileNames {
		assertFile(name, n)
	}
}
