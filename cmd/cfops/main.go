package main

import (
	"fmt"
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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "logLevel",
			Value: "info",
			Usage: "log level: debug, info, error or fatal",
		},
	}
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
					Name:   "hostname, host",
					Value:  "",
					Usage:  "hostname for Ops Manager",
					EnvVar: "",
				},
				cli.StringFlag{
					Name:   "username, u",
					Value:  "",
					Usage:  "username for Ops Manager",
					EnvVar: "",
				},
				cli.StringFlag{
					Name:   "password, p",
					Value:  "",
					Usage:  "password for Ops Manager",
					EnvVar: "",
				},
				cli.StringFlag{
					Name:   "tempestpassword, tp",
					Value:  "",
					Usage:  "password for the Ops Manager tempest user",
					EnvVar: "",
				},
				cli.StringFlag{
					Name:   "destination, d",
					Value:  "",
					Usage:  "directory where the Cloud Foundry backup should be stored",
					EnvVar: "",
				},
			},
			Action: func(c *cli.Context) {
				var err error

				if c.String("hostname") == "" || c.String("username") == "" || c.String("password") == "" || c.String("tempestpassword") == "" || c.String("destination") == "" {
					cli.ShowCommandHelp(c, "backup")

				} else {
					err = cfbackup.RunBackupPipeline(c.String("hostname"), c.String("username"), c.String("password"), c.String("tempestpassword"), c.String("destination"))
				}

				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Backup complete successfully.")
				}
			},
		},
	}...)

	return app
}
