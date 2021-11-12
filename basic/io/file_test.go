package io

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试路径操作
func TestPathOperation(t *testing.T) {
	// 获取当前路径
	path, err := os.Getwd()
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(path, "/basic/io"))

	// 根据文件或路径名获取
	path, err = exec.LookPath("./file_path_test.go")
	assert.NoError(t, err)
	assert.NotEmpty(t, path)
}
