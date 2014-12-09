package persistence

import (
	"fmt"
	"io"
	"os"

	"github.com/pivotalservices/cfops/command"
	"github.com/pivotalservices/cfops/osutils"
)

type MysqlDump struct {
	Ip         string
	Username   string
	Password   string
	DbFile     string
	ConfigFile string
	Caller     command.CmdExecuter
}

func NewMysqlDump(ip, username, password, dbFile string) *MysqlDump {
	m := &MysqlDump{
		Ip:         ip,
		Username:   username,
		Password:   password,
		DbFile:     dbFile,
		ConfigFile: "~/.my.cnf",
		Caller:     command.NewLocalExecuter(),
	}
	return m
}

func (s *MysqlDump) Dump(dest io.Writer) (err error) {

	if err = s.setupConfigFile(); err == nil {
		err = s.executeDumpToWriter(dest)
	}
	return
}

func (s *MysqlDump) setupConfigFile() (err error) {
	os.Remove(s.ConfigFile)
	b, err := osutils.SafeCreate(s.ConfigFile)
	defer b.Close()
	err = s.Caller.Execute(b, s.getConfigCommand())
	return
}

func (s *MysqlDump) executeDumpToWriter(dest io.Writer) (err error) {
	err = s.Caller.Execute(dest, s.getDumpCommand())
	return
}

func (s *MysqlDump) getDumpCommand() string {
	return fmt.Sprintf("mysqldump -u %s -h %s --all-databases",
		s.Username,
		s.Ip,
	)
}

func (s *MysqlDump) getConfigCommand() string {
	formatString := `echo "[mysqldump]
user=%s
password=%s"
`
	return fmt.Sprintf(formatString,
		s.Username,
		s.Password,
	)
}
