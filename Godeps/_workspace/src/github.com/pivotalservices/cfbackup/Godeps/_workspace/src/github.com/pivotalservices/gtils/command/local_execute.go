package command

import (
	"io"
	"os/exec"
	"strings"
)

type localExecute func(name string, arg ...string) *exec.Cmd

func (cmd localExecute) Execute(destination io.Writer, command string) (err error) {
	commandArr := strings.Split(command, " ")
	c := cmd(commandArr[0], commandArr[1:]...)
	c.Stdout = destination
	err = c.Run()
	return
}

func NewLocalExecuter() Executer {
	return localExecute(exec.Command)
}
