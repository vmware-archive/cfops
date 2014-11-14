package install

import (
	"github.com/pivotalservices/cfops/system"
)

type Installer struct {
	Commands []system.Command
}

func New(factory system.CommandFactory, runner system.CommandRunner) Installer {

	installer := Installer{
		Commands: []system.Command{
			AddCommand{
				CommandRunner: runner,
			},
			DestroyCommand{
				CommandRunner: runner,
			},
			DumpCommand{
				CommandRunner: runner,
			},
			MoveCommand{
				CommandRunner: runner,
			},
		},
	}

	factory.Register("install", installer)
	return installer
}

func (cmd Installer) Metadata() system.CommandMetadata {
	return system.CommandMetadata{
		Name:        "install",
		ShortName:   "in",
		Usage:       "install cloud foundry to an iaas",
		Description: "Begin the installation of Cloud Foundry to a selected iaas",
	}
}

func (cmd Installer) Subcommands() (commands []system.Command) {
	return cmd.Commands
}

func (cmd Installer) Run(args []string) (err error) {
	return
}
