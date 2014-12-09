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
	pgsuccessCounter int
	pgfailureCounter int
)

type pgMockSuccessCall struct{}

func (s pgMockSuccessCall) Execute(destination io.Writer, command string) (err error) {
	pgsuccessCounter++
	return
}

type pgMockFailFirstCall struct{}

func (s pgMockFailFirstCall) Execute(destination io.Writer, command string) (err error) {
	pgfailureCounter++
	err = fmt.Errorf("random mock error")
	return
}

var _ = Describe("Mysql", func() {

	var (
		pgDumpInstance *PgDump
		ip             string = "0.0.0.0"
		username       string = "testuser"
		password       string = "testpass"
		dbFile         string = "testfile"
	)

	Context("Dump function call success", func() {
		BeforeEach(func() {
			pgsuccessCounter = 0
			pgfailureCounter = 0
			pgDumpInstance = &PgDump{
				Ip:       ip,
				Username: username,
				Password: password,
				DbFile:   dbFile,
				Caller:   &pgMockSuccessCall{},
			}
		})

		AfterEach(func() {
			pgDumpInstance = nil
			pgsuccessCounter = 0
			pgfailureCounter = 0
			os.Remove(dbFile)
		})

		It("Should return nil error on success", func() {
			controlSuccessCount := 1
			controlFailureCount := 0
			b, _ := osutils.SafeCreate(dbFile)
			defer b.Close()
			err := pgDumpInstance.Dump(b)
			Ω(err).Should(BeNil())
			Ω(pgsuccessCounter).Should(Equal(controlSuccessCount))
			Ω(pgfailureCounter).Should(Equal(controlFailureCount))
		})
	})

	Context("Dump function call failure", func() {
		BeforeEach(func() {
			pgsuccessCounter = 0
			pgfailureCounter = 0
			pgDumpInstance = &PgDump{
				Ip:       ip,
				Username: username,
				Password: password,
				DbFile:   dbFile,
				Caller:   &pgMockFailFirstCall{},
			}
		})

		AfterEach(func() {
			pgDumpInstance = nil
			pgsuccessCounter = 0
			pgfailureCounter = 0
			os.Remove(dbFile)
		})

		It("Should return non nil error on failure", func() {
			controlSuccessCount := 0
			controlFailureCount := 1
			b, _ := osutils.SafeCreate(dbFile)
			defer b.Close()
			err := pgDumpInstance.Dump(b)
			Ω(err).ShouldNot(BeNil())
			Ω(pgsuccessCounter).Should(Equal(controlSuccessCount))
			Ω(pgfailureCounter).Should(Equal(controlFailureCount))
		})
	})
})
