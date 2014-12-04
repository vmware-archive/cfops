package main

import (
	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfops/backup"
)

func main() {
	// Call NewApp!
}

func NewApp() *cli.App {

	app := cli.NewApp()
	app.Name = "cfops"
	app.Usage = "Cloud Foundry Operations tool for IaaS installation, deployment, and management automation"
	app.Commands = append(app.Commands, []cli.Command{
		{
			Name:        "backup",
			ShortName:   "b",
			Usage:       "backup a Cloud Foundry deployment",
			Description: "Backup a Cloud Foundry deployment, including Ops Manager configuration, databases, and blob store",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "hostname, h",
					Value:  "",
					Usage:  "hostname for Ops Manager",
					EnvVar: "",
				},
			},
			// Hostname      string
			// Username      string
			// Password      string
			// TPassword     string
			// Target        string
			Action: func(c *cli.Context) {
				context := &backup.BackupContext{}
				context.Run()
			},
		},
	}...)

	return app
}
