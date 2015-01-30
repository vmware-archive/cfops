package cfbackup

import (
	"fmt"
	"github.com/pivotal-golang/lager"
	"os"
)

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
)

var logger = getLogger("info")

func Logger() lager.Logger {
	return logger
}

func getLogger(minLogLevel string) lager.Logger {
	var minLagerLogLevel lager.LogLevel
	switch minLogLevel {
	case DEBUG:
		minLagerLogLevel = lager.DEBUG
	case INFO:
		minLagerLogLevel = lager.INFO
	case ERROR:
		minLagerLogLevel = lager.ERROR
	case FATAL:
		minLagerLogLevel = lager.FATAL
	default:
		panic(fmt.Errorf("unknown log level: %s", minLogLevel))
	}

	logger := lager.NewLogger("TestLogger")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, minLagerLogLevel))

	return logger
}
