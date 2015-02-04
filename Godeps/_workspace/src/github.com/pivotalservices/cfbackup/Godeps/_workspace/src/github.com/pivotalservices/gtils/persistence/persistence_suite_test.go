package persistence_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/gtils/mock"

	"testing"
)

func TestPersistance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TestPersistance Suite")
}

var (
	uploadError error = errors.New("file upload failed")
)

type mockRemoteOps struct {
	Err    error
	Writer io.Writer
}

func (s *mockRemoteOps) UploadFile(lfile io.Reader) error {

	if s.Writer == nil {
		s.Writer = mock.NewReadWriteCloser(nil, nil, nil)
	}

	if s.Err == nil {
		io.Copy(s.Writer, lfile)
	}
	return s.Err
}

type MockSuccessCall struct{}

func (s MockSuccessCall) Execute(destination io.Writer, command string) (err error) {
	io.Copy(destination, bytes.NewReader([]byte(command)))
	return
}

type MockFailCall struct {
	CatchCommand string
}

func (s MockFailCall) Execute(destination io.Writer, command string) (err error) {
	err = fmt.Errorf("random mock error")
	return
}
