package fake

import (
	"os"

	"github.com/pkg/sftp"
)

//MockSFTPClient - structure
type MockSFTPClient struct {
	FakeError error
}

func (s *MockSFTPClient) Create(path string) (file *sftp.File, err error) {
	err = s.FakeError
	return
}
func (s *MockSFTPClient) Mkdir(path string) (err error) {
	err = s.FakeError
	return
}
func (s *MockSFTPClient) ReadDir(p string) (fileInfo []os.FileInfo, err error) {
	err = s.FakeError
	return
}
func (s *MockSFTPClient) Remove(path string) (err error) {
	err = s.FakeError
	return
}
