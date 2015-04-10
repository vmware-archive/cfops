package osutils

import (
	"os"
	"path"
	"path/filepath"

	"github.com/pivotalservices/sftp"
)

type sftpClient interface {
	Create(path string) (*sftp.File, error)
	Mkdir(path string) error
	ReadDir(p string) ([]os.FileInfo, error)
}

// SafeCreateSSH creates a file, creating parent directories if needed on a remote machine via an ssh client
func SafeCreateSSH(client sftpClient, name ...string) (file *sftp.File, err error) {
	ssh := sshClientBucket{
		client: client,
	}
	fpath := path.Join(name...)
	basepath := filepath.Dir(fpath)

	if err = ssh.remoteSafeMkdirAll(basepath); err == nil {
		file, err = client.Create(fpath)
	}
	return
}

type sshClientBucket struct {
	client sftpClient
}

func (s sshClientBucket) remoteSafeMkdirAll(base string) (err error) {

	if !s.exists(base) {
		parentdir := filepath.Dir(base)

		if !s.exists(parentdir) {
			err = s.remoteSafeMkdirAll(filepath.Dir(parentdir))
		}

		if err == nil {
			err = s.client.Mkdir(base)
		}
	}
	return
}

func (s sshClientBucket) exists(fpath string) (ok bool) {
	ok = false

	if _, err := s.client.ReadDir(fpath); err == nil {
		ok = true
	}
	return
}

// SafeCreate creates a file, creating parent directories if needed
func SafeCreate(name ...string) (file *os.File, err error) {
	p, e := ensurePath(path.Join(name...))

	if e != nil {
		return nil, e
	}
	return os.Create(p)
}

func ensurePath(path string) (string, error) {
	base := filepath.Dir(path)
	e, _ := Exists(base)
	if e {
		return path, nil
	}

	// Create missing directory recursively
	err := os.MkdirAll(base, 0777)
	return path, err
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
