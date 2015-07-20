package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfops"
)

const (
	backup_full_name  string = "backup"
	backup_short_name        = "b"
	backup_usage             = "backup --opsmanagerhost <host> --adminuser <usr> --adminpass <pass> --opsmanageruser <opsuser> --opsmanagerpass <opspass> -d <dir> --tl 'opsmanager, er'"
	backup_descr             = "backup a Cloud Foundry deployment, including Ops Manager configuration, databases, and blob store"
)

var backupCli = cli.Command{
	Name:        backup_full_name,
	ShortName:   backup_short_name,
	Usage:       backup_usage,
	Description: backup_descr,
	Flags:       backupRestoreFlags,
	Action: func(c *cli.Context) {
		var (
			err error
			fs  = &flagSet{
				host:           c.String(opsManagerHostFlag[0]),
				adminUser:      c.String(adminUserFlag[0]),
				adminPass:      c.String(adminPassFlag[0]),
				opsManagerUser: c.String(opsManagerUserFlag[0]),
				opsManagerPass: c.String(opsManagerPassFlag[0]),
				dest:           c.String(destFlag[0]),
				tilelist:       c.String(tilelistFlag[0]),
			}
		)

		if hasValidBackupRestoreFlags(fs) {
			cfops.SetupSupportedTiles(fs)
			err = cfops.RunPipeline(fs, cfops.Backup)

			if err != nil {
				fmt.Println(err)

			} else {
				fmt.Println(backup_full_name, " completed successfully.")
			}

		} else {
			cli.ShowCommandHelp(c, backup_full_name)
		}
	},
}
