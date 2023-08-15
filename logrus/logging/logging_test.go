package logging

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bou.ke/monkey"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

type TestingHook struct {
	buffer *bytes.Buffer
}

func (h *TestingHook) String() string {
	return h.buffer.String()
}

func newTestingHook() *TestingHook {
	return &TestingHook{
		buffer: bytes.NewBuffer(make([]byte, 0, 4096)),
	}
}

func (h *TestingHook) Levels() []log.Level {
	return []log.Level{
		log.DebugLevel,
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
	}
}

func (h *TestingHook) Fire(e *log.Entry) error {
	s, err := e.String()
	if err != nil {
		return err
	}

	_, err = h.buffer.Write([]byte(s))
	return err
}

var (
	hook    *TestingHook
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

	// 实例化测试用的日志拦截器
	hook = newTestingHook()
	os.Exit(m.Run())
}

// 测试 `logrus` 库自带的格式
func TestTextFormatterWithColorificOutput(t *testing.T) {
	logger := LogInit(&LogSetting{
		logger: log.New(),
		Level:  log.DebugLevel,
		// 使用 `nested-logrus-formatter` 格式化插件
		Formatter: &log.TextFormatter{
			// 强制输出色彩信息
			ForceColors:   true,
			DisableColors: false,
            FullTimestamp: true,
            Time
		},
		Hooks: []log.Hook{hook},
	})
}

// 测试 `nested-logrus-formatter` 库的日志格式化插件
func TestNestedFormatter(t *testing.T) {
	logger := LogInit(&LogSetting{
		logger: log.New(),
		Level:  log.DebugLevel,
		// 使用 `nested-logrus-formatter` 格式化插件
		Formatter: &nested.Formatter{
			FieldsOrder:           []string{},
			HideKeys:              true,
			NoColors:              false,
			NoFieldsColors:        false,
			TimestampFormat:       time.RFC3339,
			ShowFullLevel:         true,
			NoUppercaseLevel:      false,
			NoFieldsSpace:         false,
			TrimMessages:          true,
			CallerFirst:           false,
			CustomCallerFormatter: nil,
		},
		Hooks: []log.Hook{hook},
	})

	logger.Debug("Hello")
	assert.Equal(t, "2023-08-15T17:13:00+08:00\x1b[37m [DEBUG] \x1b[0mHello\n", hook.String())
}
