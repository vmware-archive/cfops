package main

import (
	"github.com/cloudfoundry/gosteno"
	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfops/app"
	"github.com/pivotalservices/cfops/backup"
	"github.com/pivotalservices/cfops/install"
	"github.com/pivotalservices/cfops/system"
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

	commandFactory := system.NewCommandFactory(logger)

	commandRunner := system.OSCommandRunner{}
	commandRunner.Logger = logger

	installer = install.New(commandFactory, commandRunner)
	backuper = backup.New(commandFactory, commandRunner)

	app := app.New(commandFactory)

	cli.CommandHelpTemplate = getCommandLineTemplate()

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

func getCommandLineTemplate() string {
	return `NAME:
   {{.Name}} - {{.Description}}
{{with .ShortName}}

ALIAS:
   {{.}}
{{end}}

USAGE:
   cfops {{.Name}}{{if .Flags}} [command options]{{end}} [arguments...]{{if .Description}}

OPTIONS:
{{range .Flags}}{{.}}
{{end}}{{ end }}`
}
