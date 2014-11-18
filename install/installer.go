package install

import (
	"github.com/cloudfoundry/gosteno"
	"github.com/pivotalservices/cfops/cli"
)

type Installer struct {
	Commands []cli.Command
	Logger   *gosteno.Logger
}

func New(factory cli.CommandFactory) Installer {

	installer := Installer{
		Commands: []cli.Command{
			AddCommand{
				Logger: factory.GetLogger(),
			},
			DestroyCommand{
				Logger: factory.GetLogger(),
			},
			DumpCommand{
			// Logger: factory.GetLogger(),
			},
			MoveCommand{
				Logger: factory.GetLogger(),
			},
		},
		Logger: factory.GetLogger(),
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
