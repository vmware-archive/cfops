package database_test

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/pivotalservices/cfops/database"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Postgres", func() {
	var (
		context database.Context
	)

	BeforeEach(func() {
		context = database.NewContext("username", "password", "database", "host", 5000)
	})

	Describe("Create postgres command", func() {
		Context("From database context and writer", func() {
			It("Should create the command", func() {
				writer := bytes.NewBufferString("teststring")
				command := database.PostGresBackupCommand(context, writer)
				osCommand := command.(*exec.Cmd)
				Ω(len(osCommand.Args)).Should(Equal(8))
				Ω(osCommand.Args[0]).Should(Equal("pg_dump"))
				Ω(osCommand.Args[1]).Should(Equal("-h"))
				Ω(osCommand.Args[2]).Should(Equal("host"))
				Ω(osCommand.Args[3]).Should(Equal("-U"))
				Ω(osCommand.Args[4]).Should(Equal("username"))
				Ω(osCommand.Args[5]).Should(Equal("-p"))
				Ω(osCommand.Args[6]).Should(Equal("5000"))
				Ω(osCommand.Args[7]).Should(Equal("database"))
				Ω(os.Getenv("PGPASSWORD")).Should(Equal("password"))
				//Test writer assigned to the command
				resultWriter := osCommand.Stdout.(*bytes.Buffer)
				Ω(resultWriter.String()).Should(Equal("teststring"))
			})
		})
	})
})
