package assertion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 判断文件或路径是否存在
func FileOrPathExist(t *testing.T, path string) {
	_, err := os.Stat(path)
	assert.Nilf(t, err, "file %s not exist", path)
}

// 跨平台的路径比较
func EqualPath(t *testing.T, p1 string, p2 string) {
	p1 = filepath.Clean(p1)
	p1 = strings.ReplaceAll(p1, "\\", "/")

	p2 = filepath.Clean(p2)
	p2 = strings.ReplaceAll(p2, "\\", "/")

	assert.Equal(t, p1, p2)
}

// 跨平台判读路径是否以指定字符串结尾
func PathEndsWith(t *testing.T, p string, suffix string) bool {
	return strings.HasSuffix(filepath.Clean(p), filepath.Clean(suffix))
}
