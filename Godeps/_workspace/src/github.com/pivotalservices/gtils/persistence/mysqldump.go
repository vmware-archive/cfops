package persistence

import (
	"fmt"
	"io"

	"github.com/pivotalservices/gtils/command"
)

type MysqlDump struct {
	Ip         string
	Username   string
	Password   string
	DbFile     string
	ConfigFile string
	Caller     command.Executer
}

func NewMysqlDump(ip, username, password string) *MysqlDump {
	m := &MysqlDump{
		Ip:       ip,
		Username: username,
		Password: password,
		Caller:   command.NewLocalExecuter(),
	}
	return m
}

func NewRemoteMysqlDump(username, password string, sshCfg command.SshConfig) (*MysqlDump, error) {
	remoteExecuter, err := command.NewRemoteExecutor(sshCfg)
	return &MysqlDump{
		Ip:       "localhost",
		Username: username,
		Password: password,
		Caller:   remoteExecuter,
	}, err
}

func (s *MysqlDump) Dump(dest io.Writer) (err error) {
	err = s.Caller.Execute(dest, s.getDumpCommand())
	return
}

func (s *MysqlDump) getDumpCommand() string {
	return fmt.Sprintf("mysqldump -u %s -h %s --password=%s --all-databases",
		s.Username,
		s.Ip,
		s.Password,
	)
}
