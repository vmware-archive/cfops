package ssh

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

type SshConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

type SshCommand interface {
	execute(*ssh.Session) (err error)
}

func DialSsh(sshConfig *SshConfig, command SshCommand) (err error) {
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

	err = command.execute(session)
	if err != nil {
		err = errors.New("Failed to execute ssh command" + err.Error())
		return
	}
	return
}

type SshRemoteCopy struct {
	Command  string
	Filepath string
}

func (command *SshRemoteCopy) execute(session *ssh.Session) (err error) {

	f, err := os.Create(command.Filepath)
	defer f.Close()
	if err != nil {
		err = errors.New("Failed to create file" + err.Error())
		return
	}
	var b bytes.Buffer
	buf := make([]byte, 1024)
	session.Stdout = &b
	if err := session.Run(command.Command); err != nil {
		err = errors.New("Failed to run command: " + err.Error())
		return err
	}
	for {
		n, err := b.Read(buf)
		if err == io.EOF {
			break
		}
		_, err = f.Write(buf[:n])
		if err != nil {
			err = errors.New("Failed to wrtie to file" + err.Error())
			return err
		}
	}
	return
}
