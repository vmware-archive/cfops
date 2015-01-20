package command

import (
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

// Config for the SSH connection
type SshConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}

type ClientInterface interface {
	NewSession() (SSHSession, error)
}

type DefaultRemoteExecutor struct {
	Client ClientInterface
}

//Wrapper of ssh client to match client interface signature, since client.NewSession() does not use an interface
type SshClientWrapper struct {
	sshclient *ssh.Client
}

func NewClientWrapper(client *ssh.Client) *SshClientWrapper {
	return &SshClientWrapper{
		sshclient: client,
	}
}

func (c *SshClientWrapper) NewSession() (SSHSession, error) {
	return c.sshclient.NewSession()
}

// This method creates executor based on ssh, it has concrete ssh reference
func NewRemoteExecutor(sshCfg SshConfig) (executor Executer, err error) {
	clientconfig := &ssh.ClientConfig{
		User: sshCfg.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshCfg.Password),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshCfg.Host, sshCfg.Port), clientconfig)
	if err != nil {
		return
	}
	c := NewClientWrapper(client)
	executor = &DefaultRemoteExecutor{
		Client: c,
	}
	return
}

type SSHSession interface {
	Start(cmd string) error
	Wait() error
	StdoutPipe() (io.Reader, error)
	Close() error
}

// Copy the output from a command to the specified io.Writer
func (executor *DefaultRemoteExecutor) Execute(dest io.Writer, command string) (err error) {
	session, err := executor.Client.NewSession()
	defer session.Close()
	if err != nil {
		return
	}
	stdoutReader, err := session.StdoutPipe()
	if err != nil {
		return
	}
	err = session.Start(command)
	if err != nil {
		return
	}
	_, err = io.Copy(dest, stdoutReader)
	if err != nil {
		return
	}
	err = session.Wait()
	return
}
