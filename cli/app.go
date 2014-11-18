package cli

import (
	"github.com/codegangsta/cli"
)

type Flag struct {
	Name       string
	Value      interface{}
	Type       string
	StringFlag cli.StringFlag
	BoolFlag   cli.BoolFlag
}

func NewFlag(name string, value interface{}, usage string, envVar string) (f *Flag) {
	f = &Flag{
		Name: name,
	}
	switch value.(type) {
	case bool:
		f.Type = "bool"
		f.BoolFlag = cli.BoolFlag{
			Name:   f.Name,
			Usage:  usage,
			EnvVar: envVar,
		}
		return
	case string:
		f.Type = "string"
		f.StringFlag = cli.StringFlag{
			Name:   f.Name,
			Value:  value.(string),
			Usage:  usage,
			EnvVar: envVar,
		}
		return
	}
	return
}

func (flag *Flag) ParseString() string {
	// fmt.Printf("string flag %s has value %s\n", flag.Name, flag.Value)
	return flag.Value.(string)
}

func (flag *Flag) ParseBool() bool {
	// fmt.Printf("bool flag %s has value %v\n", flag.Name, flag.Value)
	return flag.Value.(bool)
}

func NewApp(cmdFactory CommandFactory, globalFlags []*Flag) *cli.App {

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
	for _, flag := range globalFlags {
		if flag.Type == "bool" {
			app.Flags = append(app.Flags, flag.BoolFlag)
		} else {
			app.Flags = append(app.Flags, flag.StringFlag)
		}
	}

	// Create all the CLI commands
	for _, cmd := range cmdFactory.Commands() {
		cliCommand := createCLICommand(cmd)
		subcommands := cmdFactory.Subcommands(cmd)
		if subcommands != nil {
			for _, subCmd := range subcommands {
				subCliCommand := createCLICommand(subCmd)
				setAction(subCmd, &subCliCommand)
				cliCommand.Subcommands = append(cliCommand.Subcommands, subCliCommand)
			}
		} else {
			setAction(cmd, &cliCommand)
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

	cli.CommandHelpTemplate = getCommandLineTemplate()

	return app
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

func createCLICommand(cmd Command) (cliCommand cli.Command) {
	return cli.Command{
		Name:        cmd.Metadata().Name,
		ShortName:   cmd.Metadata().ShortName,
		Usage:       cmd.Metadata().Usage,
		Description: cmd.Metadata().Description,
		Flags:       cmd.Metadata().Flags,
	}
}

func setAction(cmd Command, cliCommand *cli.Command) {
	cliCommand.Action = func(c *cli.Context) {
		err := cmd.Run(c.Args())
		if err != nil {
			panic(err)
		}
	}
}
