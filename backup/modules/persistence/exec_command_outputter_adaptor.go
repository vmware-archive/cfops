package persistence

import (
	"os/exec"
	"strings"
)

type ExecCommandOutputterAdaptor func(name string, arg ...string) *exec.Cmd

func (cmd ExecCommandOutputterAdaptor) Output(cmdString string) ([]byte, error) {
	commandArr := strings.Split(cmdString, " ")
	return cmd(commandArr[0], commandArr[1:]...).Output()
}
