package backup

import (
	"github.com/pivotalservices/cfops/cli"
	"github.com/pivotalservices/cfops/system"
)

type RestoreCommand struct {
	CommandRunner system.CommandRunner
}

func (cmd RestoreCommand) Metadata() cli.CommandMetadata {
	return cli.CommandMetadata{
		Name:        "restore",
		ShortName:   "r",
		Usage:       "restore an deployment from a backup",
		Description: "restore an existing cloud foundry foundation deployment from a backup",
	}
}

func (cmd RestoreCommand) Run(args []string) (err error) {
	println("restoring backup: " + args[0])
	return
}
