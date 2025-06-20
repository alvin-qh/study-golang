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

// 定义日志文件的最大尺寸, 超出此尺寸后会自动建立新的日志文件
//
// 参数:
//   - `maxSize` (`int`): 最大日志长度
//
// 返回:
//   - `rollingFileHookOption`: 可选参数对象
func withMaxSize(maxSize int) rollingFileHookOption {
	return func(option *lumberjack.Logger) {
		option.MaxSize = maxSize
	}
}

// 定义最大备份日志文件数量, 超出此数量, 则最早的备份日志文件会被删除
//
// 参数:
//   - `maxBackups` (`int`): 最大备份文件数量
//
// 返回:
//   - `rollingFileHookOption`: 可选参数对象
func withMaxBackups(maxBackups int) rollingFileHookOption {
	return func(option *lumberjack.Logger) {
		option.MaxBackups = maxBackups
	}
}

// 定义备份日志文件保存的最长时间 (单位为 天), 超出此时间的备份日志文件会被删除
//
// 参数:
//   - `maxAge` (`int`): 最大日志备份天数
//
// 返回:
//   - `rollingFileHookOption`: 可选参数对象
func withMaxAge(maxAge int) rollingFileHookOption {
	return func(option *lumberjack.Logger) {
		option.MaxAge = maxAge
	}
}

// 是否启用备份压缩, 如果为 `true`, 则备份的日志文件会被压缩 (gzip)
//
// 参数:
//   - `compress` (`bool`): 是否启用日志备份压缩
//
// 返回:
//   - `rollingFileHookOption`: 可选参数对象
func withCompress(compress bool) rollingFileHookOption {
	return func(option *lumberjack.Logger) {
		option.Compress = compress
	}
}
