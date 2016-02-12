package persistence

import (
	"fmt"
	"io"

	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/osutils"
	"github.com/xchapter7x/lo"
)

//NewPgDump - a pgdump object initialized for local fs
func NewPgDump(ip string, port int, database, username, password string) *PgDump {
	return &PgDump{
		IP:       ip,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
		Caller:   command.NewLocalExecuter(),
	}
}

func NewPgRemoteDumpWithPath(port int, database, username, password string, sshCfg command.SshConfig, remoteArchivePath string) (*PgDump, error) {
	remoteExecuter, err := command.NewRemoteExecutor(sshCfg)
	remoteOps := osutils.NewRemoteOperations(sshCfg)
	if len(remoteArchivePath) > 0 {
		remoteOps.SetPath(remoteArchivePath)
	}
	return &PgDump{
		sshCfg:    sshCfg,
		IP:        "localhost",
		Port:      port,
		Database:  database,
		Username:  sshCfg.Username,
		Password:  sshCfg.Password,
		Caller:    remoteExecuter,
		RemoteOps: remoteOps,
	}, err
}

//NewPgRemoteDump - a pgdump initialized for remote fs
func NewPgRemoteDump(port int, database, username, password string, sshCfg command.SshConfig) (*PgDump, error) {
	return NewPgRemoteDumpWithPath(port, database, username, password, sshCfg, "")
}

//Import - allows us to import a pgdmp file in the form of a reader
func (s *PgDump) Import(lfile io.Reader) (err error) {
	lo.G.Debug("pgdump Import being called")

	if err = s.RemoteOps.UploadFile(lfile); err == nil {
		err = s.restore()
	}
	return
}

func (s *PgDump) restore() (err error) {
	callList := []string{
		s.getRestoreCommand(),
	}
	err = executeList(callList, s.Caller)
	lo.G.Debug("pgdump restore called: ", callList, err)
	return
}

//Dump - allows us to create a backup and pipe it to the given writer
func (s *PgDump) Dump(dest io.Writer) (err error) {
	err = s.Caller.Execute(dest, s.getDumpCommand())
	lo.G.Debug("pgdump Dump called: ", s.getDumpCommand(), err)
	return
}

func (s *PgDump) dumpConnectionDecorator(command string) string {
	return fmt.Sprintf("PGPASSWORD=%s %s -Fc -h %s -U %s -p %d %s",
		s.Password,
		command,
		s.IP,
		s.Username,
		s.Port,
		s.Database,
	)
}

func (s *PgDump) restoreConnectionDecorator(command string) string {
	return fmt.Sprintf("PGPASSWORD=%s %s -h %s -U %s -x -p %d -c -d %s %s",
		s.Password,
		command,
		s.IP,
		s.Username,
		s.Port,
		s.Database,
		s.RemoteOps.Path(),
	)
}

func (s *PgDump) getRestoreCommand() string {
	return s.restoreConnectionDecorator(PGDmpRestoreBin)
}

func (s *PgDump) getDumpCommand() string {
	return s.dumpConnectionDecorator(PGDmpDumpBin)
}
