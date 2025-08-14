package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"study/web/gin/core/conf"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

var (
	_DIR, _ = os.Getwd()
)

func resolveFile(file string) string {
	file, err := filepath.Rel(_DIR, file)
	if err != nil {
		return ""
	}
	return file
}

// 初始化日志
func init() {
	// 设置日志输出, 默认输出为控制台
	log.SetOutput(os.Stdout)
	// 设置日志格式
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,         // 不显示 key 名称
		TimestampFormat: time.RFC3339, // 设置时间显示格式
		NoColors:        true,         // 不显示颜色
		CallerFirst:     true,         // 优先显示调用位置
		ShowFullLevel:   true,         // 输出完整的日志等级
		CustomCallerFormatter: func(f *runtime.Frame) string {
			// 设置调用位置的日志格式
			return fmt.Sprintf(" (%v:%v)", resolveFile(f.File), f.Line)
		},
	})

	// 从配置文件中解析最小日志等级
	level, err := log.ParseLevel(conf.Config.Logger.Level)
	if err != nil {
		log.Fatalf("invalid log level in config %v", conf.Config.Logger.Level)
	}
	// 设置最小日志等级
	log.SetLevel(level)
	// 要求日志中包含调用位置
	log.SetReportCaller(conf.Config.Logger.ShowCaller)

	// 查看是否运行在单元测试
	if len(conf.Config.Logger.File) > 0 {
		// 添加日志 hook, 用于将日志记录到文件中
		log.AddHook(newRollingFileHook(
			conf.Config.Logger.File,
			withMaxSize(conf.Config.Logger.MaxSize),
			withCompress(conf.Config.Logger.Compress),
			withMaxAge(conf.Config.Logger.MaxAge),
			withMaxBackups(conf.Config.Logger.MaxBackups),
		))
	}
}
