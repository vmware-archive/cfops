package main

import (
	"github.com/cloudfoundry/gosteno"
	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfops/app"
	"github.com/pivotalservices/cfops/backup"
	"github.com/pivotalservices/cfops/install"
	"os"
)

var (
	installer *install.Installer
	backuper  *backup.Backuper
)

func init() {
	c := &gosteno.Config{
		Sinks: []gosteno.Sink{
			gosteno.NewFileSink("./cfops.log"),
			gosteno.NewIOSink(os.Stdout),
		},
		Level:     gosteno.LOG_INFO,
		Codec:     gosteno.NewJsonPrettifier(gosteno.EXCLUDE_DATA),
		EnableLOC: true,
	}
	gosteno.Init(c)
}

// To get the base foundation configuration see the Pivotal CF Data Collector @
// https://docs.google.com/a/pivotal.io/spreadsheets/d/1XHKSrJiQIu5MWGpPYWbMY8M09eqe-GV8MQsl_mqw1RM/edit#gid=0
func main() {

	logger := gosteno.NewLogger("cfops")
	installer = install.New(logger)
	backuper = backup.New(logger)
	app := app.New(logger)

	app.Run(os.Args)
}

// Installer Commands
func StartDeployment() func(c *cli.Context) {
	return func(c *cli.Context) {
		installer.StartDeployment(c.Args())
	}
}

// Backup Commands
func BackupDeployment() func(c *cli.Context) {
	return func(c *cli.Context) {
		backuper.ValidateSoftware(c.Args())
	}
}
