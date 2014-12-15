package backup

import (
	"io"

	"github.com/pivotalservices/cfops/command"
	"github.com/pivotalservices/cfops/osutils"
	"github.com/pivotalservices/cfops/utils"
)

func backupNfs(jsonfile, destDir string) (err error) {
	arguments := []string{jsonfile, "cf", "nfs_server", "vcap"}
	password := utils.GetPassword(arguments)
	ip := utils.GetIP(arguments)
	file, _ := osutils.SafeCreate(destDir, "nfs.tar.gz")
	defer file.Close()

	if nfsb, err := NewNFSBackup(password, ip); err == nil {
		err = nfsb.Dump(file)
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

func (s *NFSBackup) Dump(dest io.Writer) (err error) {
	command := "cd /var/vcap/store && tar cz shared"
	err = s.Caller.Execute(dest, command)
	return
}
