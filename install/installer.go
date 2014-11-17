package install

import (
	"github.com/pivotalservices/cfops/cli"
	"github.com/pivotalservices/cfops/system"
)

type Installer struct {
	Commands []cli.Command
}

func New(factory cli.CommandFactory, runner system.CommandRunner) Installer {

	installer := Installer{
		Commands: []cli.Command{
			AddCommand{},
			DestroyCommand{},
			DumpCommand{},
			MoveCommand{},
		},
	}

	factory.Register("install", installer)
	return installer
}

func (cmd Installer) Metadata() cli.CommandMetadata {
	return cli.CommandMetadata{
		Name:        "install",
		ShortName:   "in",
		Usage:       "install cloud foundry to an iaas",
		Description: "Begin the installation of Cloud Foundry to a selected iaas",
	}
}

func (installer Installer) Subcommands() []cli.Command {
	return installer.Commands
}

func (cmd Installer) Run(args []string) error {
	return nil
}
