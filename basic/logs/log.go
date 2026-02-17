package logs

import (
	"errors"
	"io"
	"log"
	"sync"
)

// 定义日志级别类型
type LogLevel int

// 定义日志级别常量枚举
const (
	LEVEL_DEBUG LogLevel = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
)

const (
	Ldate         = log.Ldate
	Ltime         = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile     = log.Llongfile
	Lshortfile    = log.Lshortfile
	LUTC          = log.LUTC
	Lmsgprefix    = log.Lmsgprefix
	LstdFlags     = log.LstdFlags
)

// 定义日志级别转字符串
func (l LogLevel) String() string {
	switch l {
	case LEVEL_DEBUG:
		return "DEBUG"
	case LEVEL_INFO:
		return "INFO"
	case LEVEL_WARN:
		return "WARN"
	case LEVEL_ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// 定义日志错误
var (
	ErrNoAppender = errors.New("no appender available")
)

// 定义日志内容结构体
type logContent struct {
	level  LogLevel      // 日志级别
	format string        // 日志内容
	args   []interface{} // 格式化日志内容的参数列表
}

// 定义日志结构体
type Logger struct {
	loggers map[LogLevel]([]*log.Logger) // 保持日志级别和日志对象的对应关系
	mut     sync.Mutex                   // 用于锁定 loggers 字段的互斥锁
	writeCh chan logContent              // 写入日志的 channel
	closeCh chan struct{}                // 关闭日志的 channel
}

// 创建新的日志结构体对象
func New() *Logger {
	// 构建日志结构体
	logger := &Logger{
		loggers: make(map[LogLevel]([]*log.Logger)),
		mut:     sync.Mutex{},
		writeCh: make(chan logContent, 100), // 写入 channel, 设置 1000 个缓冲
		closeCh: make(chan struct{}),
	}

	// 启动协程, 等待日志 channel, 并将日志写入规定的 appender 中
	go func() {
		// 循环, 不断从 channel 中读取日志内容, 直到 channel 被关闭
		for {
			content, ok := <-logger.writeCh // 从 channel 中读取日志内容
			if !ok {
				break
			}
			// 写入日志
			logger.writeLog(content.level, content.format, content.args...)
		}

		// 循环结束, 表示日志 channel 已被关闭, 此时关闭日志 channel, 报告日志已正确关闭
		if logger.closeCh != nil {
			close(logger.closeCh)
		}
	}()

	return logger
}

// 关闭日志
func (l *Logger) Close() {
	if l.writeCh != nil {
		defer func() {
			l.writeCh = nil
			l.closeCh = nil
		}()

		// 关闭日志读取 channel, 令协程结束
		close(l.writeCh)

		// 等待日志 channel 关闭, 表示协程完全结束
		<-l.closeCh
	}
}

// 为日志添加新的 Appender, 用于记录日志
func (l *Logger) AddNewAppender(w io.Writer, level LogLevel, flags int) {
	// 添加指定 level 的日志对象, 例如 level 为 DEBUG, 则添加 DEBUG, INFO, WARN 和 ERROR 级别的日志对象
	for ; level <= LEVEL_ERROR; level++ {
		// 获取日志级别字符串作为日志前缀
		prefix := level.String()
		// 新建日志对象
		logger := log.New(w, prefix+" ", flags)

		// 添加该 Level 的日志对象
		l.addNewAppender(logger, level)
	}
}

// 加入新的 appender 对象
func (l *Logger) addNewAppender(logger *log.Logger, level LogLevel) {
	l.mut.Lock()
	defer l.mut.Unlock()

	// 获取指定 level 的 appender 集合
	ls, ok := l.loggers[level]
	if !ok {
		// 该 level 暂无 appender, 创建新的 appender 列表
		ls = make([]*log.Logger, 0)
	}

	// 将 logger 对象加入 appender 列表
	l.loggers[level] = append(ls, logger)
}

// 将 log 数据根据 level 写入对应的 appender 中
func (l *Logger) writeLog(level LogLevel, format string, v ...interface{}) error {
	l.mut.Lock()
	defer l.mut.Unlock()

	if ls, ok := l.loggers[level]; !ok {
		// level 对应的 appender 不存在, 返回错误
		return ErrNoAppender
	} else {
		// 依次写入该 level 下所有的 appender 中
		for _, log := range ls {
			log.Printf(format, v...)
		}
	}
	return nil
}

// 根据所给的 level 和文本内容, 写入 log
func (l *Logger) Log(level LogLevel, format string, args ...interface{}) (err error) {
	defer func() {
		if e, ok := recover().(error); ok {
			err = e
		}
	}()

	// 将 log 内容发往 channel 中, 由读取 channel 的协程完成实际的 log 写入工作
	l.writeCh <- logContent{level: level, format: format, args: args}
	return nil
}

// 写入 DEBUG 级别的 log
func (l *Logger) Debug(format string, args ...interface{}) error {
	return l.Log(LEVEL_DEBUG, format, args...)
}

// 写入 INFO 级别的 log
func (l *Logger) Info(format string, args ...interface{}) error {
	return l.Log(LEVEL_INFO, format, args...)
}

// 写入 WARN 级别的 log
func (l *Logger) Warn(format string, args ...interface{}) error {
	return l.Log(LEVEL_WARN, format, args...)
}

// 写入 ERROR 级别的 log
func (l *Logger) Error(format string, args ...interface{}) error {
	return l.Log(LEVEL_ERROR, format, args...)
}
