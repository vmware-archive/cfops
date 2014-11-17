package start

import (
	"github.com/pivotalservices/cfops/cli"
	"github.com/pivotalservices/cfops/system"
)

type Starter struct {
}

func New(factory cli.CommandFactory, runner system.CommandRunner) Starter {

	factory.Register("start", StartCommand{})

	return Starter{}
}
