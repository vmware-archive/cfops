package system

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cloudfoundry/gosteno"
)

type CommandRunner interface {
	Run(name string, args ...string) error
	SetLogger(logger *gosteno.Logger)
}

type OSCommandRunner struct {
	Logger *gosteno.Logger
}

func (runner *OSCommandRunner) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	runner.Logger.Info(fmt.Sprint(name, " ", strings.Join(args, " ")))
	err := cmd.Run()
	runner.Logger.Info(fmt.Sprintf("result: %s", out.String()))
	return err
}

func (runner *OSCommandRunner) SetLogger(logger *gosteno.Logger) {
	runner.Logger = logger
}
