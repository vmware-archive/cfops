package backup

import (
	"github.com/pivotalservices/cfops/system"
)

type BackupCommand struct {
	CommandRunner system.CommandRunner
}

func (cmd BackupCommand) Metadata() system.CommandMetadata {
	return system.CommandMetadata{
		Name:        "backup",
		ShortName:   "b",
		Usage:       "backup an existing deployment",
		Description: "backup an existing cloud foundry foundation deployment from the iaas",
	}
}

func (cmd BackupCommand) Subcommands() (commands []system.Command) {
	return
}

func (cmd BackupCommand) Run(args []string) error {
	params := make([]string, len(args)+1)
	params = append(params, "validate_software")
	params = append(params, args...)
	err := cmd.CommandRunner.Run("./backup/scripts/backup.sh", params...)
	if err != nil {
		return err
	}
	return nil
}
