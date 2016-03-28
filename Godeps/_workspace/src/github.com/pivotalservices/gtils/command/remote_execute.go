package command

import (
	"fmt"
	"io"
	"sync"

	"github.com/xchapter7x/lo"

	"golang.org/x/crypto/ssh"
)

// Config for the SSH connection
type SshConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	SSLKey   string
}

func (s *SshConfig) GetAuthMethod() (authMethod []ssh.AuthMethod) {

	if s.SSLKey == "" {
		lo.G.Debug("using password for authn")
		authMethod = []ssh.AuthMethod{
			ssh.Password(s.Password),
		}

	} else {
		lo.G.Debug("using sslkey for authn")
		keySigner, _ := ssh.ParsePrivateKey([]byte(s.SSLKey))

		authMethod = []ssh.AuthMethod{
			ssh.PublicKeys(keySigner),
		}
	}
	return
}

type ClientInterface interface {
	NewSession() (SSHSession, error)
}

type DefaultRemoteExecutor struct {
	Client         ClientInterface
	LazyClientDial func()
	once           sync.Once
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
		Auth: sshCfg.GetAuthMethod(),
	}
	remoteExecutor := &DefaultRemoteExecutor{}
	remoteExecutor.LazyClientDial = func() {
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshCfg.Host, sshCfg.Port), clientconfig)
		if err != nil {
			lo.G.Error("ssh connection issue:", err)
			return
		}
		remoteExecutor.Client = NewClientWrapper(client)
	}
	executor = remoteExecutor
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
	var session SSHSession
	var stdoutReader io.Reader

	if executor.once.Do(executor.LazyClientDial); executor.Client != nil {
		session, err = executor.Client.NewSession()
		defer session.Close()
		if err != nil {
			return
		}
		stdoutReader, err = session.StdoutPipe()
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
	} else {
		err = fmt.Errorf("un-initialized client executor")
		lo.G.Error(err.Error())
	}
	return
}
