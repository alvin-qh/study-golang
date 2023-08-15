package logging

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// 定义日志配置信息
type LogSetting struct {
	// 要配置的 log 对象, 如果为 nil, 则配置默认全局日志对象
	Logger *log.Logger
	// 设置的日志格式
	Formatter log.Formatter
	// 设置的日志级别
	Level log.Level
	// 要添加的日志拦截器
	Hooks []log.Hook
}

// 初始化日志
//
// 参数:
//   - `s`: 日志配置信息对象
//
// 返回:
//
//	日志对象
func LogInit(s *LogSetting) *log.Logger {
	// 获取要设置的 log 对象
	logger := s.Logger
	if logger == nil {
		// 如果未设置 log 对象, 则使用默认的标准日志对象
		logger = log.StandardLogger()
	}

	// 设置日志的格式
	withFormatter(logger, s.Formatter)
	// 设置日志的级别
	withLevel(logger, s.Level)
	// 设置日志的 Hooks
	withHooks(logger, s.Hooks)

	// 设置默认的输出为系统标准输出
	log.SetOutput(os.Stdout)

	return logger
}

// 设置日志的格式
func withFormatter(logger *log.Logger, f log.Formatter) {
	if f == nil {
		// 如果未设置日志格式, 则使用 `TextFormatter` 作为默认格式
		f = &log.TextFormatter{
			ForceColors:      true,
			TimestampFormat:  time.RFC3339,
			DisableTimestamp: false,
		}
	}
	logger.SetFormatter(f)
}

// 设置日志级别
func withLevel(logger *log.Logger, l log.Level) {
	if l == 0 {
		// 如果未设置日志级别, 则使用 `InfoLevel` 作为默认级别
		l = log.InfoLevel
	}
	logger.SetLevel(l)
}

// 设置日志拦截器
func withHooks(logger *log.Logger, hs []log.Hook) {
	if len(hs) > 0 {
		// 添加设置的日志拦截器
		for _, h := range hs {
			logger.AddHook(h)
		}
	}
}
