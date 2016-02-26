package cfbackup

import (
	"fmt"

	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/persistence"
)

//Get - a getter for a systeminfo object
func (s *SystemInfo) Get(name string) string {
	return s.GetSet.Get(s, name).(string)
}

//Set - a setter for a systeminfo object
func (s *SystemInfo) Set(name string, val string) {
	s.GetSet.Set(s, name, val)
}

//GetPersistanceBackup - the constructor for a new nfsinfo object
func (s *NfsInfo) GetPersistanceBackup() (dumper PersistanceBackup, err error) {
	return NewNFSBackup(s.Pass, s.Ip, s.SSHPrivateKey, s.RemoteArchivePath)
}

//GetPersistanceBackup - the constructor for a new DirectorInfo object
func (s *DirectorInfo) GetPersistanceBackup() (dumper PersistanceBackup, err error) {
	sshConfig := command.SshConfig{
		Username: s.VcapUser,
		Password: s.VcapPass,
		Host:     s.Ip,
		Port:     22,
		SSLKey:   s.SSHPrivateKey,
	}
	return persistence.NewPgRemoteDumpWithPath(2544, s.Database, s.User, s.Pass, sshConfig, s.RemoteArchivePath)
}

//GetPersistanceBackup - the constructor for a new mysqlinfo object
func (s *MysqlInfo) GetPersistanceBackup() (dumper PersistanceBackup, err error) {
	sshConfig := command.SshConfig{
		Username: s.VcapUser,
		Password: s.VcapPass,
		Host:     s.Ip,
		Port:     22,
		SSLKey:   s.SSHPrivateKey,
	}
	return persistence.NewRemoteMysqlDumpWithPath(s.User, s.Pass, sshConfig, s.RemoteArchivePath)
}

//GetPersistanceBackup - the constructor for a new pginfo object
func (s *PgInfo) GetPersistanceBackup() (dumper PersistanceBackup, err error) {
	sshConfig := command.SshConfig{
		Username: s.VcapUser,
		Password: s.VcapPass,
		Host:     s.Ip,
		Port:     22,
		SSLKey:   s.SSHPrivateKey,
	}
	return persistence.NewPgRemoteDump(2544, s.Database, s.User, s.Pass, sshConfig)
}

//GetPersistanceBackup - the constructor for a systeminfo object
func (s *SystemInfo) GetPersistanceBackup() (dumper PersistanceBackup, err error) {
	panic("you have to extend SystemInfo and implement GetPersistanceBackup method on the child")
	return
}

//Error - method making systeminfo implement the error interface
func (s *SystemInfo) Error() (err error) {
	if s.Product == "" ||
		s.Component == "" ||
		s.Identifier == "" ||
		s.Ip == "" ||
		s.User == "" ||
		s.Pass == "" ||
		s.VcapUser == "" ||
		s.VcapPass == "" {
		err = fmt.Errorf("invalid or incomplete system info object: %+v", s)
	}
	return
}
