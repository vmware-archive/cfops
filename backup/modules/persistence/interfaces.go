package persistence

import "io"

type PersistanceBackup interface {
	Dump(io.Writer) error
}
