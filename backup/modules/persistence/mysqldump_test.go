package persistence_test

import (
	"fmt"
	"io"
	"os"

	. "github.com/pivotalservices/cfops/backup/modules/persistence"
	"github.com/pivotalservices/cfops/osutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	successCounter int
	failureCounter int
)

type MockSuccessCall struct{}

func (s MockSuccessCall) Execute(destination io.Writer, command string) (err error) {
	successCounter++
	return
}

type MockFailFirstCall struct{}

func (s MockFailFirstCall) Execute(destination io.Writer, command string) (err error) {
	failureCounter++
	err = fmt.Errorf("random mock error")
	return
}

type MockFailSecondCall struct{}

func (s MockFailSecondCall) Execute(destination io.Writer, command string) (err error) {
	if successCounter < 1 {
		successCounter++
	} else {
		failureCounter++
		err = fmt.Errorf("mock failure")
	}
	return
}

var _ = Describe("Mysql", func() {
	var (
		mysqlDumpInstance *MysqlDump
		ip                string = "0.0.0.0"
		username          string = "testuser"
		password          string = "testpass"
		dbFile            string = "testfile"
		configFile        string = ".my.cnf"
	)

	Context("Dump function call success", func() {
		BeforeEach(func() {
			successCounter = 0
			failureCounter = 0
			mysqlDumpInstance = &MysqlDump{
				Ip:         ip,
				Username:   username,
				Password:   password,
				DbFile:     dbFile,
				ConfigFile: configFile,
				Caller:     &MockSuccessCall{},
			}
		})

		AfterEach(func() {
			mysqlDumpInstance = nil
			successCounter = 0
			failureCounter = 0
			os.Remove(configFile)
			os.Remove(dbFile)
		})

		It("Should return nil error on success", func() {
			controlSuccessCount := 2
			controlFailureCount := 0
			f, _ := osutils.SafeCreate(dbFile)
			defer f.Close()
			err := mysqlDumpInstance.Dump(f)
			Ω(err).Should(BeNil())
			Ω(successCounter).Should(Equal(controlSuccessCount))
			Ω(failureCounter).Should(Equal(controlFailureCount))
		})
	})

	Context("Dump function call failure", func() {
		BeforeEach(func() {
			successCounter = 0
			failureCounter = 0
			mysqlDumpInstance = &MysqlDump{
				Ip:         ip,
				Username:   username,
				Password:   password,
				DbFile:     dbFile,
				ConfigFile: configFile,
				Caller:     &MockFailFirstCall{},
			}
		})

		AfterEach(func() {
			mysqlDumpInstance = nil
			successCounter = 0
			failureCounter = 0
			os.Remove(configFile)
			os.Remove(dbFile)
		})

		It("Should return non nil error on failure", func() {
			controlSuccessCount := 0
			controlFailureCount := 1
			f, _ := osutils.SafeCreate(dbFile)
			defer f.Close()
			err := mysqlDumpInstance.Dump(f)
			Ω(err).ShouldNot(BeNil())
			Ω(successCounter).Should(Equal(controlSuccessCount))
			Ω(failureCounter).Should(Equal(controlFailureCount))
		})
	})

	Context("Dump function call partial failure", func() {
		BeforeEach(func() {
			successCounter = 0
			failureCounter = 0
			mysqlDumpInstance = &MysqlDump{
				Ip:         ip,
				Username:   username,
				Password:   password,
				DbFile:     dbFile,
				ConfigFile: configFile,
				Caller:     &MockFailSecondCall{},
			}
		})

		AfterEach(func() {
			mysqlDumpInstance = nil
			successCounter = 0
			failureCounter = 0
			os.Remove(configFile)
			os.Remove(dbFile)
		})

		It("Should return non nil error on failure", func() {
			controlSuccessCount := 1
			controlFailureCount := 1
			f, _ := osutils.SafeCreate(dbFile)
			defer f.Close()
			err := mysqlDumpInstance.Dump(f)
			Ω(err).ShouldNot(BeNil())
			Ω(successCounter).Should(Equal(controlSuccessCount))
			Ω(failureCounter).Should(Equal(controlFailureCount))
		})
	})

})
