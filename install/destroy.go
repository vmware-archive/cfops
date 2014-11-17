package install

import (
	"github.com/pivotalservices/cfops/cli"
)

type DestroyCommand struct{}

func (cmd DestroyCommand) Metadata() cli.CommandMetadata {
	return cli.CommandMetadata{
		Name:        "destroy",
		Usage:       "destroy an existing deployment",
		Description: "destroy an existing cloud foundry foundation deployment from the iaas",
	}
}

func (cmd DestroyCommand) Run(args []string) (err error) {
	println("destroyed deployment: " + args[0])
	return
}
