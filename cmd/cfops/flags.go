package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

const (
	opsManagerHostDescr string = "hostname for Ops Manager"
	directorUserDescr   string = "username for Ops Mgr Director VM"
	directorPassDescr   string = "password for Ops Mgr Director VM"
	opsManagerUserDescr string = "username for Ops Manager"
	opsManagerPassDescr string = "password for Ops Manager"
	destdescr           string = "directory of the Cloud Foundry backup archive"
	tilelistdescr       string = "a csv list of the tiles you would like to run the operation on"
)

var (
	opsManagerHostFlag = []string{"opsmanagerhost", "omh"}
	directorUserFlag   = []string{"directoruser", "du"}
	directorPassFlag   = []string{"directorpass", "dp"}
	opsManagerUserFlag = []string{"opsmanageruser", "omu"}
	opsManagerPassFlag = []string{"opsmanagerpass", "omp"}
	destFlag           = []string{"destination", "d"}
	tilelistFlag       = []string{"tilelist", "tl"}
)

type flagSet struct {
	host           string
	directorUser   string
	directorPass   string
	opsManagerUser string
	opsManagerPass string
	dest           string
	tilelist       string
}

func (s *flagSet) Host() string {
	return s.host
}

func (s *flagSet) DirectorUser() string {
	return s.directorUser
}

func (s *flagSet) DirectorPass() string {
	return s.directorPass
}

func (s *flagSet) OpsManagerUser() string {
	return s.opsManagerUser
}

func (s *flagSet) OpsManagerPass() string {
	return s.opsManagerPass
}

func (s *flagSet) Dest() string {
	return s.dest
}

func (s *flagSet) Tilelist() string {
	return s.tilelist
}

func hasValidBackupRestoreFlags(fs *flagSet) bool {
	res := (fs.Host() != "" && fs.DirectorUser() != "" && fs.DirectorPass() != "" && fs.OpsManagerUser() != "" && fs.OpsManagerPass() != "" && fs.Dest() != "")

	if res == false {
		fmt.Println("OpsManagerHost: ", fs.Host())
		fmt.Println("DirectorUser: ", fs.DirectorUser())
		fmt.Println("DirectorPass: ", fs.DirectorPass())
		fmt.Println("OpsManagerUser: ", fs.OpsManagerUser())
		fmt.Println("OpsManagerPass: ", fs.OpsManagerPass())
		fmt.Println("Destination: ", fs.Dest())
	}
	return res
}

var backupRestoreFlags = []cli.Flag{
	cli.StringFlag{
		Name:   strings.Join(opsManagerHostFlag, ", "),
		Value:  "",
		Usage:  opsManagerHostDescr,
		EnvVar: "",
	},
	cli.StringFlag{
		Name:   strings.Join(directorUserFlag, ", "),
		Value:  "",
		Usage:  directorUserDescr,
		EnvVar: "",
	},
	cli.StringFlag{
		Name:   strings.Join(directorPassFlag, ", "),
		Value:  "",
		Usage:  directorPassDescr,
		EnvVar: "",
	},
	cli.StringFlag{
		Name:   strings.Join(opsManagerUserFlag, ", "),
		Value:  "",
		Usage:  opsManagerUserDescr,
		EnvVar: "",
	},
	cli.StringFlag{
		Name:   strings.Join(opsManagerPassFlag, ", "),
		Value:  "",
		Usage:  opsManagerPassDescr,
		EnvVar: "",
	},
	cli.StringFlag{
		Name:   strings.Join(destFlag, ", "),
		Value:  "",
		Usage:  destdescr,
		EnvVar: "",
	},
	cli.StringFlag{
		Name:   strings.Join(tilelistFlag, ", "),
		Value:  "",
		Usage:  tilelistdescr,
		EnvVar: "",
	},
}
