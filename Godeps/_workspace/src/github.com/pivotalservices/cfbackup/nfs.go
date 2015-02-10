package cfbackup

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/osutils"
)

const (
	NFS_DIR_PATH         string = "/var/vcap/store"
	NFS_ARCHIVE_DIR      string = "shared"
	NFS_DEFAULT_SSH_USER string = "vcap"
)

type remoteOpsInterface interface {
	UploadFile(lfile io.Reader) (err error)
	Path() string
}

func BackupNfs(password, ip string, dest io.Writer) (err error) {
	var nfsb *NFSBackup

	if nfsb, err = NewNFSBackup(password, ip); err == nil {
		err = nfsb.Dump(dest)
	}
	return
}

type NFSBackup struct {
	Caller    command.Executer
	RemoteOps remoteOpsInterface
}

var NfsNewRemoteExecuter func(command.SshConfig) (command.Executer, error) = command.NewRemoteExecutor

func NewNFSBackup(password, ip string) (nfs *NFSBackup, err error) {
	config := command.SshConfig{
		Username: NFS_DEFAULT_SSH_USER,
		Password: password,
		Host:     ip,
		Port:     22,
	}
	var remoteExecuter command.Executer

	if remoteExecuter, err = NfsNewRemoteExecuter(config); err == nil {
		nfs = &NFSBackup{
			Caller:    remoteExecuter,
			RemoteOps: osutils.NewRemoteOperations(config),
		}
	}
	return
}

func (s *NFSBackup) Dump(dest io.Writer) (err error) {
	err = s.Caller.Execute(dest, s.getDumpCommand())
	return
}

func (s *NFSBackup) Import(lfile io.Reader) (err error) {
	if err = s.RemoteOps.UploadFile(lfile); err == nil {
		err = s.Caller.Execute(ioutil.Discard, s.getRestoreCommand())
	}
	return
}

func (s *NFSBackup) getRestoreCommand() string {
	return fmt.Sprintf("cd %s && tar zx %s", NFS_DIR_PATH, s.RemoteOps.Path())
}

func (s *NFSBackup) getDumpCommand() string {
	return fmt.Sprintf("cd %s && tar cz %s", NFS_DIR_PATH, NFS_ARCHIVE_DIR)
}
