package start

import (
	"github.com/pivotalservices/cfops/system"
)

type StartCommand struct {
	CommandRunner system.CommandRunner
	Starter
}

func (cmd StartCommand) Metadata() system.CommandMetadata {
	return system.CommandMetadata{
		Name:        "start",
		ShortName:   "s",
		Usage:       "start up an entire cloud foundry foundation",
		Description: "start all the VMs in an existing cloud foundry deployment",
	}
}

func (cmd StartCommand) Run(args []string) error {
	err := cmd.CommandRunner.Run("echo", "WHOOOA", "slow down!")
	if err != nil {
		return err
	}
	return nil
}

func (cmd StartCommand) Subcommands() (commands []system.Command) {
	return
}
