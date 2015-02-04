package persistence

import (
	"fmt"
	"io"

	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/osutils"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func getRemoteFile(sshCfg command.SshConfig) (rfile io.WriteCloser, err error) {
	var (
		sshconn    *ssh.Client
		sftpclient *sftp.Client
	)

	clientconfig := &ssh.ClientConfig{
		User: sshCfg.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshCfg.Password),
		},
	}

	if sshconn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshCfg.Host, sshCfg.Port), clientconfig); err == nil {

		if sftpclient, err = sftp.NewClient(sshconn); err == nil {
			rfile, err = osutils.SafeCreateSSH(sftpclient, PGDMP_REMOTE_IMPORT_PATH)
		}
	}
	return
}
