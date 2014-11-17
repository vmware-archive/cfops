package install

import (
	"github.com/pivotalservices/cfops/cli"
)

type MoveCommand struct{}

func (cmd MoveCommand) Metadata() cli.CommandMetadata {
	return cli.CommandMetadata{
		Name:        "move",
		Usage:       "move an existing deployment to another iaas location",
		Description: "destroy an existing cloud foundry foundation deployment from the iaas",
	}
}

func (cmd MoveCommand) Run(args []string) (err error) {
	println("move deployment: " + args[0])
	return
}
