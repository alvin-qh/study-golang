package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"study-gin/core/conf"
	"time"

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
func Init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: time.RFC3339,
		NoColors:        true,
		CallerFirst:     true,
		ShowFullLevel:   true,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			return fmt.Sprintf(" (%v:%v)", resolveFile(f.File), f.Line)
		},
	})

	level, err := log.ParseLevel(conf.Default(conf.Config.Logger.Level, "DEBUG"))
	if err != nil {
		log.Fatalf("invalid log level in config %v", conf.Config.Logger.Level)
	}

	log.SetLevel(level)
	log.SetReportCaller(conf.Config.Logger.ShowCaller)

	log.AddHook(newRollingFileHook(
		conf.Config.Logger.File,
		withMaxSize(conf.Default(conf.Config.Logger.MaxSize, 100)),
		withCompress(conf.Config.Logger.Compress),
		withMaxAge(conf.Config.Logger.MaxAge),
		withMaxBackups(conf.Config.Logger.MaxBackups),
	))
}
