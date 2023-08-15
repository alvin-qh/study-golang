package main

import (
	"logrus/logging"

	log "github.com/sirupsen/logrus"
)

func main() {
	logging.LogInit(&logging.LogSetting{
		Level: log.DebugLevel,
	})
	log.WithField("name", "ball").WithField("say", "hi").Info("info log")
}
