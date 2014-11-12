package system

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cloudfoundry/gosteno"
)

type CommandRunner interface {
	Run(name string, args ...string) error
}

type OSCommandRunner struct {
	Logger *gosteno.Logger
}

func (runner OSCommandRunner) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	runner.Logger.Info(fmt.Sprint(name, " ", strings.Join(args, " ")))
	return cmd.Run()
}
