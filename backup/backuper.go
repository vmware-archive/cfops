package backup

import (
	"github.com/cloudfoundry/gosteno"
	"github.com/pivotalservices/cfops/system"
)

type Backuper struct {
	CommandRunner system.CommandRunner
}

func New(logger *gosteno.Logger) *Backuper {
	commandRunner := new(system.OSCommandRunner)
	commandRunner.Logger = logger
	return &Backuper{
		CommandRunner: commandRunner,
	}
}

func (installer *Backuper) ValidateSoftware(args []string) error {
	params := make([]string, len(args)+1)
	params = append(params, "validate_software")
	params = append(params, args...)
	err := installer.CommandRunner.Run("./backup/scripts/backup.sh", params...)
	if err != nil {
		return err
	}
	return nil
}
