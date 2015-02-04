package persistence

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/osutils"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	PGDMP_REMOTE_IMPORT_PATH string = "/tmp/pgdump.sql"
	PGDMP_DUMP_BIN           string = "pg_dump"
	PGDMP_SQL_BIN                   = "psql"
	PGDMP_DROP_CMD                  = "drop schema public cascade;"
	PGDMP_CREATE_CMD                = "create schema public;"
)

type PgDump struct {
	sshCfg        command.SshConfig
	Ip            string
	Port          int
	Database      string
	Username      string
	Password      string
	DbFile        string
	Caller        command.Executer
	GetRemoteFile func(command.SshConfig) (io.WriteCloser, error)
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
		sshCfg:        sshCfg,
		Ip:            "localhost",
		Port:          port,
		Database:      database,
		Username:      username,
		Password:      password,
		Caller:        remoteExecuter,
		GetRemoteFile: getRemoteFile,
	}, err
}

func (s *PgDump) Import(lfile io.Reader) (err error) {

	if err = s.uploadBackupFile(lfile); err == nil {
		err = s.restore()
	}
	return
}

func (s *PgDump) restore() (err error) {
	var byteWriter bytes.Buffer

	callList := []string{
		s.getDropCommand(),
		s.getCreateCommand(),
		s.getImportCommand(),
	}

	for _, callstring := range callList {

		if err = s.Caller.Execute(&byteWriter, callstring); err != nil {
			break
		}
	}
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
	return fmt.Sprintf("%s < %s", connect, PGDMP_REMOTE_IMPORT_PATH)
}

func (s *PgDump) uploadBackupFile(lfile io.Reader) (err error) {
	var rfile io.WriteCloser

	if rfile, err = s.GetRemoteFile(s.sshCfg); err == nil {
		defer rfile.Close()
		_, err = io.Copy(rfile, lfile)
	}
	return
}

func (s *PgDump) Dump(dest io.Writer) (err error) {
	err = s.Caller.Execute(dest, s.getDumpCommand())
	return
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

func getRemoteFile(sshCfg command.SshConfig) (rfile io.WriteCloser, err error) {
	var (
		sshconn    *ssh.Client
		sftpclient *sftp.Client
	)

	clientconfig := &ssh.ClientConfig{
		User: sshCfg.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshCfg.Password),
		},
	}

	if sshconn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshCfg.Host, sshCfg.Port), clientconfig); err == nil {

		if sftpclient, err = sftp.NewClient(sshconn); err == nil {
			rfile, err = osutils.SafeCreateSSH(sftpclient, PGDMP_REMOTE_IMPORT_PATH)
		}
	}
	return
}
