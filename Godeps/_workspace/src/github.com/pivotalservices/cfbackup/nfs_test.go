package cfbackup_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	. "github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/gtils/command"
	"github.com/pivotalservices/gtils/mock"
	"github.com/pivotalservices/gtils/osutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("nfs", func() {
	Describe("Import NFS Restore", func() {
		var (
			nfs           *NFSBackup
			buffer        *gbytes.Buffer
			controlString string = "test of local file"
			err           error
		)

		BeforeEach(func() {
			err = nil
		})

		AfterEach(func() {
			err = nil
		})

		Context("successful call to import", func() {

			BeforeEach(func() {
				lf := strings.NewReader(controlString)
				buffer = gbytes.NewBuffer()
				nfs = getNfs(buffer, &SuccessMockNFSExecuter{})
				err = nfs.Import(lf)
			})

			It("should return nil error", func() {
				Ω(err).Should(BeNil())
			})

			It("should write the local file contents to the remote", func() {
				Ω(buffer).Should(gbytes.Say(controlString))
			})
		})

		Context("error on command execution", func() {

			BeforeEach(func() {
				lf := strings.NewReader(controlString)
				buffer = gbytes.NewBuffer()
				nfs = getNfs(buffer, &FailureMockNFSExecuter{})
				err = nfs.Import(lf)
			})

			It("should return non-nil execution error", func() {
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(Equal(mockNfsCommandError))
			})

			It("should write the local file contents to the remote", func() {
				Ω(buffer).Should(gbytes.Say(controlString))
			})
		})

		Context("error on file upload", func() {
			BeforeEach(func() {
				buffer = gbytes.NewBuffer()
				nfs = getNfs(buffer, &SuccessMockNFSExecuter{})
			})

			Context("Read failure", func() {
				BeforeEach(func() {
					lf := mock.NewReadWriteCloser(mock.READ_FAIL_ERROR, nil, nil)
					err = nfs.Import(lf)
				})

				It("should return non-nil execution error", func() {
					Ω(err).ShouldNot(BeNil())
					Ω(err).Should(Equal(mock.READ_FAIL_ERROR))
				})

				It("should write the local file contents to the remote", func() {
					Ω(buffer).ShouldNot(gbytes.Say(controlString))
				})
			})

			Context("Writer related failure", func() {
				Context("Write failure", func() {
					BeforeEach(func() {
						lf := mock.NewReadWriteCloser(nil, mock.WRITE_FAIL_ERROR, nil)
						nfs = getNfs(lf, &SuccessMockNFSExecuter{})
						err = nfs.Import(lf)
					})

					It("should return non-nil execution error", func() {
						Ω(err).ShouldNot(BeNil())
						Ω(err).Should(Equal(mock.WRITE_FAIL_ERROR))
					})

					It("should write the local file contents to the remote", func() {
						Ω(buffer).ShouldNot(gbytes.Say(controlString))
					})
				})

				Context("Close failure", func() {
					BeforeEach(func() {
						lf := mock.NewReadWriteCloser(nil, nil, mock.CLOSE_FAIL_ERROR)
						nfs = getNfs(lf, &SuccessMockNFSExecuter{})
						err = nfs.Import(lf)
					})

					It("should return non-nil execution error", func() {
						Ω(err).ShouldNot(BeNil())
						Ω(err).Should(Equal(io.ErrShortWrite))
					})

					It("should write the local file contents to the remote", func() {
						Ω(buffer).ShouldNot(gbytes.Say(controlString))
					})
				})
			})
		})
	})

	Describe("BackupNfs", func() {
		var origExecuterFunction func(command.SshConfig) (command.Executer, error)
		var tmpfile *os.File
		var tmpfilepath string
		Context("called with valid arguments", func() {
			BeforeEach(func() {
				origExecuterFunction = NfsNewRemoteExecuter
				NfsNewRemoteExecuter = func(command.SshConfig) (ce command.Executer, err error) {
					ce = &SuccessMockNFSExecuter{}
					return
				}

				tmpdir, _ := ioutil.TempDir("/tmp", "spec")
				filename := "nfs.tar.gz"
				tmpfilepath = path.Join(tmpdir, filename)
				tmpfile, _ = osutils.SafeCreate(tmpfilepath)
			})

			AfterEach(func() {
				NfsNewRemoteExecuter = origExecuterFunction
				tmpfile.Close()
				os.Remove(tmpfilepath)
			})

			It("should return nil error and write success output to an outfile", func() {
				err := BackupNfs("pass", "1.2.3.4", tmpfile)
				b, _ := ioutil.ReadFile(tmpfilepath)
				Ω(b).Should(Equal([]byte(nfsSuccessString)))
				Ω(err).Should(BeNil())
			})
		})

		Context("called with invalid arguments", func() {
			BeforeEach(func() {
				origExecuterFunction = NfsNewRemoteExecuter
				NfsNewRemoteExecuter = func(command.SshConfig) (ce command.Executer, err error) {
					ce = &FailureMockNFSExecuter{}
					return
				}

				tmpdir, _ := ioutil.TempDir("/tmp", "spec")
				filename := "nfs.tar.gz"
				tmpfilepath = path.Join(tmpdir, filename)
				tmpfile, _ = osutils.SafeCreate(tmpfilepath)
			})

			AfterEach(func() {
				NfsNewRemoteExecuter = origExecuterFunction
				tmpfile.Close()
				os.Remove(tmpfilepath)
			})

			It("should return non nil error and write failure output to file", func() {
				err := BackupNfs("pass", "1.2.3.4", tmpfile)
				b, _ := ioutil.ReadFile(tmpfilepath)
				Ω(b).Should(Equal([]byte(nfsFailureString)))
				Ω(err).ShouldNot(BeNil())
			})
		})

	})

	Describe("NFSBackup", func() {
		var nfs *NFSBackup

		BeforeEach(func() {
			nfs = &NFSBackup{}
		})

		Context("sucessfully calling Dump", func() {
			BeforeEach(func() {
				nfs.Caller = &SuccessMockNFSExecuter{}
			})

			It("Should return nil error and a success message in the writer", func() {
				var b bytes.Buffer
				err := nfs.Dump(&b)
				Ω(err).Should(BeNil())
				Ω(b.String()).Should(Equal(nfsSuccessString))
			})
		})

		Context("failed calling Dump", func() {
			BeforeEach(func() {
				nfs.Caller = &FailureMockNFSExecuter{}
			})

			It("Should return non nil error and a failure output in the writer", func() {
				var b bytes.Buffer
				err := nfs.Dump(&b)
				Ω(err).ShouldNot(BeNil())
				Ω(b.String()).Should(Equal(nfsFailureString))
			})
		})

		Describe("NewNFSBackup", func() {
			Context("when executer is created successfully", func() {
				var origExecuterFunction func(command.SshConfig) (command.Executer, error)

				BeforeEach(func() {
					origExecuterFunction = NfsNewRemoteExecuter
					NfsNewRemoteExecuter = func(command.SshConfig) (command.Executer, error) {
						return &SuccessMockNFSExecuter{}, nil
					}
				})

				AfterEach(func() {
					NfsNewRemoteExecuter = origExecuterFunction
				})

				It("should return a nil error and a non-nil NFSBackup object", func() {
					n, err := NewNFSBackup("pass", "0.0.0.0")
					Ω(err).Should(BeNil())
					Ω(n).Should(BeAssignableToTypeOf(&NFSBackup{}))
					Ω(n).ShouldNot(BeNil())
				})
			})

			Context("when executer fails to be created properly", func() {
				var origExecuterFunction func(command.SshConfig) (command.Executer, error)

				BeforeEach(func() {
					origExecuterFunction = NfsNewRemoteExecuter
					NfsNewRemoteExecuter = func(command.SshConfig) (ce command.Executer, err error) {
						ce = &FailureMockNFSExecuter{}
						err = fmt.Errorf("we have an error")
						return
					}
				})

				AfterEach(func() {
					NfsNewRemoteExecuter = origExecuterFunction
				})

				It("should return a nil error and a NFSBackup object that is nil", func() {
					n, err := NewNFSBackup("pass", "0.0.0.0")
					Ω(err).ShouldNot(BeNil())
					Ω(n).Should(BeNil())
					Ω(n).Should(BeAssignableToTypeOf(&NFSBackup{}))
					Ω(n).Should(BeNil())
				})
			})
		})
	})
})
