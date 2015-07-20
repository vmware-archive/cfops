package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

const (
	opsManagerHostDescr string = "hostname for Ops Manager"
	adminUserDescr      string = "username for Ops Mgr admin VM"
	adminPassDescr      string = "password for Ops Mgr admin VM"
	opsManagerUserDescr string = "username for Ops Manager"
	opsManagerPassDescr string = "password for Ops Manager"
	destdescr           string = "adminy of the Cloud Foundry backup archive"
	tilelistdescr       string = "a csv list of the tiles you would like to run the operation on"
)

var (
	opsManagerHostFlag = []string{"opsmanagerhost", "omh"}
	adminUserFlag      = []string{"adminuser", "du"}
	adminPassFlag      = []string{"adminpass", "dp"}
	opsManagerUserFlag = []string{"opsmanageruser", "omu"}
	opsManagerPassFlag = []string{"opsmanagerpass", "omp"}
	destFlag           = []string{"destination", "d"}
	tilelistFlag       = []string{"tilelist", "tl"}
)

type flagSet struct {
	host           string
	adminUser      string
	adminPass      string
	opsManagerUser string
	opsManagerPass string
	dest           string
	tilelist       string
}

func (s *flagSet) Host() string {
	return s.host
}

func (s *flagSet) AdminUser() string {
	return s.adminUser
}

func (s *flagSet) AdminPass() string {
	return s.adminPass
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
	res := (fs.Host() != "" && fs.AdminUser() != "" && fs.AdminPass() != "" && fs.OpsManagerUser() != "" && fs.OpsManagerPass() != "" && fs.Dest() != "")

	if res == false {
		fmt.Println("OpsManagerHost: ", fs.Host())
		fmt.Println("adminUser: ", fs.AdminUser())
		fmt.Println("adminPass: ", fs.AdminPass())
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
		Name:   strings.Join(adminUserFlag, ", "),
		Value:  "",
		Usage:  adminUserDescr,
		EnvVar: "",
	},
	cli.StringFlag{
		Name:   strings.Join(adminPassFlag, ", "),
		Value:  "",
		Usage:  adminPassDescr,
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
