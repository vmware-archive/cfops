package main

import (
	// "fmt"
	"github.com/cloudfoundry/gosteno"
	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfops/backup"
	. "github.com/pivotalservices/cfops/cli"
	"github.com/pivotalservices/cfops/install"
	"github.com/pivotalservices/cfops/start"
	"github.com/pivotalservices/cfops/system"
	"os"
	"strings"
)

var (
	logFilePath    = NewFlag("logFile", "", "The CFOPS log file, defaults to STDOUT", "CFOPS_LOG")
	configFilePath = NewFlag("configFile", "config/cfops.json", "Location of the CFOPS config json file", "CFOPS_CONFIG")
	debug          = NewFlag("debug", false, "Debug logging", "CFOPS_TRACE")
	iaas           = NewFlag("iaas, i", "aws, vsphere, vcloud, openstack", "Sets the IaaS type to target for deployment", "CFOPS_IAAS")
	lang           = NewFlag("lang, l", "en, es", "Language for the cfops cli responses", "CFOPS_LANG")
)

type Config struct {
	Index uint
}

// To get the base foundation configuration see the Pivotal CF Data Collector @
// https://docs.google.com/a/pivotal.io/spreadsheets/d/1XHKSrJiQIu5MWGpPYWbMY8M09eqe-GV8MQsl_mqw1RM/edit#gid=0
func main() {

	commandFactory := NewCommandFactory()

	commandRunner := &system.OSCommandRunner{}

	start.New(commandFactory)
	install.New(commandFactory)
	backup.New(commandFactory, commandRunner)

	globalFlags := []*Flag{logFilePath, configFilePath, debug, iaas, lang}

	app := NewApp(commandFactory, globalFlags)

	app.Before = func(c *cli.Context) error {
		for _, flag := range globalFlags {
			if flag.Type == "bool" {
				flag.Value = c.Bool(flag.Name)
				// fmt.Printf("BEFORE flag %s has value %v\n", flag.Name, flag.Value)
			} else {
				flag.Value = c.String(flag.Name)
				// fmt.Printf("BEFORE flag %s has value %s\n", flag.Name, flag.Value)
			}
		}
		_, logger := parseConfig(debug.ParseBool(), configFilePath.ParseString(), logFilePath.ParseString())

		commandRunner.SetLogger(logger)
		commandFactory.SetLogger(logger)
		return nil
	}

	app.RunAndExitOnError()
}

func parseConfig(debug bool, configFile, logFilePath string) (Config, *gosteno.Logger) {
	config := Config{}
	// err := config.LoadConfig(&config, configFile)
	// if err != nil {
	// 	panic(err)
	// }

	logger := NewLogger(debug, logFilePath, "cfops", config)
	logger.Info("Setting up CFOPS profiler")

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
		Codec:     gosteno.NewJsonPrettifier(gosteno.EXCLUDE_DATA),
		EnableLOC: true}

	if strings.TrimSpace(logFilePath) == "" {
		loggingConfig.Sinks[0] = gosteno.NewIOSink(os.Stdout)
	} else {
		loggingConfig.Sinks[0] = gosteno.NewFileSink(logFilePath)
	}

	gosteno.Init(loggingConfig)
	logger := gosteno.NewLogger(name)
	logger.Debugf("Component %s in debug mode!", name)
	return logger
}
