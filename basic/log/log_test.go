package log

import (
	"bytes"
	"io"
	logger "log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试创建日志实例
func TestLog_New(t *testing.T) {
	log := New()

	buf := bytes.NewBuffer(make([]byte, 0))

	// 设定日志输出标记
	fl := logger.Ldate | logger.Ltime | logger.Lshortfile
	log.AddNewAppender(os.Stdout, LEVEL_DEBUG, fl)

	// 将内存缓冲作为日志输出
	log.AddNewAppender(buf, LEVEL_DEBUG, fl)

	// 输出各种级别的日志
	log.Debug("Test Debug Log")
	log.Info("Test Info Log")
	log.Warn("Test Warn Log")
	log.Error("Test Error Log")

	log.Close()

	lines := strings.Split(buf.String(), "\n")
	assert.True(t, strings.HasPrefix(lines[0], "DEBUG"))
	assert.True(t, strings.HasSuffix(lines[0], "Test Debug Log"))

	assert.True(t, strings.HasPrefix(lines[1], "INFO"))
	assert.True(t, strings.HasSuffix(lines[1], "Test Info Log"))

	assert.True(t, strings.HasPrefix(lines[2], "WARN"))
	assert.True(t, strings.HasSuffix(lines[2], "Test Warn Log"))

	assert.True(t, strings.HasPrefix(lines[3], "ERROR"))
	assert.True(t, strings.HasSuffix(lines[3], "Test Error Log"))
}

const (
	LOG_FILE_NAME = "log.log" // 保存 log 的文件
)

// 测试将日志写入文件
func TestLog_WithFile(t *testing.T) {
	defer os.Remove(LOG_FILE_NAME)

	f, err := os.Create(LOG_FILE_NAME)
	assert.Nil(t, err)

	defer f.Close()

	log := New()

	// 设定日志输出标记
	fl := logger.Ldate | logger.Ltime | logger.Lshortfile
	log.AddNewAppender(os.Stdout, LEVEL_DEBUG, fl)

	// 将日志输出到文件
	log.AddNewAppender(f, LEVEL_DEBUG, fl)

	// 输出各种级别的日志
	log.Debug("Test Debug Log")
	log.Info("Test Info Log")
	log.Warn("Test Warn Log")
	log.Error("Test Error Log")

	log.Close()

	_, err = f.Seek(0, io.SeekStart)
	assert.Nil(t, err)

	data, err := io.ReadAll(f)
	assert.Nil(t, err)

	lines := strings.Split(string(data), "\n")
	assert.True(t, strings.HasPrefix(lines[0], "DEBUG"))
	assert.True(t, strings.HasSuffix(lines[0], "Test Debug Log"))

	assert.True(t, strings.HasPrefix(lines[1], "INFO"))
	assert.True(t, strings.HasSuffix(lines[1], "Test Info Log"))

	assert.True(t, strings.HasPrefix(lines[2], "WARN"))
	assert.True(t, strings.HasSuffix(lines[2], "Test Warn Log"))

	assert.True(t, strings.HasPrefix(lines[3], "ERROR"))
	assert.True(t, strings.HasSuffix(lines[3], "Test Error Log"))
}
