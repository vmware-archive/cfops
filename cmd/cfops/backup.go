package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var backupFlags []cli.Flag = []cli.Flag{
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
}

type commandFunc func(c *cli.Context)

type runPipeline func(string, string, string, string, string) error

func runBackupRestore(command string, run runPipeline) commandFunc {
	return func(c *cli.Context) {
		runBackupRestoreCmd(command, c, run)
	}
}

func runBackupRestoreCmd(command string, c *cli.Context, run runPipeline) {
	var err error

	if c.String("hostname") == "" || c.String("username") == "" || c.String("password") == "" || c.String("tempestpassword") == "" || c.String("destination") == "" {
		cli.ShowCommandHelp(c, command)

	} else {
		err = run(c.String("hostname"), c.String("username"), c.String("password"), c.String("tempestpassword"), c.String("destination"))
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(fmt.Sprintf("%s compeleted successfully", command))
	}
}
