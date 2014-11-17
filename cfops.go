package main

import (
	"flag"
	"github.com/cloudfoundry/gosteno"
	"github.com/pivotalservices/cfops/backup"
	"github.com/pivotalservices/cfops/cli"
	"github.com/pivotalservices/cfops/install"
	"github.com/pivotalservices/cfops/start"
	"github.com/pivotalservices/cfops/system"
	"os"
	"strings"
)

var (
	logFilePath    = flag.String("logFile", "", "The CFOPS log file, defaults to STDOUT")
	configFilePath = flag.String("configFile", "config/cfops.json", "Location of the CFOPS config json file")
	debug          = flag.Bool("debug", false, "Debug logging")
)

var Logger *gosteno.Logger

func init() {
	// c := &gosteno.Config{
	// 	Sinks: []gosteno.Sink{
	// 		gosteno.NewIOSink(os.Stdout),
	// 	},
	// 	Level:     gosteno.LOG_INFO,
	// 	Codec:     gosteno.NewJsonPrettifier(gosteno.EXCLUDE_DATA),
	// 	EnableLOC: true,
	// }
	// if len(*logFilePath) > 0 {
	// 	c.Sinks = append(c.Sinks, gosteno.NewFileSink(*logFilePath))
	// }

	// gosteno.Init(c)
}

// To get the base foundation configuration see the Pivotal CF Data Collector @
// https://docs.google.com/a/pivotal.io/spreadsheets/d/1XHKSrJiQIu5MWGpPYWbMY8M09eqe-GV8MQsl_mqw1RM/edit#gid=0
func main() {
	_, logger := parseConfig(*debug, *configFilePath, *logFilePath)
	flag.Parse()

	commandFactory := cli.NewCommandFactory(logger)

	commandRunner := system.OSCommandRunner{}
	commandRunner.Logger = logger

	start.New(commandFactory)
	install.New(commandFactory)
	backup.New(commandFactory, commandRunner)

	app := cli.NewApp(commandFactory)

	app.RunAndExitOnError()
}

type Config struct {
	Index uint
}

func parseConfig(debug bool, configFile, logFilePath string) (Config, *gosteno.Logger) {
	config := Config{}
	// err := config.ReadConfigInto(&config, configFile)
	// if err != nil {
	// 	panic(err)
	// }

	logger := NewLogger(debug, logFilePath, "cfops", config)
	logger.Info("Startup: Setting up the CFOPS profiler")

	return config, logger
}

func NewLogger(verbose bool, logFilePath, name string, config Config) *gosteno.Logger {
	level := gosteno.LOG_INFO

	if verbose {
		level = gosteno.LOG_DEBUG
	}

	loggingConfig := &gosteno.Config{
		Sinks:     make([]gosteno.Sink, 1),
		Level:     level,
		Codec:     gosteno.NewJsonCodec(),
		EnableLOC: true}

	if strings.TrimSpace(logFilePath) == "" {
		loggingConfig.Sinks[0] = gosteno.NewIOSink(os.Stdout)
	} else {
		loggingConfig.Sinks[0] = gosteno.NewFileSink(logFilePath)
	}

	gosteno.Init(loggingConfig)
	logger := gosteno.NewLogger(name)
	logger.Debugf("Component %s in debug mode!", name)
	setGlobalLogger(logger)
	return logger
}

func setGlobalLogger(logger *gosteno.Logger) {
	Logger = logger
}
