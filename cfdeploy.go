package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {


	// To get the base foundation configuration see the Pivotal CF Data Collector @
	// https://docs.google.com/a/pivotal.io/spreadsheets/d/1XHKSrJiQIu5MWGpPYWbMY8M09eqe-GV8MQsl_mqw1RM/edit#gid=0

	app := cli.NewApp()
	app.Name = "cfdeploy"
	app.Usage = "Cloud Foundry Deployment tool for IaaS install and deployment automation"

	// The `cfdeploy` command default without argument
	app.Action = func(c *cli.Context) {
		arg := ""
		if len(c.Args()) > 0 {
			arg = c.Args()[0]
			println("Try using 'cfdeploy help'.  Invalid argument: ", arg)
		} else if len(c.Args()) == 0 {
			println("To get started, try using 'cfdeploy help'")
		}
	}

	// Global application flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "iaas, i",
			Value:  "aws, vsphere, vcloud, openstack",
			Usage:  "set the IaaS type to target for deployment",
			EnvVar: "CF_IAAS",
		},
		cli.StringFlag{
			Name:   "debug, d",
			Value:  "true, false",
			Usage:  "enable/disable debug output",
			EnvVar: "CF_TRACE",
		},
		cli.StringFlag{
			Name:   "lang, l",
			Value:  "en, es",
			Usage:  "language for the cfdeploy cli responses",
			EnvVar: "CF_LANG",
		},
	}

	// CLI arg functionality
	app.Commands = []cli.Command{
		{
			Name:        "survey",
			ShortName:   "s",
			Usage:       "analyze and inspect the deployment environment",
			Description: "Survey the target IaaS environment to determine suitability for deploying a cloud foundry foundation",
			Subcommands: []cli.Command{
				{
					Name:        "verify",
					Usage:       "verify that the target IaaS is Cloud Foundry ready",
					Description: "analyze the readiness of the Cloud Foundry IaaS target",
					Action: func(c *cli.Context) {
						println("Verified that the target IaaS is Cloud Foundry: ", c.Args().First())
					},
				},
				{
					Name:        "report",
					Usage:       "produce a report on IaaS related to Cloud Foundry",
					Description: "produce a report against the target IaaS environment in the context of Cloud Foundry deployment",
					Action: func(c *cli.Context) {
						println("report produced: ", c.Args().First())
					},
				},
				{
					Name:        "stats",
					Usage:       "produce statistics on IaaS",
					Description: "produce useful statistics against the target IaaS environment",
					Action: func(c *cli.Context) {
						println("new jumpbox deployed to IaaS: ", c.Args().First())
					},
				},

			},
		},
		{
			Name:        "prepare",
			ShortName:   "p",
			Usage:       "prepare the deployment environment",
			Description: "Build and configure an environment that will be used to run the cloud foundry deployment from",
			Subcommands: []cli.Command{
				{
					Name:        "jumpbox",
					Usage:       "add a new jumpbox on the IaaS",
					Description: "add a jumpbox to the target IaaS environment which can be used to deploy Cloud Foundry",
					Action: func(c *cli.Context) {
						println("new jumpbox deployed to IaaS: ", c.Args().First())
					},
				},
			},
		},
		{
			Name:        "install",
			ShortName:   "in",
			Usage:       "install cloud foundry to an iaas",
			Description: "Begin the installation of Cloud Foundry to a selected iaas",
			Subcommands: []cli.Command{
				{
					Name:        "add",
					Usage:       "add a new deployment",
					Description: "use the provided deployment template to deploy a new cloud foundry foundation into the iaas",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "release",
							Value: "latest, <version number>",
							Usage: "Set the release to use",
						},
					},
					Action: func(c *cli.Context) {
						println("new deployment with template: ", c.Args().First())
					},
				},
				{
					Name:        "destroy",
					Usage:       "destroy an existing deployment",
					Description: "destroy an existing cloud foundry foundation deployment from the iaas",
					Action: func(c *cli.Context) {
						println("destroyed deployment: ", c.Args().First())
					},
				},
				{
					Name:        "move",
					Usage:       "move an existing deployment to another iaas location",
					Description: "destroy an existing cloud foundry foundation deployment from the iaas",
					Action: func(c *cli.Context) {
						println("move deployment: ", c.Args().First())
					},
				},
				{
					Name:        "dump",
					Usage:       "dump the configuration information of an existing deployment",
					Description: "dump an existing cloud foundry foundation deployment configuration from the iaas",
					Action: func(c *cli.Context) {
						println("dump deployment: ", c.Args().First())
					},
				},
				{
					Name:        "backup",
					Usage:       "backup an existing deployment",
					Description: "backup an existing cloud foundry foundation deployment from the iaas",
					Action: func(c *cli.Context) {
						println("backup deployment: ", c.Args().First())
					},
				},
				{
					Name:        "restore",
					Usage:       "restore an deployment from a backup",
					Description: "restore an existing cloud foundry foundation deployment from a backup",
					Action: func(c *cli.Context) {
						println("restore deployment: ", c.Args().First())
					},
				},

			},
		},
		{
			Name:        "start",
			ShortName:   "s",
			Usage:       "start up an entire cloud foundry foundation",
			Description: "start all the VMs in an existing cloud foundry deployment",
			Action: func(c *cli.Context) {
				println("starting: ", c.Args().First())
			},
		},
		{
			Name:        "restart",
			ShortName:   "r",
			Usage:       "shutdown and restart an entire cloud foundry foundation",
			Description: "shutdown and restart all the VMs in an existing cloud foundry deployment",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "force",
					Value: "true, false",
					Usage: "Force the restart without prompting for confirmation",
				},
			},
			Action: func(c *cli.Context) {
				println("restarting: ", c.Args().First())
			},
		},
		{
			Name:        "shutdown",
			ShortName:   "stop",
			Usage:       "shutdown and stop an entire cloud foundry foundation",
			Description: "shutdown all the VMs in an existing cloud foundry deployment",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "force",
					Value: "true, false",
					Usage: "Force the shutdown without prompting for confirmation",
				},
			},
			Action: func(c *cli.Context) {
				println("shutting down: ", c.Args().First())
			},
		},
		{
			Name:        "test",
			ShortName:   "t",
			Usage:       "test the Cloud Foundry deployment and underlying IaaS environment",
			Description: "Run various tests against the target IaaS environment",
			Subcommands: []cli.Command{
				{
					Name:        "chaos",
					Usage:       "run chaos tests against the IaaS and Cloud Foundry foundation",
					Description: "Chaos testing for the IaaS and Cloud Foundry foundation",
					Action: func(c *cli.Context) {
						println("Chaos testing: ", c.Args().First())
					},
				},
				{
					Name:        "vulnerability",
					Usage:       "run vulnerability testing against the IaaS and Cloud Foundry foundation",
					Description: "Chaos testing for the IaaS and Cloud Foundry foundation",
					Action: func(c *cli.Context) {
						println("vulnerability testing: ", c.Args().First())
					},
				},
				{
					Name:        "report",
					Usage:       "produce a test report on IaaS related to Cloud Foundry",
					Description: "produce a report against the target IaaS environment in the context of Cloud Foundry deployment",
					Action: func(c *cli.Context) {
						println("test report produced: ", c.Args().First())
					},
				},
				{
					Name:        "stats",
					Usage:       "produce statistics on IaaS",
					Description: "produce useful statistics against the target IaaS environment",
					Action: func(c *cli.Context) {
						println("test stats: ", c.Args().First())
					},
				},

			},
		},
	}

	cli.CommandHelpTemplate = `NAME:
   {{.Name}} - {{.Description}}
{{with .ShortName}}

ALIAS:
   {{.}}
{{end}}

USAGE:
   cfdeploy {{.Name}}{{if .Flags}} [command options]{{end}} [arguments...]{{if .Description}}

OPTIONS:
{{range .Flags}}{{.}}
{{end}}{{ end }}`

	app.Run(os.Args)

}
