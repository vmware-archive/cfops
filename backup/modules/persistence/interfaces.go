package persistence

type PersistanceBackup interface {
	Dump() error
}

type CmdOutputter interface {
	Output(cmd string) ([]byte, error)
}
