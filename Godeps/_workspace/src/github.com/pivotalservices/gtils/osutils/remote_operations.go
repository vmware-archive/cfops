package osutils

import (
	"fmt"
	"io"

	"github.com/pivotalservices/gtils/command"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	REMOTE_IMPORT_PATH string = "/tmp/archive.backup"
)

func NewRemoteOperations(sshCfg command.SshConfig) *remoteOperations {
	return &remoteOperations{
		sshCfg:     sshCfg,
		remotePath: REMOTE_IMPORT_PATH,
	}
}

type remoteOperations struct {
	sshCfg     command.SshConfig
	remotePath string
}

func (s *remoteOperations) UploadFile(lfile io.Reader) (err error) {
	var rfile io.WriteCloser

	if rfile, err = s.GetRemoteFile(); err == nil {
		defer rfile.Close()
		_, err = io.Copy(rfile, lfile)
	}
	return
}

func (s *remoteOperations) SetPath(p string) {
	s.remotePath = p
}

func (s *remoteOperations) Path() string {
	return s.remotePath
}

func (s *remoteOperations) GetRemoteFile() (rfile io.WriteCloser, err error) {
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
