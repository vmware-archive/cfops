package persistence

type PersistanceBackup interface {
	Dump() error
}
