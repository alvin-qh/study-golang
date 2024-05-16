package logging

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bou.ke/monkey"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

// 用于测试的日志拦截器结构体
type TestingHook struct {
	// 存储日志的缓冲区
	buffer *bytes.Buffer
}

// 将日志拦截器拦截的日志内容转为字符串
func (h *TestingHook) String() string {
	return h.buffer.String()
}

// 创建新的测试日志拦截器对象
//
// 返回值:
//
//	拦截器对象
func newTestingHook() *TestingHook {
	return &TestingHook{
		buffer: bytes.NewBuffer(make([]byte, 0, 4096)),
	}
}

// 设置该日志拦截器对应的日志级别
//
// 返回值:
//
//	可以应用此拦截器的日志级别集合
func (h *TestingHook) Levels() []log.Level {
	return []log.Level{
		log.DebugLevel,
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
	}
}

// 日志拦截方法
//
// 参数:
//   - `e` 日志实体指针
//
// 返回值:
//
//	错误对象, `nil` 表示没有错误
func (h *TestingHook) Fire(e *log.Entry) error {
	s, err := e.String()
	if err != nil {
		return err
	}

	_, err = h.buffer.Write([]byte(s))
	return err
}

var (
	// 日志时间时区
	zone, _ = time.LoadLocation("Asia/Shanghai")
)

// 在每个测试前执行, 对 `time.Now()` 方法进行 patch
func TestMain(m *testing.M) {
	// 替换 `time.Now()` 方法, 返回固定时间
	guard := monkey.Patch(time.Now, func() time.Time {
		return time.Date(2023, 8, 15, 17, 13, 0, 0, zone)
	})
	// 执行完毕后, 取消 patch
	defer guard.Unpatch()

	os.Exit(m.Run())
}

// 测试 `logrus` 库自带的格式
func TestTextFormatterWithColorificOutput(t *testing.T) {
	// 实例化测试用的日志拦截器
	hook := newTestingHook()

	logger := LogInit(&LogSetting{
		Logger: log.New(),
		Level:  log.DebugLevel,
		// 使用 `nested-logrus-formatter` 格式化插件
		Formatter: &log.TextFormatter{
			// 强制输出色彩信息
			ForceColors:   true,
			DisableColors: false,
			FullTimestamp: true,
		},
		Hooks: []log.Hook{hook},
	})
	logger.Debug("")
}

// 测试 `nested-logrus-formatter` 库的日志格式化插件
func TestNestedFormatter(t *testing.T) {
	// 实例化测试用的日志拦截器
	hook := newTestingHook()

	logger := LogInit(&LogSetting{
		Logger: log.New(),
		Level:  log.DebugLevel,
		// 使用 `nested-logrus-formatter` 格式化插件
		Formatter: &nested.Formatter{
			FieldsOrder:           []string{},
			HideKeys:              true,
			NoColors:              false,
			NoFieldsColors:        false,
			TimestampFormat:       fmt.Sprintf("[%s]", time.RFC3339),
			ShowFullLevel:         true,
			NoUppercaseLevel:      false,
			NoFieldsSpace:         false,
			TrimMessages:          true,
			CallerFirst:           true,
			CustomCallerFormatter: nil,
		},
		Hooks: []log.Hook{hook},
	})
	logger.SetReportCaller(true)

	// 获取
	_, currentFilename, line, _ := runtime.Caller(0)
	logger.Debug("Hello")

	assert.Equal(t, fmt.Sprintf(
		"[2023-08-15T17:13:00+08:00] "+
			"(%s:%d study/thirdpart/logrus/logging.TestNestedFormatter)"+
			"\x1b[37m [DEBUG] \x1b[0mHello\n", currentFilename, line+1), hook.String())
}
