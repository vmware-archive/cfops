package install

import (
	"github.com/pivotalservices/cfops/system"
)

type MoveCommand struct {
	CommandRunner system.CommandRunner
	Installer
}

func (cmd MoveCommand) Metadata() system.CommandMetadata {
	return system.CommandMetadata{
		Name:        "move",
		Usage:       "move an existing deployment to another iaas location",
		Description: "destroy an existing cloud foundry foundation deployment from the iaas",
	}
}

func (cmd MoveCommand) Run(args []string) (err error) {
	println("move deployment: " + args[0])
	return
}
