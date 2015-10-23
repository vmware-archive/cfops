package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/gtils/log"
)

const (
	logLevelEnv = "LOG_LEVEL"
)

var (
	VERSION  string
	ExitCode int
	logLevel string
	logger   log.Logger
)

func init() {
	ExitCode = cleanExitCode
}

func main() {
	app := NewApp()
	app.Run(os.Args)
	os.Exit(ExitCode)
}

// NewApp creates a new cli app
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Version = VERSION
	app.Name = "cfops"
	app.Usage = "Cloud Foundry Operations Tool"
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   "logLevel",
			Value:  "info",
			Usage:  "set the log level by setting the LOG_LEVEL environment variable",
			EnvVar: "LOG_LEVEL",
		},
	)
	app.Commands = append(app.Commands, []cli.Command{
		cli.Command{
			Name: "version",
			Action: func(c *cli.Context) {
				cli.ShowVersion(c)
			},
		},
		backupCli,
		restoreCli,
	}...)
	return app
}
