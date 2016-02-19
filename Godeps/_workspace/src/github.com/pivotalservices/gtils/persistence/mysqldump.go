package persistence

import (
	"fmt"
	"io"

	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/osutils"
	"github.com/xchapter7x/lo"
)

//NewMysqlDump - will initialize a mysqldump for local execution
func NewMysqlDump(ip, username, password string) *MysqlDump {
	lo.G.Debug("setting up a new local mysqldump object")
	m := &MysqlDump{
		IP:       ip,
		Username: username,
		Password: password,
		Caller:   command.NewLocalExecuter(),
	}
	return m
}

func NewRemoteMysqlDumpWithPath(username, password string, sshCfg command.SshConfig, remoteArchivePath string) (*MysqlDump, error) {
    lo.G.Debug("setting up a new remote MyslDump object")
	remoteExecuter, err := command.NewRemoteExecutor(sshCfg)

    remoteOps := osutils.NewRemoteOperations(sshCfg)
    if len(remoteArchivePath) > 0 {
        remoteOps.SetPath(remoteArchivePath)
    }
    
	return &MysqlDump{
		IP:        "localhost",
		Username:  username,
		Password:  password,
		Caller:    remoteExecuter,
		RemoteOps: remoteOps,
	}, err
}
//NewRemoteMysqlDump - will initialize a mysqldmp for remote execution
func NewRemoteMysqlDump(username, password string, sshCfg command.SshConfig) (*MysqlDump, error) {
	return NewRemoteMysqlDumpWithPath(username, password, sshCfg, "")
}

//Import - will import to mysql from the given reader
func (s *MysqlDump) Import(lfile io.Reader) (err error) {
	if err = s.RemoteOps.UploadFile(lfile); err == nil {
		err = s.restore()
	}
	lo.G.Debug("mysqldump Import called: ", err)
	
	if err != nil {
	    return
	}
	
	err = s.RemoteOps.RemoveRemoteFile()
	
	lo.G.Debug("mysqldump remove remote file called: ", err)
	return
}

//Dump - will dump a mysql to the given writer
func (s *MysqlDump) Dump(dest io.Writer) (err error) {
	err = s.Caller.Execute(dest, s.getDumpCommand())
	lo.G.Debug("mysqldump Dump called: ", s.getDumpCommand(), err)
	return
}

func (s *MysqlDump) restore() (err error) {
	callList := []string{
		s.getImportCommand(),
	}
	err = executeList(callList, s.Caller)
	lo.G.Debug("mysqldump restore called: ", callList, err)
	return
}

func (s *MysqlDump) getImportCommand() string {
	return fmt.Sprintf(MySQLDmpCreateCmd, s.getConnectCommand(MySQLDmpSQLBin), s.RemoteOps.Path())
}

func (s *MysqlDump) getFlushCommand() string {
	return fmt.Sprintf(MySQLDmpFlushCmd, s.getConnectCommand(MySQLDmpSQLBin))
}

func (s *MysqlDump) getDumpCommand() string {
	return fmt.Sprintf(MySQLDmpDumpCmd, s.getConnectCommand(MySQLDmpDumpBin))
}

func (s *MysqlDump) getConnectCommand(bin string) string {
	return fmt.Sprintf(MySQLDmpConnectCmd,
		bin,
		s.Username,
		s.IP,
		s.Password,
	)
}
