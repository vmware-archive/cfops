package persistence_test

import (
	"bytes"
	"fmt"
	"io"

	. "github.com/pivotalservices/cfops/backup/modules/persistence"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	successCounter    int
	failureCounter    int
	mysqlCatchCommand string
)

type MockSuccessCall struct {
}

func (s MockSuccessCall) Execute(destination io.Writer, command string) (err error) {
	mysqlCatchCommand = command
	successCounter++
	return
}

type MockFailCall struct{}

func (s MockFailCall) Execute(destination io.Writer, command string) (err error) {
	failureCounter++
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

	Context("Dump function call success", func() {
		BeforeEach(func() {
			successCounter = 0
			failureCounter = 0
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
			successCounter = 0
			failureCounter = 0
		})

		It("Should return nil error on success", func() {
			controlSuccessCount := 1
			controlFailureCount := 0
			err := mysqlDumpInstance.Dump(&writer)
			Ω(err).Should(BeNil())
			Ω(successCounter).Should(Equal(controlSuccessCount))
			Ω(failureCounter).Should(Equal(controlFailureCount))
			Ω(mysqlCatchCommand).Should(Equal("mysqldump -u testuser -h 0.0.0.0 --password=testpass --all-databases"))
		})
	})

	Context("Dump function call failure", func() {
		BeforeEach(func() {
			successCounter = 0
			failureCounter = 0
			mysqlDumpInstance = &MysqlDump{
				Ip:       ip,
				Username: username,
				Password: password,
				Caller:   &MockFailCall{},
			}
		})

		AfterEach(func() {
			mysqlDumpInstance = nil
			successCounter = 0
			failureCounter = 0
		})

		It("Should return non nil error on failure", func() {
			controlSuccessCount := 0
			controlFailureCount := 1
			err := mysqlDumpInstance.Dump(&writer)
			Ω(err).ShouldNot(BeNil())
			Ω(successCounter).Should(Equal(controlSuccessCount))
			Ω(failureCounter).Should(Equal(controlFailureCount))
		})
	})

})
