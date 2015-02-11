package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := NewApp()
	app.Run(os.Args)
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
