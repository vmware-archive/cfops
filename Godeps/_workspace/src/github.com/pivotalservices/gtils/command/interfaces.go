package command

import "io"

type Executer interface {
	Execute(destination io.Writer, command string) error
}
