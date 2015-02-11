package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfbackup"
)

const (
	restore_full_name  string = "restore"
	restore_short_name        = "r"
	restore_usage             = "restore -host <host> -u <usr> -p <pass> -tp <tpass> -d <dir>"
	restore_descr             = "Restore a Cloud Foundry deployment, including Ops Manager configuration, databases, and blob store"
)

var restoreCli = cli.Command{
	Name:        restore_full_name,
	ShortName:   restore_short_name,
	Usage:       restore_usage,
	Description: restore_descr,
	Flags:       backupRestoreFlags,
	Action: func(c *cli.Context) {
		var err error

		if hasValidBackupRestoreFlags(c) {
			err = cfbackup.RunRestorePipeline(c.String(hostflag[0]), c.String(userflag[0]), c.String(passflag[0]), c.String(tpassflag[0]), c.String(destflag[0]))

		} else {
			cli.ShowCommandHelp(c, backup_full_name)
		}

		if err != nil {
			fmt.Println(err)

		} else {
			fmt.Println(restore_full_name, " completed successfully.")
		}
	},
}
