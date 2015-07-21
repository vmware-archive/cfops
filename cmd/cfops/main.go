package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfops"
	"github.com/pivotalservices/gtils/log"
)

const (
	logLevelEnv = "LOG_LEVEL"
)

var (
	ExitCode int
	logLevel string
	logger   log.Logger
)

func init() {
	ExitCode = cleanExitCode

	if logLevel = os.Getenv(logLevelEnv); logLevel != "" {
		log.SetLogLevel(logLevel)
		logger = log.LogFactory("cfops cli", log.Lager, os.Stdout)
		logger.Debug("log level set", log.Data{"level": logLevel})
		cfops.SetLogger(logger)
	}
}

func main() {
	app := NewApp()
	app.Run(os.Args)
	os.Exit(ExitCode)
}

// NewApp creates a new cli app
func NewApp() *cli.App {

	app := cli.NewApp()
	app.Name = "cfops"
	app.Usage = "Cloud Foundry Operations tool for IaaS installation, deployment, and management automation"
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:  "logLevel",
			Value: "info",
			Usage: "log level: debug, info, error or fatal",
		},
	)
	app.Commands = append(app.Commands, []cli.Command{
		backupCli,
		restoreCli,
	}...)

	return app
}
