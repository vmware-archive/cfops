package install

import (
	"github.com/pivotalservices/cfops/system"
)

type DestroyCommand struct {
	CommandRunner system.CommandRunner
	Installer
}

func (cmd DestroyCommand) Metadata() system.CommandMetadata {
	return system.CommandMetadata{
		Name:        "destroy",
		Usage:       "destroy an existing deployment",
		Description: "destroy an existing cloud foundry foundation deployment from the iaas",
	}
}

func (cmd DestroyCommand) HasSubcommands() bool {
	return false
}

func (cmd DestroyCommand) Run(args []string) (err error) {
	println("destroyed deployment: " + args[0])
	return
}
