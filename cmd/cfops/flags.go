package main

import (
	"strings"

	"github.com/codegangsta/cli"
)

const (
	hostdescr     string = "hostname for Ops Manager"
	userdescr     string = "username for Ops Manager"
	passdescr     string = "password for Ops Manager"
	tpassdescr    string = "password for the Ops Manager tempest user"
	destdescr     string = "directory of the Cloud Foundry backup archive"
	tilelistdescr string = "a csv list of the tiles you would like to run the operation on"
)

var (
	hostflag     = []string{"hostname", "host"}
	userflag     = []string{"username", "u"}
	passflag     = []string{"password", "p"}
	tpassflag    = []string{"tempestpassword", "tp"}
	destflag     = []string{"destination", "d"}
	tilelistflag = []string{"tilelist", "tl"}
)

type flagSet struct {
	host     string
	user     string
	pass     string
	tpass    string
	dest     string
	tilelist string
}

func (s *flagSet) Host() string {
	return s.host
}

func (s *flagSet) User() string {
	return s.user
}

func (s *flagSet) Pass() string {
	return s.pass
}

func (s *flagSet) Tpass() string {
	return s.tpass
}

func (s *flagSet) Dest() string {
	return s.dest
}

func (s *flagSet) Tilelist() string {
	return s.tilelist
}

func hasValidBackupRestoreFlags(fs *flagSet) bool {
	return (fs.Host() != "" && fs.User() != "" && fs.Pass() != "" && fs.Tpass() != "" && fs.Dest() != "")
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
	cli.StringFlag{
		Name:   strings.Join(tilelistflag, ", "),
		Value:  "",
		Usage:  tilelistdescr,
		EnvVar: "",
	},
}
