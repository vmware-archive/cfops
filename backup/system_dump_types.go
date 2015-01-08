package backup

import (
	"fmt"

	"github.com/pivotalservices/cfops/backup/modules/persistence"
	"github.com/pivotalservices/cfops/command"
)

type (
	SystemDump interface {
		Error() error
		GetDumper() (dumper persistence.Dumper, err error)
		GetProduct() string
		GetComponent() string
		GetIdentity() string
		GetIp() string
		SetIp(string)
		GetUser() string
		SetUser(string)
		GetPass() string
		SetPass(string)
		GetVcapUser() string
		SetVcapUser(string)
		GetVcapPass() string
		SetVcapPass(string)
	}

	SystemInfo struct {
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
	}

	MysqlInfo struct {
		SystemInfo
	}

	NfsInfo struct {
		SystemInfo
	}
)

func (s *SystemInfo) GetProduct() string {
	return s.Product
}

func (s *SystemInfo) GetComponent() string {
	return s.Component
}

func (s *SystemInfo) GetIdentity() string {
	return s.Identity
}

func (s *SystemInfo) GetIp() string {
	return s.Ip
}

func (s *SystemInfo) SetIp(in string) {
	s.Ip = in
}

func (s *SystemInfo) GetUser() string {
	return s.User
}

func (s *SystemInfo) SetUser(in string) {
	s.User = in
}

func (s *SystemInfo) GetPass() string {
	return s.Pass
}

func (s *SystemInfo) SetPass(in string) {
	s.Pass = in
}

func (s *SystemInfo) GetVcapUser() string {
	return s.VcapUser
}

func (s *SystemInfo) SetVcapUser(in string) {
	s.VcapUser = in
}

func (s *SystemInfo) GetVcapPass() string {
	return s.VcapPass
}

func (s *SystemInfo) SetVcapPass(in string) {
	s.VcapPass = in
}

func (s *NfsInfo) GetDumper() (dumper persistence.Dumper, err error) {
	return NewNFSBackup(s.Pass, s.Ip)
}

func (s *MysqlInfo) GetDumper() (dumper persistence.Dumper, err error) {
	sshConfig := command.SshConfig{
		Username: s.VcapUser,
		Password: s.VcapPass,
		Host:     s.Ip,
		Port:     22,
	}
	return persistence.NewRemoteMysqlDump(s.User, s.Pass, sshConfig)
}

func (s *PgInfo) GetDumper() (dumper persistence.Dumper, err error) {
	sshConfig := command.SshConfig{
		Username: s.VcapUser,
		Password: s.VcapPass,
		Host:     s.Ip,
		Port:     22,
	}
	return persistence.NewPgRemoteDump(2544, s.Component, s.User, s.Pass, sshConfig)
}

func (s *SystemInfo) GetDumper() (dumper persistence.Dumper, err error) {
	panic("you have to extend SystemInfo and implement GetDumper method on the child")
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
		err = fmt.Errorf("invalid or incomplete system info object: %s", s)
	}
	return
}
