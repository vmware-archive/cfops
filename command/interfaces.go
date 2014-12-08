package command

import "io"

type CmdExecuter interface {
	Execute(destination io.Writer, command string) error
}
