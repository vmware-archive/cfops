package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

const (
	errExitCode           = 1
	helpExitCode          = 2
	cleanExitCode         = 0
	opsManagerHost string = "opsmanagerHost"
	adminUser      string = "adminUser"
	adminPass      string = "adminPass"
	opsManagerUser string = "opsManagerUser"
	opsManagerPass string = "opsManagerPass"
	dest           string = "destination"
	tilelist       string = "tilelist"
)

var (
	flagList = map[string]flagBucket{
		opsManagerHost: flagBucket{
			Flag:   []string{"opsmanagerhost", "omh"},
			Desc:   "hostname for Ops Manager",
			EnvVar: "CFOPS_HOST",
		},
		adminUser: flagBucket{
			Flag:   []string{"adminuser", "du"},
			Desc:   "username for Ops Mgr admin (Ops Manager WebConsole Credentials)",
			EnvVar: "CFOPS_ADMIN_USER",
		},
		adminPass: flagBucket{
			Flag:   []string{"adminpass", "dp"},
			Desc:   "password for Ops Mgr admin (Ops Manager WebConsole Credentials)",
			EnvVar: "CFOPS_ADMIN_PASS",
		},
		opsManagerUser: flagBucket{
			Flag:   []string{"opsmanageruser", "omu"},
			Desc:   "username for Ops Manager VM Access (used for ssh connections)",
			EnvVar: "CFOPS_OM_USER",
		},
		opsManagerPass: flagBucket{
			Flag:   []string{"opsmanagerpass", "omp"},
			Desc:   "password for Ops Manager VM Access (used for ssh connections)",
			EnvVar: "CFOPS_OM_PASS",
		},
		dest: flagBucket{
			Flag:   []string{"destination", "d"},
			Desc:   "path of the Cloud Foundry backup archive",
			EnvVar: "CFOPS_BACKUP_PATH",
		},
		tilelist: flagBucket{
			Flag:   []string{"tilelist", "tl"},
			Desc:   "a csv list of the tiles you would like to run the operation on",
			EnvVar: "CFOPS_TILE_LIST",
		},
	}
)

type (
	flagSet struct {
		host           string
		adminUser      string
		adminPass      string
		opsManagerUser string
		opsManagerPass string
		dest           string
		tilelist       string
	}

	flagBucket struct {
		Flag   []string
		Desc   string
		EnvVar string
	}
)

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

var backupRestoreFlags = func() (flags []cli.Flag) {
	for _, v := range flagList {
		flags = append(flags, cli.StringFlag{
			Name:   strings.Join(v.Flag, ", "),
			Value:  "",
			Usage:  v.Desc,
			EnvVar: v.EnvVar,
		})
	}
	return
}()
