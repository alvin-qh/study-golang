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
	path, err := os.Getwd() // 获取当前路径
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(path, "/basic/io"))

	path, err = exec.LookPath("bash")
    assert.NoError(t, err)
	assert.NotEmpty(t, path)
}
