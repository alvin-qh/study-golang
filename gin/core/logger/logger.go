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
func Setup(config *conf.Config) {
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

	level, err := log.ParseLevel(conf.Default(config.Logger.Level, "DEBUG"))
	if err != nil {
		panic(err)
	}

	log.SetLevel(level)
	log.SetReportCaller(config.Logger.ShowCaller)

	log.AddHook(newRollingFileHook(
		config.Logger.File,
		withMaxSize(conf.Default(config.Logger.MaxSize, 100)),
		withCompress(config.Logger.Compress),
		withMaxAge(config.Logger.MaxAge),
		withMaxBackups(config.Logger.MaxBackups),
	))
}
