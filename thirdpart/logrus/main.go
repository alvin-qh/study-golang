package main

import (
	"errors"
	"fmt"
	"study/thirdpart/logrus/logging"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func main() {
	logging.LogInit(&logging.LogSetting{
		Level: log.TraceLevel,
		Formatter: &nested.Formatter{
			HideKeys:        true,
			TimestampFormat: fmt.Sprintf("[%s]", time.RFC3339),
			ShowFullLevel:   true,
			NoFieldsColors:  true,
			TrimMessages:    true,
			CallerFirst:     true,
		},
	}).SetReportCaller(true)

	log.Trace("This is a trace log by logrus")
	log.Debug("This is a debug log by logrus")
	log.Info("This is a info log by logrus")
	log.Warn("This is a warning log by logrus")
	log.Error("This is a error log by logrus", errors.New("test error"))
	log.Fatal("This is a fatal log by logrus")
}
