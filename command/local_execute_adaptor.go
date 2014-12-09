package command

import (
	"io"
	"os/exec"
	"strings"
)

type localExecuteAdaptor func(name string, arg ...string) *exec.Cmd

func (cmd localExecuteAdaptor) Execute(destination io.Writer, command string) (err error) {
	commandArr := strings.Split(command, " ")
	c := cmd(commandArr[0], commandArr[1:]...)
	c.Stdout = destination
	err = c.Run()
	return
}

func NewLocalExecuter() CmdExecuter {
	return localExecuteAdaptor(exec.Command)
}
