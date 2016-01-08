package cfbackup

import (
	"io"
	"os"
	ospath "path"

	"github.com/pivotalservices/gtils/osutils"
)

// DiskProvider is a storage provider that stores your Docker images on local disk.
type DiskProvider struct {
	Directory string
}

// NewDiskProvider creates a new disk storage provider instance
func NewDiskProvider() StorageProvider {
	return &DiskProvider{}
}

// Reader returns an io.ReadCloser for the specified path
func (d *DiskProvider) Reader(path ...string) (io.ReadCloser, error) {
	filePath := ospath.Join(path...)
	return os.Open(filePath)
}

// Writer returns an io.WriteCloser for the specified path
func (d *DiskProvider) Writer(path ...string) (io.WriteCloser, error) {
	return osutils.SafeCreate(path...)
}
