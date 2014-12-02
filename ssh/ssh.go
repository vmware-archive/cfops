package ssh

import (
	"bufio"
	"bytes"
	"errors"
	"golang.org/x/crypto/ssh"
	"io"
)

type SshConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

type DumpOutput interface {
	Execute(io.Reader) (err error)
}

func DialSsh(sshConfig *SshConfig, command string, output DumpOutput) (err error) {
	config := &ssh.ClientConfig{
		User: sshConfig.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshConfig.Password),
		},
	}
	client, err := ssh.Dial("tcp", sshConfig.Host+":"+sshConfig.Port, config)
	if err != nil {
		err = errors.New("Failed to dial: " + err.Error())
		return
	}
	session, err := client.NewSession()
	if err != nil {
		err = errors.New("Failed to create session: " + err.Error())
		return
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		err = errors.New("Failed to run command: " + err.Error())
		return err
	}

	err = output.Execute(&b)
	if err != nil {
		err = errors.New("Failed to execute ssh command" + err.Error())
		return
	}
	return
}

type DumpToWriter struct {
	Writer io.Writer
}

func (w *DumpToWriter) Execute(r io.Reader) (err error) {
	reader := bufio.NewReader(r)
	_, err = (&reader).WriteTo(w.Writer)
	return
}
