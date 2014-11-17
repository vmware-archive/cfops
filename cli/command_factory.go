package cli

import (
	"github.com/cloudfoundry/gosteno"
)

type SubcommandProvider interface {
	Subcommands() []Command
}

type CommandFactory interface {
	Register(commandName string, cmd Command) CommandFactory
	Commands() []Command
	Subcommands(command Command) []Command
}

type factory struct {
	commands map[string](Command)
	logger   *gosteno.Logger
}

func NewCommandFactory(logger *gosteno.Logger) *factory {
	return &factory{
		commands: make(map[string]Command),
		logger:   logger,
	}
}

func (f factory) Register(commandName string, cmd Command) CommandFactory {
	f.commands[commandName] = cmd
	return f
}

func (f factory) Commands() (commands []Command) {
	for _, command := range f.commands {
		commands = append(commands, command)
	}
	return
}

func (f factory) Subcommands(cmd Command) (subcommands []Command) {
	switch cmd.(type) {
	case SubcommandProvider:
		return cmd.(SubcommandProvider).Subcommands()
	}
	return nil
}
