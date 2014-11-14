package start

import (
	"github.com/pivotalservices/cfops/system"
)

type Starter struct {
}

func New(factory system.CommandFactory, runner system.CommandRunner) Starter {

	factory.Register("start", StartCommand{
		CommandRunner: runner,
	})

	return Starter{}
}
