package system

import "github.com/codegangsta/cli"

type CommandMetadata struct {
	Name        string
	ShortName   string
	Usage       string
	Description string
	Flags       []cli.Flag
}

type Command interface {
	Subcommands() []Command
	Metadata() CommandMetadata
	Run(args []string) error
}
