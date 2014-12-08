package command

import (
	"io"
	"os/exec"
	"strings"
)

type localExecuteAdaptor func(name string, arg ...string) *exec.Cmd

func (cmd localExecuteAdaptor) Execute(destination io.Writer, command string) error {
	commandArr := strings.Split(command, " ")
	byteArray, err := cmd(commandArr[0], commandArr[1:]...).Output()
	destination.Write(byteArray)
	return err
}

func NewLocalExecuter() CmdExecuter {
	return localExecuteAdaptor(exec.Command)
}
