package install

import (
	"github.com/cloudfoundry/gosteno"
	"github.com/pivotalservices/cfops/system"
)

type Installer struct {
	CommandRunner system.CommandRunner
}

func New(logger *gosteno.Logger) *Installer {
	commandRunner := new(system.OSCommandRunner)
	commandRunner.Logger = logger
	return &Installer{
		CommandRunner: commandRunner,
	}
}

func (installer *Installer) StartDeployment(args []string) error {
	err := installer.CommandRunner.Run("echo", "WHOOOA", "slow down!")
	if err != nil {
		return err
	}
	return nil
}
