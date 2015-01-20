package persistence

import (
	"fmt"
	"io"

	"github.com/pivotalservices/gtils/command"
)

type PgDump struct {
	Ip       string
	Port     int
	Database string
	Username string
	Password string
	DbFile   string
	Caller   command.Executer
}

type Dumper interface {
	Dump(io.Writer) error
}

func NewPgDump(ip string, port int, database, username, password string) Dumper {
	return &PgDump{
		Ip:       ip,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
		Caller:   command.NewLocalExecuter(),
	}
}

func NewPgRemoteDump(port int, database, username, password string, sshCfg command.SshConfig) (Dumper, error) {
	remoteExecuter, err := command.NewRemoteExecutor(sshCfg)
	return &PgDump{
		Ip:       "localhost",
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
		Caller:   remoteExecuter,
	}, err
}

func (s *PgDump) Dump(dest io.Writer) (err error) {
	err = s.Caller.Execute(dest, s.getDumpCommand())
	return
}

func (s *PgDump) getDumpCommand() string {
	return fmt.Sprintf("PGPASSWORD=%s pg_dump -h %s -U %s -p %d %s",
		s.Password,
		s.Ip,
		s.Username,
		s.Port,
		s.Database,
	)
}
