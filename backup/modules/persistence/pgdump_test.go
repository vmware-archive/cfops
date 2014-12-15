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

var _ = Describe("Mysql", func() {

	var (
		pgDumpInstance *PgDump
		ip             string = "0.0.0.0"
		username       string = "testuser"
		password       string = "testpass"
		writer         bytes.Buffer
	)

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
