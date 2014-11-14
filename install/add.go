package install

import (
	"github.com/pivotalservices/cfops/system"
)

type AddCommand struct {
	CommandRunner system.CommandRunner
	Installer
}

func (cmd AddCommand) Metadata() system.CommandMetadata {
	return system.CommandMetadata{
		Name:        "add",
		Usage:       "add a new deployment",
		Description: "use the provided deployment template to deploy a new cloud foundry foundation into the iaas",
	}
}

func (cmd AddCommand) Run(args []string) (err error) {
	println("new deployment with template: " + args[0])
	return
}
