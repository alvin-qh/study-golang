package logging

import (
	"testing"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func TestLoggingToFile(t *testing.T) {
	logger := LogInit(&LogSetting{
		Logger: log.New(),
		Level:  log.DebugLevel,
		Formatter: &nested.Formatter{
			HideKeys:        true,
			TimestampFormat: time.RFC3339,
			ShowFullLevel:   true,
			TrimMessages:    true,
			CallerFirst:     true,
			NoColors:        true,
		},
	})
	logger.SetReportCaller(true)
	logger.AddHook(NewRollingFileHook("logs/test.log"))

	logger.Debug("Hello")
}
