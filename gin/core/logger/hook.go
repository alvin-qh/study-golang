package logger

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 定义滚动文件 Hook
type rollingFileHook struct {
	logger *lumberjack.Logger
}

// 可变参数类型
type rollingFileHookOption = func(option *lumberjack.Logger)

// 创建滚动文件日志拦截器
//
// 参数:
//   - `filename` 日志文件名
//   - `options` 日志记录选项, 参考 `RollingFileHookOption` 选项类型
//
// 返回:
//   - 日志拦截器对象
func newRollingFileHook(filename string, options ...rollingFileHookOption) *rollingFileHook {
	opt := lumberjack.Logger{
		Filename:   filename,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}

	for _, fn := range options {
		fn(&opt)
	}

	opt.Rotate()

	return &rollingFileHook{
		logger: &opt,
	}
}

// 设置该日志拦截器对应的日志级别
//
// 返回值:
//   - 可以应用此拦截器的日志级别集合
func (h *rollingFileHook) Levels() []log.Level {
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
//   - 错误对象, `nil` 表示没有错误
func (h *rollingFileHook) Fire(e *log.Entry) error {
	fmt := e.Logger.Formatter

	c, err := fmt.Format(e)
	if err != nil {
		return err
	}

	_, err = h.logger.Write(c)
	return err
}

// 定义 `NewRollingFileHook` 函数的选项

func withMaxSize(maxSize int) rollingFileHookOption {
	return func(option *lumberjack.Logger) {
		option.MaxSize = maxSize
	}
}

func withMaxBackups(maxBackups int) rollingFileHookOption {
	return func(option *lumberjack.Logger) {
		option.MaxBackups = maxBackups
	}
}

func withMaxAge(maxAge int) rollingFileHookOption {
	return func(option *lumberjack.Logger) {
		option.MaxAge = maxAge
	}
}

func withCompress(compress bool) rollingFileHookOption {
	return func(option *lumberjack.Logger) {
		option.Compress = compress
	}
}
