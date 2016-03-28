package osutils

import (
	"fmt"
	"io"

	"github.com/pivotalservices/gtils/command"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	//RemoteImportPath - path to place the temporary remote import files when uploaded
	RemoteImportPath string = "/tmp/archive.backup"
)

//NewRemoteOperations - a constructor for a remoteoperations object
func NewRemoteOperations(sshCfg command.SshConfig) *RemoteOperations {
	return NewRemoteOperationsWithPath(sshCfg, RemoteImportPath)
}

func NewRemoteOperationsWithPath(sshCfg command.SshConfig, remoteImportPath string) *RemoteOperations {
	if len(remoteImportPath) == 0 {
		panic("remoteImportPath cannot be blank")
	}
	return &RemoteOperations{
		sshCfg:           sshCfg,
		remotePath:       remoteImportPath,
		GetSSHConnection: newSSHConnection,
	}
}

//RemoteOperations - an object which allows us to execute operations on a remote system
type RemoteOperations struct {
	sshCfg           command.SshConfig
	remotePath       string
	GetSSHConnection func(command.SshConfig, *ssh.ClientConfig) (SFTPClient, error)
}

//UploadFile - allows us to upload the contents of the given reader
func (s *RemoteOperations) UploadFile(lfile io.Reader) (err error) {
	var rfile io.WriteCloser

	if rfile, err = s.GetRemoteFile(); err == nil {
		defer rfile.Close()
		_, err = io.Copy(rfile, lfile)
	}
	return
}

//SetPath - allows us to set the remote path of the upload
func (s *RemoteOperations) SetPath(p string) {
	s.remotePath = p
}

//Path - allows us to get the path of the remote upload
func (s *RemoteOperations) Path() string {
	return s.remotePath
}

func (s *RemoteOperations) getClient() (sftpclient SFTPClient, err error) {
	clientconfig := &ssh.ClientConfig{
		User: s.sshCfg.Username,
		Auth: s.sshCfg.GetAuthMethod(),
	}
	sftpclient, err = s.GetSSHConnection(s.sshCfg, clientconfig)
	return
}

func newSSHConnection(config command.SshConfig, clientConfig *ssh.ClientConfig) (sftpclient SFTPClient, err error) {
	var sshconn *ssh.Client
	if sshconn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), clientConfig); err == nil {
		sftpclient, err = sftp.NewClient(sshconn)
	}
	return
}

//Remove Remote File - get a file from a remote system and return a writecloser to it
func (s *RemoteOperations) RemoveRemoteFile() (err error) {
	var sftpclient SFTPClient
	sftpclient, err = s.getClient()

	if err == nil {
		err = SafeRemoveSSH(sftpclient, s.remotePath)
	}

	return
}

//GetRemoteFile - get a file from a remote system and return a writecloser to it
func (s *RemoteOperations) GetRemoteFile() (rfile io.WriteCloser, err error) {
	var sftpclient SFTPClient
	sftpclient, err = s.getClient()

	if err == nil {
		rfile, err = SafeCreateSSH(sftpclient, s.remotePath)
	}

	return
}
