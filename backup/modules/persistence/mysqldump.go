package persistence

import (
	"bytes"
	"fmt"

	"github.com/pivotalservices/cfops/command"
)

type MysqlDump struct {
	Ip       string
	Username string
	Password string
	DbFile   string
	Caller   command.CmdExecuter
}

func NewMysqlDump(ip, username, password, dbFile string) *MysqlDump {
	return &MysqlDump{
		Ip:       ip,
		Username: username,
		Password: password,
		DbFile:   dbFile,
		Caller:   command.NewLocalExecuter(),
	}
}

func (s *MysqlDump) Dump() (err error) {

	if err = s.setupConfigFile(); err == nil {
		err = s.executeDumpToFile()
	}
	return
}

func (s *MysqlDump) setupConfigFile() (err error) {
	var b bytes.Buffer
	err = s.Caller.Execute(&b, s.getConfigCommand())
	return
}

func (s *MysqlDump) executeDumpToFile() (err error) {
	var b bytes.Buffer
	err = s.Caller.Execute(&b, s.getDumpCommand())
	return
}

func (s *MysqlDump) getDumpCommand() string {
	return fmt.Sprintf("mysqldump -u %s -h %s --all-databases > %s",
		s.Username,
		s.Ip,
		s.DbFile,
	)
}

func (s *MysqlDump) getConfigCommand() string {
	formatString := `echo "[mysqldump]
user=%s
password=%s" > ~/.my.cnf
`
	return fmt.Sprintf(formatString,
		s.Username,
		s.Password,
	)
}
