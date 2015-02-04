package cfbackup

import (
	"io"

	"github.com/pivotalservices/gtils/command"
)

func BackupNfs(password, ip string, dest io.Writer) (err error) {
	var nfsb *NFSBackup

	if nfsb, err = NewNFSBackup(password, ip); err == nil {
		err = nfsb.Dump(dest)
	}
	return
}

type NFSBackup struct {
	Caller command.Executer
}

var NfsNewRemoteExecuter func(command.SshConfig) (command.Executer, error) = command.NewRemoteExecutor

func NewNFSBackup(password, ip string) (nfs *NFSBackup, err error) {
	config := command.SshConfig{
		Username: "vcap",
		Password: password,
		Host:     ip,
		Port:     22,
	}
	var remoteExecuter command.Executer

	if remoteExecuter, err = NfsNewRemoteExecuter(config); err == nil {
		nfs = &NFSBackup{
			Caller: remoteExecuter,
		}
	}
	return
}

func (s *NFSBackup) Import(io.Reader) (err error) {
	panic("you need to implement this")
	return
}

func (s *NFSBackup) Dump(dest io.Writer) (err error) {
	command := "cd /var/vcap/store && tar cz shared"
	err = s.Caller.Execute(dest, command)
	return
}
