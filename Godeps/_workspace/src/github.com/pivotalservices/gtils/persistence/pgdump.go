package persistence

import (
	"fmt"
	"io"

	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/osutils"
	"github.com/xchapter7x/lo"
)

const (
	PGDMP_DROP_CMD   string = "drop schema public cascade;"
	PGDMP_CREATE_CMD        = "create schema public;"
)

var (
	PGDMP_DUMP_BIN string = "/var/vcap/packages/postgres/bin/pg_dump"
	PGDMP_SQL_BIN         = "/var/vcap/packages/postgres/bin/psql -v ON_ERROR_STOP=1"
)

type PgDump struct {
	sshCfg    command.SshConfig
	Ip        string
	Port      int
	Database  string
	Username  string
	Password  string
	DbFile    string
	Caller    command.Executer
	RemoteOps remoteOperationsInterface
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

func NewPgRemoteDump(port int, database, username, password string, sshCfg command.SshConfig) (*PgDump, error) {
	remoteExecuter, err := command.NewRemoteExecutor(sshCfg)
	return &PgDump{
		sshCfg:    sshCfg,
		Ip:        "localhost",
		Port:      port,
		Database:  database,
		Username:  username,
		Password:  password,
		Caller:    remoteExecuter,
		RemoteOps: osutils.NewRemoteOperations(sshCfg),
	}, err
}

func (s *PgDump) Import(lfile io.Reader) (err error) {
	lo.G.Debug("pgdump Import being called")

	if err = s.RemoteOps.UploadFile(lfile); err == nil {
		err = s.restore()
	}
	return
}

func (s *PgDump) restore() (err error) {
	callList := []string{
		s.getAddSuperUserCommand(),
		s.getDropCommand(),
		s.getCreateCommand(),
		s.getImportCommand(),
	}
	err = execute_list(callList, s.Caller)
	lo.G.Debug("pgdump restore called: ", callList, err)
	return
}

func (s *PgDump) getDropCommand() string {
	connect := s.getPostgresConnect(PGDMP_SQL_BIN)
	return fmt.Sprintf("%s -c '%s'", connect, PGDMP_DROP_CMD)
}

func (s *PgDump) getCreateCommand() string {
	connect := s.getPostgresConnect(PGDMP_SQL_BIN)
	return fmt.Sprintf("%s -c '%s'", connect, PGDMP_CREATE_CMD)
}

func (s *PgDump) getImportCommand() string {
	connect := s.getPostgresConnect(PGDMP_SQL_BIN)
	return fmt.Sprintf("%s < %s", connect, s.RemoteOps.Path())
}

func (s *PgDump) Dump(dest io.Writer) (err error) {
	err = s.Caller.Execute(dest, s.getDumpCommand())
	lo.G.Debug("pgdump Dump called: ", s.getDumpCommand(), err)
	return
}

func (s *PgDump) getAddSuperUserCommand() string {
	return fmt.Sprintf("echo \"%s\" | sudo -u postgres %s -c 'alter role %s superuser;'",
		s.Password,
		PGDMP_SQL_BIN,
		s.Username,
	)
}

func (s *PgDump) getPostgresConnect(command string) string {
	return fmt.Sprintf("PGPASSWORD=%s %s -h %s -U %s -p %d %s",
		s.Password,
		command,
		s.Ip,
		s.Username,
		s.Port,
		s.Database,
	)
}

func (s *PgDump) getDumpCommand() string {
	return s.getPostgresConnect(PGDMP_DUMP_BIN)
}
