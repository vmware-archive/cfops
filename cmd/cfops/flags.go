package main

import (
	"strings"

	"github.com/codegangsta/cli"
)

const (
	hostdescr  string = "hostname for Ops Manager"
	userdescr  string = "username for Ops Manager"
	passdescr  string = "password for Ops Manager"
	tpassdescr string = "password for the Ops Manager tempest user"
	destdescr  string = "directory of the Cloud Foundry backup archive"
)

var (
	hostflag  = []string{"hostname", "host"}
	userflag  = []string{"username", "u"}
	passflag  = []string{"password", "p"}
	tpassflag = []string{"tempestpassword", "tp"}
	destflag  = []string{"destination", "d"}
)

func hasValidBackupRestoreFlags(c *cli.Context) bool {
	return (c.String(hostflag[0]) != "" && c.String(userflag[0]) != "" && c.String(passflag[0]) != "" && c.String(tpassflag[0]) != "" && c.String(destflag[0]) != "")
}

var backupRestoreFlags = []cli.Flag{
	cli.StringFlag{
		Name:   strings.Join(hostflag, ", "),
		Value:  "",
		Usage:  hostdescr,
		EnvVar: "",
	},
	cli.StringFlag{
		Name:   strings.Join(userflag, ", "),
		Value:  "",
		Usage:  userdescr,
		EnvVar: "",
	},
	cli.StringFlag{
		Name:   strings.Join(passflag, ", "),
		Value:  "",
		Usage:  passdescr,
		EnvVar: "",
	},
	cli.StringFlag{
		Name:   strings.Join(tpassflag, ", "),
		Value:  "",
		Usage:  tpassdescr,
		EnvVar: "",
	},
	cli.StringFlag{
		Name:   strings.Join(destflag, ", "),
		Value:  "",
		Usage:  destdescr,
		EnvVar: "",
	},
}
