package start

import (
	"github.com/pivotalservices/cfops/cli"
)

type StartCommand struct{}

func (cmd StartCommand) Metadata() cli.CommandMetadata {
	return cli.CommandMetadata{
		Name:        "start",
		ShortName:   "s",
		Usage:       "start up an entire cloud foundry foundation",
		Description: "start all the VMs in an existing cloud foundry deployment",
	}
}

func (cmd StartCommand) Run(args []string) error {
	return nil
}
