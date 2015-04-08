package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfops"
)

const (
	backup_full_name  string = "backup"
	backup_short_name        = "b"
	backup_usage             = "backup -host <host> -u <usr> -p <pass> --tp <tpass> -d <dir> --tl 'opsmanager, er'"
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
				host:     c.String(hostflag[0]),
				user:     c.String(userflag[0]),
				pass:     c.String(passflag[0]),
				tpass:    c.String(tpassflag[0]),
				dest:     c.String(destflag[0]),
				tilelist: c.String(tilelistflag[0]),
			}
		)

		if hasValidBackupRestoreFlags(fs) {
			cfops.SetupSupportedTiles(fs)
			err = cfops.RunPipeline(fs, cfops.Backup)

		} else {
			cli.ShowCommandHelp(c, backup_full_name)
		}

		if err != nil {
			fmt.Println(err)

		} else {
			fmt.Println(backup_full_name, " completed successfully.")
		}
	},
}
