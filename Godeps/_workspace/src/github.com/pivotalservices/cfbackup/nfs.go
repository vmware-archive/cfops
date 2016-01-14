package cfbackup

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/osutils"
)

const (
	//NfsDirPath - this is where the nfs store lives
	NfsDirPath string = "/var/vcap/store"
	//NfsArchiveDir - this is the archive dir name
	NfsArchiveDir string = "shared"
	//NfsDefaultSSHUser - this is the default ssh user for nfs
	NfsDefaultSSHUser string = "vcap"
)

type remoteOpsInterface interface {
	UploadFile(lfile io.Reader) (err error)
	Path() string
}

//BackupNfs - this function will execute the nfs backup process
func BackupNfs(password, ip string, dest io.Writer) (err error) {
	var nfsb *NFSBackup

	if nfsb, err = NewNFSBackup(password, ip); err == nil {
		err = nfsb.Dump(dest)
	}
	return
}

//NFSBackup - this is a nfs backup object
type NFSBackup struct {
	Caller    command.Executer
	RemoteOps remoteOpsInterface
}

//NfsNewRemoteExecuter - this is a function which is able to execute a remote command against the nfs server
var NfsNewRemoteExecuter func(command.SshConfig) (command.Executer, error) = command.NewRemoteExecutor

//NewNFSBackup - constructor for an nfsbackup object
func NewNFSBackup(password, ip string) (nfs *NFSBackup, err error) {
	config := command.SshConfig{
		Username: NfsDefaultSSHUser,
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

//Dump - will dump the output of a executed command to the given writer
func (s *NFSBackup) Dump(dest io.Writer) (err error) {
	err = s.Caller.Execute(dest, s.getDumpCommand())
	return
}

//Import - will upload the contents of the given io.reader to the remote execution target and execute the restore command against the uploaded file.
func (s *NFSBackup) Import(lfile io.Reader) (err error) {
	if err = s.RemoteOps.UploadFile(lfile); err == nil {
		err = s.Caller.Execute(ioutil.Discard, s.getRestoreCommand())
	}
	return
}

func (s *NFSBackup) getRestoreCommand() string {
	return fmt.Sprintf("cd %s && tar zx %s", NfsDirPath, s.RemoteOps.Path())
}

func (s *NFSBackup) getDumpCommand() string {
	return fmt.Sprintf("cd %s && tar cz %s", NfsDirPath, NfsArchiveDir)
}
