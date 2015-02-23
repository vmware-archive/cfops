package command

import (
	"fmt"
	"io"
	"log"

	"github.com/dynport/gossh"
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

type GoSshExecutor struct {
	client *gossh.Client
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

func NewSshExecutor(sshCfg SshConfig) (Executer, error) {
	client := &gossh.Client{
		Host: sshCfg.Host,
		User: sshCfg.Username,
	}
	client.SetPassword(sshCfg.Password)
	client.DebugWriter = MakeLogger("DEBUG")
	client.InfoWriter = MakeLogger("INFO ")
	client.ErrorWriter = MakeLogger("ERROR")
	return &GoSshExecutor{client}, nil
}

func (executor *GoSshExecutor) Execute(dest io.Writer, command string) (err error) {
	client := executor.client
	defer client.Close()
	rsp, err := client.Execute("uptime")
	if err != nil {
		client.ErrorWriter(err.Error())
	}
	client.InfoWriter(rsp.String())

	// rsp, e = client.Execute("echo -n $(cat /proc/loadavg); cat /does/not/exists")
	rsp, err = client.Execute(command)
	if err != nil {
		client.ErrorWriter(err.Error())
		client.ErrorWriter("STDOUT: " + rsp.Stdout())
		client.ErrorWriter("STDERR: " + rsp.Stderr())
	}
	client.InfoWriter(rsp.String())
	return
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
	if exitError, ok := err.(*ssh.ExitError); ok {
		fmt.Printf("error status: %s", exitError.ExitStatus())
	}
	return
}

// returns a function of type gossh.Writer func(...interface{})
// MakeLogger just adds a prefix (DEBUG, INFO, ERROR)
func MakeLogger(prefix string) gossh.Writer {
	return func(args ...interface{}) {
		log.Println((append([]interface{}{prefix}, args...))...)
	}
}
