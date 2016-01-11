package cfbackup

import (
	"fmt"
	"io"
	"os"

	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/persistence"
	"github.com/xchapter7x/goutil"
)

func init() {
	SetPGDumpUtilVersions()
}

func SetPGDumpUtilVersions() {
	switch os.Getenv(ER_VERSION_ENV_FLAG) {
	case ER_VERSION_16:
		persistence.PGDMP_DUMP_BIN = "/var/vcap/packages/postgres-9.4.2/bin/pg_dump"
		persistence.PGDMP_RESTORE_BIN = "/var/vcap/packages/postgres-9.4.2/bin/pg_restore"
	default:
		persistence.PGDMP_DUMP_BIN = "/var/vcap/packages/postgres/bin/pg_dump"
		persistence.PGDMP_RESTORE_BIN = "/var/vcap/packages/postgres/bin/pg_restore"
	}
}

const (
	SD_PRODUCT   string = "Product"
	SD_COMPONENT string = "Component"
	SD_IDENTITY  string = "Identity"
	SD_IP        string = "Ip"
	SD_USER      string = "User"
	SD_PASS      string = "Pass"
	SD_VCAPUSER  string = "VcapUser"
	SD_VCAPPASS  string = "VcapPass"
)

type (
	PersistanceBackup interface {
		Dump(io.Writer) error
		Import(io.Reader) error
	}

	stringGetterSetter interface {
		Get(string) string
		Set(string, string)
	}

	SystemDump interface {
		stringGetterSetter
		Error() error
		GetPersistanceBackup() (dumper PersistanceBackup, err error)
	}

	SystemInfo struct {
		goutil.GetSet
		Product   string
		Component string
		Identity  string
		Ip        string
		User      string
		Pass      string
		VcapUser  string
		VcapPass  string
	}

	PgInfo struct {
		SystemInfo
		Database string
	}

	MysqlInfo struct {
		SystemInfo
		Database string
	}

	NfsInfo struct {
		SystemInfo
	}
)

func (s *SystemInfo) Get(name string) string {
	return s.GetSet.Get(s, name).(string)
}

func (s *SystemInfo) Set(name string, val string) {
	s.GetSet.Set(s, name, val)
}

func (s *NfsInfo) GetPersistanceBackup() (dumper PersistanceBackup, err error) {
	return NewNFSBackup(s.Pass, s.Ip)
}

func (s *MysqlInfo) GetPersistanceBackup() (dumper PersistanceBackup, err error) {
	sshConfig := command.SshConfig{
		Username: s.VcapUser,
		Password: s.VcapPass,
		Host:     s.Ip,
		Port:     22,
	}
	return persistence.NewRemoteMysqlDump(s.User, s.Pass, sshConfig)
}

func (s *PgInfo) GetPersistanceBackup() (dumper PersistanceBackup, err error) {
	sshConfig := command.SshConfig{
		Username: s.VcapUser,
		Password: s.VcapPass,
		Host:     s.Ip,
		Port:     22,
	}
	return persistence.NewPgRemoteDump(2544, s.Database, s.User, s.Pass, sshConfig)
}

func (s *SystemInfo) GetPersistanceBackup() (dumper PersistanceBackup, err error) {
	panic("you have to extend SystemInfo and implement GetPersistanceBackup method on the child")
	return
}

func (s *SystemInfo) Error() (err error) {
	if s.Product == "" ||
		s.Component == "" ||
		s.Identity == "" ||
		s.Ip == "" ||
		s.User == "" ||
		s.Pass == "" ||
		s.VcapUser == "" ||
		s.VcapPass == "" {
		err = fmt.Errorf("invalid or incomplete system info object: %+v", s)
	}
	return
}
