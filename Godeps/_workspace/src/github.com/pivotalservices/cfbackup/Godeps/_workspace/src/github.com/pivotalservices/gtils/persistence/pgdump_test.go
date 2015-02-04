package persistence_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/mock"
	"github.com/pivotalservices/gtils/osutils"
	. "github.com/pivotalservices/gtils/persistence"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	pgCatchCommand string
)

type pgMockSuccessCall struct{}

func (s pgMockSuccessCall) Execute(destination io.Writer, command string) (err error) {
	pgCatchCommand = command
	return
}

type pgMockFailFirstCall struct{}

func (s pgMockFailFirstCall) Execute(destination io.Writer, command string) (err error) {
	err = fmt.Errorf("random mock error")
	return
}

var _ = Describe("PgDump", func() {

	var (
		pgDumpInstance *PgDump
		ip             string = "0.0.0.0"
		username       string = "testuser"
		password       string = "testpass"
		writer         bytes.Buffer
	)
	Context("Import", func() {
		var (
			remoteFilePath string
			localFilePath  string
			dir            string
			sftpFailErr    error = errors.New("failed to make sftp connection")
		)

		BeforeEach(func() {
			dir, _ = ioutil.TempDir("", "spec")
			remoteFilePath = path.Join(dir, "rfile")
			localFilePath = path.Join(dir, "lfile")

			pgDumpInstance = &PgDump{
				Ip:       ip,
				Username: username,
				Password: password,
				Caller:   &pgMockSuccessCall{},
			}
		})

		AfterEach(func() {
			os.RemoveAll(dir)
		})

		Context("called w/ successful sftp connection", func() {
			BeforeEach(func() {
				pgDumpInstance.GetRemoteFile = func(command.SshConfig) (w io.WriteCloser, err error) {
					w, err = osutils.SafeCreate(remoteFilePath)
					return
				}
			})

			It("should copy local file to remote file and return nil error", func() {
				controlString := "hello there"
				l, _ := osutils.SafeCreate(localFilePath)
				l.WriteString(controlString)
				l.Close()
				l, _ = os.Open(localFilePath)
				err := pgDumpInstance.Import(l)
				l.Close()
				rf, _ := os.Open(remoteFilePath)
				defer rf.Close()
				rarray, _ := ioutil.ReadAll(rf)
				lf, _ := os.Open(localFilePath)
				defer lf.Close()
				larray, _ := ioutil.ReadAll(lf)
				Ω(err).Should(BeNil())
				Ω(rarray).Should(Equal(larray))
			})
		})

		Context("called w/ failed sftp connection", func() {
			BeforeEach(func() {
				pgDumpInstance.GetRemoteFile = func(command.SshConfig) (w io.WriteCloser, err error) {
					err = sftpFailErr
					return
				}
			})

			It("should return sftp connection error", func() {
				controlString := "hello there"
				l, _ := osutils.SafeCreate(localFilePath)
				l.WriteString(controlString)
				l.Close()
				l, _ = os.Open(localFilePath)
				err := pgDumpInstance.Import(l)
				l.Close()
				rf, _ := os.Open(remoteFilePath)
				defer rf.Close()
				rarray, _ := ioutil.ReadAll(rf)
				lf, _ := os.Open(localFilePath)
				defer lf.Close()
				larray, _ := ioutil.ReadAll(lf)

				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(Equal(sftpFailErr))
				Ω(rarray).ShouldNot(Equal(larray))
			})
		})

		Context("called w/ failed copy to remote", func() {
			BeforeEach(func() {
				pgDumpInstance.GetRemoteFile = func(command.SshConfig) (w io.WriteCloser, err error) {
					w = mock.NewReadWriteCloser(mock.READ_FAIL_ERROR, nil, nil)
					return
				}
			})

			It("should return failed copy error", func() {
				l := mock.NewReadWriteCloser(mock.READ_FAIL_ERROR, nil, nil)
				err := pgDumpInstance.Import(l)
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(Equal(mock.READ_FAIL_ERROR))
			})
		})

		Context("remote call w/ failed result from first call", func() {
			BeforeEach(func() {
				pgDumpInstance.Caller = &pgMockFailFirstCall{}
				pgDumpInstance.GetRemoteFile = func(command.SshConfig) (w io.WriteCloser, err error) {
					w = mock.NewReadWriteCloser(nil, nil, nil)
					return
				}
			})

			It("should return a call error", func() {
				l := mock.NewReadWriteCloser(nil, nil, nil)
				err := pgDumpInstance.Import(l)
				Ω(err).ShouldNot(BeNil())
			})
		})
	})

	Context("With caller successfully execute the command", func() {
		BeforeEach(func() {
			pgDumpInstance = &PgDump{
				Ip:       ip,
				Username: username,
				Password: password,
				Caller:   &pgMockSuccessCall{},
			}
			pgCatchCommand = ""
		})

		AfterEach(func() {
			pgDumpInstance = nil
		})

		It("Should execute the pg command", func() {
			pgDumpInstance.Dump(&writer)
			Ω(pgCatchCommand).Should(Equal("PGPASSWORD=testpass pg_dump -h 0.0.0.0 -U testuser -p 0 "))
		})

		It("Should return nil error", func() {
			err := pgDumpInstance.Dump(&writer)
			Ω(err).Should(BeNil())
		})
	})

	Context("With caller failed to execute command", func() {
		BeforeEach(func() {
			pgDumpInstance = &PgDump{
				Ip:       ip,
				Username: username,
				Password: password,
				Caller:   &pgMockFailFirstCall{},
			}
		})

		AfterEach(func() {
			pgDumpInstance = nil
		})

		It("Should return non nil error", func() {
			err := pgDumpInstance.Dump(&writer)
			Ω(err).ShouldNot(BeNil())
		})
	})
})
