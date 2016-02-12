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
		sshCfg:     sshCfg,
		remotePath: remoteImportPath,
	}
}

//RemoteOperations - an object which allows us to execute operations on a remote system
type RemoteOperations struct {
	sshCfg     command.SshConfig
	remotePath string
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

//GetRemoteFile - get a file from a remote system and return a writecloser to it
func (s *RemoteOperations) GetRemoteFile() (rfile io.WriteCloser, err error) {
	var (
		sshconn    *ssh.Client
		sftpclient *sftp.Client
	)

	clientconfig := &ssh.ClientConfig{
		User: s.sshCfg.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.sshCfg.Password),
		},
	}

	if sshconn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.sshCfg.Host, s.sshCfg.Port), clientconfig); err == nil {

		if sftpclient, err = sftp.NewClient(sshconn); err == nil {
			rfile, err = SafeCreateSSH(sftpclient, s.remotePath)
		}
	}
	return
}
