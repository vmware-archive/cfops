package backup

import (
	"github.com/pivotalservices/cfops/system"
)

type Backuper struct {
}

func New(f system.CommandFactory, commandRunner system.CommandRunner) *Backuper {

	f.Register("backup", BackupCommand{
		CommandRunner: commandRunner,
	})

	return &Backuper{}
}
