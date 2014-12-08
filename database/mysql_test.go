package database_test

import (
	"bytes"
	"os/exec"

	"github.com/pivotalservices/cfops/database"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mysql", func() {
	var (
		context database.Context
	)

	BeforeEach(func() {
		context = database.NewContext("username", "password", "--all-databases", "host", 5000)
	})

	Describe("Create mysql command", func() {
		Context("From database context and writer", func() {
			It("Should create the command for all the databases", func() {
				writer := bytes.NewBufferString("teststring")
				command := database.MysqlBackupCommand(context, writer)
				osCommand := command.(*exec.Cmd)
				Ω(len(osCommand.Args)).Should(Equal(10))
				Ω(osCommand.Args[0]).Should(Equal("mysqldump"))
				Ω(osCommand.Args[1]).Should(Equal("-h"))
				Ω(osCommand.Args[2]).Should(Equal("host"))
				Ω(osCommand.Args[3]).Should(Equal("-U"))
				Ω(osCommand.Args[4]).Should(Equal("username"))
				Ω(osCommand.Args[5]).Should(Equal("-p"))
				Ω(osCommand.Args[6]).Should(Equal("password"))
				Ω(osCommand.Args[7]).Should(Equal("-P"))
				Ω(osCommand.Args[8]).Should(Equal("5000"))
				Ω(osCommand.Args[9]).Should(Equal("--all-databases"))
				resultWriter := osCommand.Stdout.(*bytes.Buffer)
				Ω(resultWriter.String()).Should(Equal("teststring"))
			})
			It("Should create the command for the database", func() {
				context = database.NewContext("username", "password", "databaseName", "host", 5000)
				writer := bytes.NewBufferString("teststring")
				command := database.MysqlBackupCommand(context, writer)
				osCommand := command.(*exec.Cmd)
				Ω(len(osCommand.Args)).Should(Equal(10))
				Ω(osCommand.Args[9]).Should(Equal("--databases databaseName"))
				resultWriter := osCommand.Stdout.(*bytes.Buffer)
				Ω(resultWriter.String()).Should(Equal("teststring"))
			})
		})

	})
})
