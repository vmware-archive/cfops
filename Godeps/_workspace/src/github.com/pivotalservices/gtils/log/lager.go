package log

import (
	"fmt"

	"github.com/pivotal-golang/lager"
)

//NewLager - constructor for a Logger object
func NewLager(log *logger) Logger {
	var minLagerLogLevel lager.LogLevel
	switch log.LogLevel {
	case DEBUG:
		minLagerLogLevel = lager.DEBUG
	case INFO:
		minLagerLogLevel = lager.INFO
	case ERROR:
		minLagerLogLevel = lager.ERROR
	case FATAL:
		minLagerLogLevel = lager.FATAL
	default:
		panic(fmt.Errorf("unknown log level: %s", log.LogLevel))
	}

	logger := lager.NewLogger(log.Name)
	logger.RegisterSink(lager.NewWriterSink(log.Writer, minLagerLogLevel))
	log.Logger = logger

	return log
}

func (l *logger) Debug(action string, data ...Data) {
	l.Logger.Debug(action, toLagerData(data...))
}

func (l *logger) Info(action string, data ...Data) {
	l.Logger.Info(action, toLagerData(data...))
}

func (l *logger) Error(action string, err error, data ...Data) {
	l.Logger.Error(action, err, toLagerData(data...))
}

func (l *logger) Fatal(action string, err error, data ...Data) {
	l.Logger.Fatal(action, err, toLagerData(data...))
}

func toLagerData(givenData ...Data) lager.Data {
	data := lager.Data{}

	if len(givenData) > 0 {
		for _, dataArg := range givenData {
			for key, val := range dataArg {
				data[key] = val
			}
		}
	}

	return data
}
