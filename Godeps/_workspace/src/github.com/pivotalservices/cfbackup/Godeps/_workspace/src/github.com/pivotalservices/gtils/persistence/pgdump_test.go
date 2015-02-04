package persistence_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pivotalservices/gtils/mock"
	"github.com/pivotalservices/gtils/osutils"
	. "github.com/pivotalservices/gtils/persistence"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	pgCatchCommand string
)

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
				Caller:   &MockSuccessCall{},
			}
		})

		AfterEach(func() {
			os.RemoveAll(dir)
		})

		Context("called w/ successful sftp connection", func() {
			var output bytes.Buffer
			BeforeEach(func() {
				pgDumpInstance.RemoteOps = &mockRemoteOps{
					Writer: &output,
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
				lf, _ := os.Open(localFilePath)
				defer lf.Close()
				larray, _ := ioutil.ReadAll(lf)
				Ω(err).Should(BeNil())
				Ω(output.String()).Should(Equal(string(larray[:])))
			})
		})

		Context("called w/ failed sftp connection", func() {
			var output bytes.Buffer
			BeforeEach(func() {
				pgDumpInstance.RemoteOps = &mockRemoteOps{
					Err:    sftpFailErr,
					Writer: &output,
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
				lf, _ := os.Open(localFilePath)
				defer lf.Close()
				larray, _ := ioutil.ReadAll(lf)

				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(Equal(sftpFailErr))
				Ω(output.String()).ShouldNot(Equal(string(larray[:])))
			})
		})

		Context("called w/ failed copy to remote", func() {
			BeforeEach(func() {
				pgDumpInstance.RemoteOps = &mockRemoteOps{
					Err: mock.READ_FAIL_ERROR,
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
				pgDumpInstance.Caller = &MockFailCall{}
				pgDumpInstance.RemoteOps = &mockRemoteOps{}
			})

			It("should return a call error", func() {
				l := mock.NewReadWriteCloser(nil, nil, nil)
				err := pgDumpInstance.Import(l)
				Ω(err).ShouldNot(BeNil())
			})
		})
	})

	Context("Dump", func() {
		Context("With caller successfully execute the command", func() {
			BeforeEach(func() {
				pgDumpInstance = &PgDump{
					Ip:       ip,
					Username: username,
					Password: password,
					Caller:   &MockSuccessCall{},
				}
				pgCatchCommand = ""
			})

			AfterEach(func() {
				pgDumpInstance = nil
			})

			It("Should execute the pg command", func() {
				var b bytes.Buffer
				pgDumpInstance.Dump(&b)
				cmd := fmt.Sprintf("PGPASSWORD=%s %s -h %s -U %s -p 0 ", password, PGDMP_DUMP_BIN, ip, username)
				Ω(b.String()).Should(Equal(cmd))
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
					Caller:   &MockFailCall{},
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
})
