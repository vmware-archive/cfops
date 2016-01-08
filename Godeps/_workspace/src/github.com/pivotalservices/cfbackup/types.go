package cfbackup

import "io"

//BackupContext - stores the base context information for a backup/restore
type BackupContext struct {
	TargetDir string
	StorageProvider
}

// StorageProvider is responsible for obtaining/managing a reader/writer to
// a storage type (eg disk/s3)
type StorageProvider interface {
	Reader(path ...string) (io.ReadCloser, error)
	Writer(path ...string) (io.WriteCloser, error)
}

// Tile is a deployable component that can be backed up
type Tile interface {
	Backup() error
	Restore() error
}

type connBucketInterface interface {
	Host() string
	AdminUser() string
	AdminPass() string
	OpsManagerUser() string
	OpsManagerPass() string
	Destination() string
}

type action func() error

type actionAdaptor func(t Tile) action
