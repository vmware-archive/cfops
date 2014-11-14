package backup

import (
	"github.com/pivotalservices/cfops/system"
)

type RestoreCommand struct {
	CommandRunner system.CommandRunner
}

func (cmd RestoreCommand) Metadata() system.CommandMetadata {
	return system.CommandMetadata{
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

func (cmd RestoreCommand) HasSubcommands() bool {
	return false
}

func (cmd RestoreCommand) Subcommands() (commands []system.Command) {
	return
}
