package persistence_test

import (
	"bytes"
	"fmt"
	"io"

	. "github.com/pivotalservices/cfbackup/modules/persistence"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	mysqlCatchCommand string
)

type MockSuccessCall struct {
}

func (s MockSuccessCall) Execute(destination io.Writer, command string) (err error) {
	mysqlCatchCommand = command
	return
}

type MockFailCall struct{}

func (s MockFailCall) Execute(destination io.Writer, command string) (err error) {
	err = fmt.Errorf("random mock error")
	return
}

var _ = Describe("Mysql", func() {
	var (
		mysqlDumpInstance *MysqlDump
		ip                string = "0.0.0.0"
		username          string = "testuser"
		password          string = "testpass"
		writer            bytes.Buffer
		successCall       *MockSuccessCall = &MockSuccessCall{}
	)

	Context("With command execute success", func() {
		BeforeEach(func() {
			mysqlDumpInstance = &MysqlDump{
				Ip:       ip,
				Username: username,
				Password: password,
				Caller:   successCall,
			}
			mysqlCatchCommand = ""
		})

		AfterEach(func() {
			mysqlDumpInstance = nil
		})

		It("Should return nil error", func() {
			err := mysqlDumpInstance.Dump(&writer)
			Ω(err).Should(BeNil())
		})
		It("Should execute mysqldump command", func() {
			mysqlDumpInstance.Dump(&writer)
			Ω(mysqlCatchCommand).Should(Equal("mysqldump -u testuser -h 0.0.0.0 --password=testpass --all-databases"))
		})
	})

	Context("With command execute failed", func() {
		BeforeEach(func() {
			mysqlDumpInstance = &MysqlDump{
				Ip:       ip,
				Username: username,
				Password: password,
				Caller:   &MockFailCall{},
			}
		})

		AfterEach(func() {
			mysqlDumpInstance = nil
		})

		It("Should return non nil error", func() {
			err := mysqlDumpInstance.Dump(&writer)
			Ω(err).ShouldNot(BeNil())
		})
	})

})
