package log

import (
	"flag"
	"fmt"
	"io"

	"github.com/pivotal-golang/lager"
)

type Data map[string]interface{}

type Logger interface {
	Debug(message string, data ...Data)
	Info(message string, data ...Data)
	Error(message string, err error, data ...Data)
	Fatal(message string, err error, data ...Data)
}

type logger struct {
	lager.Logger
	LogLevel string
	Name     string
	Writer   io.Writer
}

type LogType uint

const (
	Lager LogType = iota
)

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
)

var (
	log         *logger
	minLogLevel string
)

func init() {
	AddFlags(flag.CommandLine)
	flag.Parse()
}

func AddFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(
		&minLogLevel,
		"logLevel",
		string(INFO),
		"log level: debug, info, error or fatal",
	)
}

func SetLogLevel(level string) {
	minLogLevel = level
}

func LogFactory(name string, logType LogType, writer io.Writer) Logger {
	log := &logger{Name: name, LogLevel: minLogLevel, Writer: writer}
	if logType == Lager {
		return NewLager(log)
	}
	return NewLager(log)
}
