package log

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	LOG_FILE_NAME = "log.log" // 保存 log 的文件
)

// 测试创建日志实例
func TestLog_New(t *testing.T) {
	defer os.Remove(LOG_FILE_NAME)

	// 创建 log 对象
	l := New()
	defer l.Close()

	// 创建写入 log 的文件
	f, err := os.Create(LOG_FILE_NAME)
	assert.Nil(t, err)

	// 设定 log 写入的内容, 包括 日期、时间和文件名
	flag := log.Ldate | log.Ltime | log.Lshortfile

	// 增加文件和标准输出流两个 appender
	l.AddNewAppender(f, LEVEL_DEBUG, flag)
	l.AddNewAppender(os.Stdout, LEVEL_DEBUG, flag)

	// 输出各种级别的日志
	l.Debug("Test Debug Log")
	l.Info("Test Info Log")
	l.Warn("Test Warn Log")
	l.Error("Test Error Log")

	// 关闭日志对象, 等待所有日志都处理完毕
	l.Close()
	f.Close()

	// 验证日志写入情况

	// 打开日志写入文件
	f, err = os.Open(LOG_FILE_NAME)
	assert.Nil(t, err)

	br := bufio.NewReader(f)

	// 逐行验证日志写入情况

	line, _, err := br.ReadLine()
	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(string(line), "DEBUG"))
	assert.True(t, strings.HasSuffix(string(line), "Test Debug Log"))

	line, _, err = br.ReadLine()
	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(string(line), "INFO"))
	assert.True(t, strings.HasSuffix(string(line), "Test Info Log"))

	line, _, err = br.ReadLine()
	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(string(line), "WARN"))
	assert.True(t, strings.HasSuffix(string(line), "Test Warn Log"))

	line, _, err = br.ReadLine()
	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(string(line), "ERROR"))
	assert.True(t, strings.HasSuffix(string(line), "Test Error Log"))

	f.Close()
}
