package app

import (
	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfops/system"
)

func New(cmdFactory system.CommandFactory) *cli.App {

	app := cli.NewApp()
	app.Name = "cfops"
	app.Usage = "Cloud Foundry Operations tool for IaaS installation, deployment, and management automation"

	// The `cfops` command default without argument
	app.Action = func(c *cli.Context) {
		arg := ""
		if len(c.Args()) > 0 {
			arg = c.Args()[0]
			println("Try using 'cfops help'.  Invalid argument: ", arg)
		} else if len(c.Args()) == 0 {
			println("To get started, try using 'cfops help'")
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
			Usage:  "language for the cfops cli responses",
			EnvVar: "CF_LANG",
		},
	}

	// Create all the CLI commands
	for _, cmd := range cmdFactory.Commands() {
		cliCommand := createCLICommand(cmd)
		if len(cmd.Subcommands()) > 0 {
			for _, subCmd := range cmd.Subcommands() {
				cliCommand.Subcommands = append(cliCommand.Subcommands, createCLICommand(subCmd))
			}
		}
		app.Commands = append(app.Commands, cliCommand)
	}

	// TODO: Move metadata into each command
	app.Commands = append(app.Commands, []cli.Command{
		{
			Name:        "survey",
			ShortName:   "sur",
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
	}...)

	return app
}

func createCLICommand(cmd system.Command) (cliCommand cli.Command) {
	return cli.Command{
		Name:        cmd.Metadata().Name,
		ShortName:   cmd.Metadata().ShortName,
		Usage:       cmd.Metadata().Usage,
		Description: cmd.Metadata().Description,
		Flags:       cmd.Metadata().Flags,
		Action: func(c *cli.Context) {
			err := cmd.Run(c.Args())
			if err != nil {
				panic(err)
			}
		},
	}
}
