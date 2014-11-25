package ssh

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

type executeCmd func(*ssh.Session) (err error)

type SshConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

func dialSsh(sshConfig *SshConfig, fn executeCmd) (err error) {
	config := &ssh.ClientConfig{
		User: sshConfig.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshConfig.Password),
		},
	}
	client, err := ssh.Dial("tcp", sshConfig.Host+":"+sshConfig.Port, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}

	defer session.Close()

	err = fn(session)
	if err != nil {
		panic("Failed to execute ssh command" + err.Error())
	}
	return
}

type SshRemoteCopy struct {
	SshConfig SshConfig
	Command   string
	Filepath  string
}

func (sshRemoteCopy *SshRemoteCopy) SshCmdCopy() (err error) {
	err = dialSsh(&sshRemoteCopy.SshConfig, func(session *ssh.Session) (err error) {
		f, err := os.Create(sshRemoteCopy.Filepath)
		defer f.Close()
		if err != nil {
			panic("Failed to create file" + err.Error())
		}
		var b bytes.Buffer
		buf := make([]byte, 1024)
		session.Stdout = &b
		if err := session.Run(sshRemoteCopy.Command); err != nil {
			panic("Failed to run: " + err.Error())
		}
		for {
			n, err := b.Read(buf)
			if err == io.EOF {
				break
			}
			_, err = f.Write(buf[:n])
			if err != nil {
				panic("Failed to wrtie to file" + err.Error())
			}

		}
		return
	})
	return
}
