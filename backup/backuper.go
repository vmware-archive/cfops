package backup

import (
	"github.com/pivotalservices/cfops/cli"
	"github.com/pivotalservices/cfops/system"
)

type Backuper struct {
}

func New(factory cli.CommandFactory, runner system.CommandRunner) Backuper {

	factory.Register("backup", BackupCommand{
		CommandRunner: runner,
		Logger:        factory.GetLogger(),
	}).Register("restore", RestoreCommand{
		CommandRunner: runner,
	})

	return Backuper{}
}
