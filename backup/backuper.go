package backup

import (
	"github.com/pivotalservices/cfops/system"
)

type Backuper struct {
	CommandRunner system.CommandRunner
}

func New(f system.CommandFactory, commandRunner system.CommandRunner) *Backuper {
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

// type Backuper struct {
// }

// func New(f system.CommandFactory, commandRunner *system.CommandRunner) *Backuper {

// 	f.Register("start", StartCommand{
// 		CommandRunner: commandRunner,
// 	})

// 	return &Backuper{}
// }
