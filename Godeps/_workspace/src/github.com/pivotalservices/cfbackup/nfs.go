package cfbackup

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/osutils"
)

//NewNFSBackup - constructor for an nfsbackup object
func NewNFSBackup(password, ip string, sslKey string, remoteArchivePath string) (nfs *NFSBackup, err error) {
	config := command.SshConfig{
		Username: NfsDefaultSSHUser,
		Password: password,
		Host:     ip,
		Port:     22,
		SSLKey:   sslKey,
	}
	var remoteExecuter command.Executer

	if remoteExecuter, err = NfsNewRemoteExecuter(config); err == nil {
		nfs = &NFSBackup{
			Caller:    remoteExecuter,
			RemoteOps: osutils.NewRemoteOperationsWithPath(config, remoteArchivePath),
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
		if err = s.Caller.Execute(ioutil.Discard, s.getRestoreCommand()); err == nil {
		    err = s.RemoteOps.RemoveRemoteFile()
		}
	}
	return
}

func (s *NFSBackup) getRestoreCommand() string {
	return fmt.Sprintf("cd %s && tar zx %s", NfsDirPath, s.RemoteOps.Path())
}

func (s *NFSBackup) getDumpCommand() string {
	return fmt.Sprintf("cd %s && tar cz %s", NfsDirPath, NfsArchiveDir)
}
