package log

import (
	"flag"
	"io"

	"github.com/pivotal-golang/lager"
)

//Data - a type representing a map
type Data map[string]interface{}

//Logger - a logger interface type
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

//LogType - a int representation of the logtype
type LogType uint

const (
	//Lager = a logtype for Lager
	Lager LogType = iota
)

//Constants for log levels
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

//AddFlags - a function to add flags to a given flag set
func AddFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(
		&minLogLevel,
		"logLevel",
		string(INFO),
		"log level: debug, info, error or fatal",
	)
}

//SetLogLevel - a function to set the log level
func SetLogLevel(level string) {
	minLogLevel = level
}

//LogFactory - a log creator
func LogFactory(name string, logType LogType, writer io.Writer) Logger {
	log := &logger{Name: name, LogLevel: minLogLevel, Writer: writer}
	if logType == Lager {
		return NewLager(log)
	}
	return NewLager(log)
}
