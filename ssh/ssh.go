package ssh

import (
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// Config for the SSH connection
type Config struct {
	Username string
	Password string
	Host     string
	Port     int
}

type Copier interface {
	Copy(dest io.Writer, src io.Reader) error
}

type sshCopier struct {
	config *Config
}

func New(username string, password string, host string, port int) *sshCopier {
	copier := &sshCopier{
		config: &Config{
			Username: username,
			Password: password,
			Host:     host,
			Port:     port,
		},
	}
	return copier
}

// Copy the output from a command to the specified io.Writer
func (copier *sshCopier) Copy(dest io.Writer, src io.Reader) error {
	// TODO: error if port <= 0
	clientconfig := &ssh.ClientConfig{
		User: copier.config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(copier.config.Password),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", copier.config.Host, copier.config.Port), clientconfig)
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	session.Stdout = dest
	command, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	if err := session.Run(string(command[:])); err != nil {
		return err
	}

	return nil
}
