package backup_test

import (
	"fmt"
	"io"

	. "github.com/pivotalservices/cfops/backup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	successControlOuput string = "successful execute"
	failureControlOuput string = "failed to execute"
)

type MockSuccessCall struct{}

func (s MockSuccessCall) Execute(destination io.Writer, command string) (err error) {
	destination.Write([]byte(successControlOuput))
	return
}

type MockFailCall struct{}

func (s MockFailCall) Execute(destination io.Writer, command string) (err error) {
	destination.Write([]byte(failureControlOuput))
	err = fmt.Errorf("random mock error")
	return
}

var _ = Describe("ToggleCCJobRunner", func() {
	Context("successful call", func() {
		var (
			username  string = "usertest"
			password  string = "passwrdtest"
			serverUrl string = "someurl.com"
		)
		It("Should return nil error and pass through the cmd output", func() {
			msg, err := ToggleCCJobRunner(username, password, serverUrl, &MockSuccessCall{})
			Ω(err).Should(BeNil())
			Ω(msg).Should(Equal(successControlOuput))
		})
	})

	Context("failure call", func() {
		var (
			username  string = "usertest"
			password  string = "passwrdtest"
			serverUrl string = "someurl.com"
		)
		It("Should return non nil error and pass through the cmd output", func() {
			msg, err := ToggleCCJobRunner(username, password, serverUrl, &MockFailCall{})
			Ω(err).ShouldNot(BeNil())
			Ω(msg).Should(Equal(failureControlOuput))
		})
	})
})
