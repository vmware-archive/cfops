package command_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	. "github.com/pivotalservices/cfops/command"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exec Command Adaptor", func() {
	Context("an adapted exec.Command function", func() {
		It("Should on success call through Command().Output() and return Output() methods response w/ nil error", func() {
			controlResponseString := "some random output"
			testcmd := fmt.Sprintf("echo %s", controlResponseString)
			syscall := NewLocalExecuter()
			commandArr := strings.Split(testcmd, " ")
			controlResponse, controlErr := exec.Command(commandArr[0], commandArr[1:]...).Output()
			var b bytes.Buffer
			err := syscall.Execute(&b, testcmd)
			Ω(err).Should(BeNil())
			Ω(controlErr).Should(BeNil())
			Ω(b.String()).Should(Equal(string(controlResponse[:])))
		})

		It("Should, on success call w/ multiple line cmd string, call through Command().Output() and return Output() methods response w/ nil error", func() {
			controlResponseString := `echo "[mysqldump]
user=%s
password=%s"
`
			testcmd := fmt.Sprintf("echo %s", controlResponseString)
			syscall := NewLocalExecuter()
			commandArr := strings.Split(testcmd, " ")
			controlResponse, controlErr := exec.Command(commandArr[0], commandArr[1:]...).Output()
			var b bytes.Buffer
			err := syscall.Execute(&b, testcmd)
			Ω(err).Should(BeNil())
			Ω(controlErr).Should(BeNil())
			Ω(b.String()).Should(Equal(string(controlResponse[:])))
		})

		It("Should on error call through Command().Output() and return Output() methods response w/ non nil error", func() {
			testcmd := "exit 1"
			syscall := NewLocalExecuter()
			commandArr := strings.Split(testcmd, " ")
			controlResponse, controlErr := exec.Command(commandArr[0], commandArr[1:]...).Output()
			var b bytes.Buffer
			err := syscall.Execute(&b, testcmd)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(controlErr))
			Ω(b.String()).Should(Equal(string(controlResponse[:])))
		})

	})
})
