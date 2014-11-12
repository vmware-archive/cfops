package install

import (
	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/pivotalservices/cfops/system"
)

type Installer struct {
}

func New() *Installer {
	return &Installer{}
}

func (installer *Installer) StartDeployment() error {
	commandRunner := new(system.OSCommandRunner)
	logger := cf_lager.New("ops-broker")
	commandRunner.Logger = logger
	err := commandRunner.Run("echo", "WHOOOA", "slow down!")
	if err != nil {
		return err
	}
	return nil
}
