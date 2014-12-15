package backup_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	. "github.com/pivotalservices/cfops/backup"
	"github.com/pivotalservices/cfops/command"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	nfsSuccessString string = "success nfs"
	nfsFailureString string = "failed nfs"
)

type SuccessMockNFSExecuter struct{}

func (s *SuccessMockNFSExecuter) Execute(dest io.Writer, cmd string) (err error) {
	io.Copy(dest, strings.NewReader(nfsSuccessString))
	return
}

type FailureMockNFSExecuter struct{}

func (s *FailureMockNFSExecuter) Execute(dest io.Writer, cmd string) (err error) {
	io.Copy(dest, strings.NewReader(nfsFailureString))
	err = fmt.Errorf("error occurred")
	return
}

var _ = Describe("NFSBackup", func() {
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
