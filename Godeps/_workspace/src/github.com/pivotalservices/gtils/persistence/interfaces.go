package persistence

import "io"

type remoteOperationsInterface interface {
	UploadFile(lfile io.Reader) (err error)
	Path() string
	RemoveRemoteFile() (err error)
}
