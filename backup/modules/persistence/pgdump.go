package persistence

import (
	"fmt"
	"io"
	"os"

	"github.com/pivotalservices/cfops/command"
)

type PgDump struct {
	Ip       string
	Port     int
	Database string
	Username string
	Password string
	DbFile   string
	Caller   command.CmdExecuter
}

func NewPgDump(ip string, port int, database, username, password string) *PgDump {
	return &PgDump{
		Ip:       ip,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
		Caller:   command.NewLocalExecuter(),
	}
}

func (s *PgDump) Dump(dest io.Writer) (err error) {

	if err = os.Setenv("PGPASSWORD", s.Password); err == nil {
		err = s.executeDumpToWriter(dest)
	}
	return
}

func (s *PgDump) executeDumpToWriter(dest io.Writer) (err error) {
	err = s.Caller.Execute(dest, s.getDumpCommand())
	return
}

func (s *PgDump) getDumpCommand() string {
	return fmt.Sprintf("pg_dump -h %s -U %s -p %s %s",
		s.Ip,
		s.Username,
		s.Port,
		s.Database,
	)
}
