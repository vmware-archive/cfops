package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfbackup"
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
		{
			Name:        "backup",
			ShortName:   "b",
			Usage:       "backup a Cloud Foundry deployment",
			Description: "Backup a Cloud Foundry deployment, including Ops Manager configuration, databases, and blob store",
			Flags:       backupFlags,

			Action: runBackupRestore("backup", func(hostname, username, password, tempestpassword, destination string) error {
				return cfbackup.RunBackupPipeline(hostname, username, password, tempestpassword, destination)
			}),
		},
		{
			Name:        "restore",
			ShortName:   "r",
			Usage:       "restore a Cloud Foundry deployment",
			Description: "Restore a Cloud Foundry deployment, including Ops Manager configuration, databases, and blob store",
			Flags:       backupFlags,

			Action: runBackupRestore("restore", func(hostname, username, password, tempestpassword, destination string) error {
				return cfbackup.RunRestorePipeline(hostname, username, password, tempestpassword, destination)
			}),
		},
	}...)

	return app
}
