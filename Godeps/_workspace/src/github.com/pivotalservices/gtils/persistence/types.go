package persistence

import "github.com/pivotalservices/gtils/command"

//PgDump - an object which can handle pgres dumps
type PgDump struct {
	sshCfg    command.SshConfig
	IP        string
	Port      int
	Database  string
	Username  string
	Password  string
	DbFile    string
	Caller    command.Executer
	RemoteOps remoteOperationsInterface
}

//MysqlDump - an object which can handle mysql dumps
type MysqlDump struct {
	IP         string
	Username   string
	Password   string
	DbFile     string
	ConfigFile string
	Caller     command.Executer
	RemoteOps  remoteOperationsInterface
}
