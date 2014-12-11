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
	pgsuccessCounter int
	pgfailureCounter int
	pgCatchCommand   string
)

type pgMockSuccessCall struct{}

func (s pgMockSuccessCall) Execute(destination io.Writer, command string) (err error) {
	pgCatchCommand = command
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
		writer         bytes.Buffer
	)

	Context("Dump function call success", func() {
		BeforeEach(func() {
			pgsuccessCounter = 0
			pgfailureCounter = 0
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
			pgsuccessCounter = 0
			pgfailureCounter = 0
		})

		It("Should return nil error on success", func() {
			controlSuccessCount := 1
			controlFailureCount := 0
			err := pgDumpInstance.Dump(&writer)
			Ω(err).Should(BeNil())
			Ω(pgsuccessCounter).Should(Equal(controlSuccessCount))
			Ω(pgfailureCounter).Should(Equal(controlFailureCount))
			Ω(pgCatchCommand).Should(Equal("PGPASSWORD=testpass pg_dump -h 0.0.0.0 -U testuser -p 0 "))
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
				Caller:   &pgMockFailFirstCall{},
			}
		})

		AfterEach(func() {
			pgDumpInstance = nil
			pgsuccessCounter = 0
			pgfailureCounter = 0
		})

		It("Should return non nil error on failure", func() {
			controlSuccessCount := 0
			controlFailureCount := 1
			err := pgDumpInstance.Dump(&writer)
			Ω(err).ShouldNot(BeNil())
			Ω(pgsuccessCounter).Should(Equal(controlSuccessCount))
			Ω(pgfailureCounter).Should(Equal(controlFailureCount))
		})
	})
})
