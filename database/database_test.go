package database_test

import (
	"bytes"
	"errors"
	"github.com/pivotalservices/cfops/database"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeCommand struct {
	success bool
}

func (command *fakeCommand) Run() error {
	if !command.success {
		return errors.New("")
	}
	return nil
}

var _ = Describe("Database", func() {
	var (
		context database.Context
	)
	BeforeEach(func() {
		context = database.NewContext("", "", "", "", 5000)
	})

	Describe("Execute the backup command", func() {
		Context("Error from command", func() {
			It("should throw an error", func() {
				var writer bytes.Buffer
				var failedCommand = func(context database.Context, writer io.Writer) database.CommandInterface {
					return &fakeCommand{
						success: false,
					}
				}
				err := context.Exec(&writer, failedCommand)
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("Got success from command", func() {
			It("should dump to the writer", func() {
				var writer bytes.Buffer
				var successCommand = func(context database.Context, writer io.Writer) database.CommandInterface {
					return &fakeCommand{
						success: true,
					}
				}
				err := context.Exec(&writer, successCommand)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})

})
