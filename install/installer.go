package install

import (
	"github.com/pivotalservices/cfops/system"
)

type Installer struct {
}

func New(f system.CommandFactory, commandRunner system.CommandRunner) *Installer {

	f.Register("start", StartCommand{
		CommandRunner: commandRunner,
	})

	return &Installer{}
}
