package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// 配置日志
func Setup() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
}
