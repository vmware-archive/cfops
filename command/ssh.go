package command

import (
	"bytes"
	"io"
	"io/ioutil"
)

// Config for the SSH connection
type Config struct {
	Username string
	Password string
	Host     string
	Port     int
}

// Copier copies from an io.Reader to an io.Writer
type Copier interface {
	Copy(dest io.Writer, src io.Reader) error
}

// DefaultCopier is an SSH copier
type DefaultCopier struct {
	session SSHSession
}

//func main() {
//clientconfig := &ssh.ClientConfig{
//User: copier.config.Username,
//Auth: []ssh.AuthMethod{
//ssh.Password(copier.config.Password),
//},
//}

//client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", copier.config.Host, copier.config.Port)
//session, err := client.NewSession()
//defer session.Close()
//sshExecuteRef := New(session)
//sshExecuteRef.Execute(writer)
//}

// New ssh copier
func New(session SSHSession) (copier *DefaultCopier) {
	copier = &DefaultCopier{
		session: session,
	}
	return
}

type SSHSession interface {
	Run(string) error
	StdoutPipe() (io.Reader, error)
}

func (copier *DefaultCopier) Copy(dest io.Writer, src io.Reader) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	s := buf.String()
	return copier.Execute(dest, s)
}

// Copy the output from a command to the specified io.Writer
func (copier *DefaultCopier) Execute(dest io.Writer, command string) (err error) {

	if stdoutReader, err := copier.session.StdoutPipe(); err == nil {

		if b, err := ioutil.ReadAll(stdoutReader); err == nil {
			dest.Write(b)
			err = copier.session.Run(command)
		}
	}
	return err
}
