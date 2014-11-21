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
	SetLogger(logger *gosteno.Logger)
	GetLogger() *gosteno.Logger
}

type factory struct {
	commands map[string](Command)
	logger   *gosteno.Logger
}

func NewCommandFactory() *factory {
	return &factory{
		commands: make(map[string]Command),
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

func (f factory) FindCommand(name string) (command Command) {
	return f.commands[name]
}

func (f factory) Subcommands(cmd Command) (subcommands []Command) {
	switch cmd.(type) {
	case SubcommandProvider:
		return cmd.(SubcommandProvider).Subcommands()
	}
	return nil
}

func (f factory) SetLogger(logger *gosteno.Logger) {
	f.logger = logger
}

func (f factory) GetLogger() *gosteno.Logger {
	return f.logger
}
