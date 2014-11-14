package start

import (
	"github.com/pivotalservices/cfops/system"
)

type Starter struct {
	Commands []system.Command
}

func New(factory system.CommandFactory, runner system.CommandRunner) Starter {

	factory.Register("start", StartCommand{
		CommandRunner: runner,
	})

	return Starter{}
}

func (cmd Starter) Subcommands() (commands []system.Command) {
	return
}
