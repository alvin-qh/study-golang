package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListFiles(t *testing.T) {
	// 列举 `STATIC_ASSETS` 实例下的所有文件或路径信息
	files, err := ListFiles(&STATIC_ASSETS)

	assert.Nil(t, err)
	assert.Equal(t, []FileItem{
		{
			Name: "asset",
			Type: DIR,
		},
		{
			Name: "asset/01",
			Type: DIR,
		},
		{
			Name: "asset/01/static1.txt",
			Type: FILE,
		},
		{
			Name: "asset/02",
			Type: DIR,
		},
		{
			Name: "asset/02/static2.txt",
			Type: FILE,
		},
		{
			Name: "asset/static.txt",
			Type: FILE,
		},
	}, files)
}
