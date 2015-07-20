package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfops"
)

const (
	restore_full_name  string = "restore"
	restore_short_name        = "r"
	restore_usage             = "restore --opsmanagerhost <host> --adminuser <usr> --adminpass <pass> --opsmanageruser <opsuser> --opsmanagerpass <opspass> -d <dir> --tl 'opsmanager, er'"
	restore_descr             = "Restore a Cloud Foundry deployment, including Ops Manager configuration, databases, and blob store"
)

var restoreCli = cli.Command{
	Name:        restore_full_name,
	ShortName:   restore_short_name,
	Usage:       restore_usage,
	Description: restore_descr,
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
			err = cfops.RunPipeline(fs, cfops.Restore)

			if err != nil {
				fmt.Println(err)

			} else {
				fmt.Println(restore_full_name, " completed successfully.")
			}

		} else {
			cli.ShowCommandHelp(c, restore_full_name)
		}
	},
}
