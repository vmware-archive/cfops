package cfbackup

import (
	"os"

	"github.com/pivotalservices/gtils/log"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
)

var logger = getLogger("info")

func Logger() log.Logger {
	return logger
}

func getLogger(minLogLevel string) log.Logger {
	log.SetLogLevel(minLogLevel)
	return log.LogFactory("TestLogger", log.Lager, os.Stdout)
}
