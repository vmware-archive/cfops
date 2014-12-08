package persistence_test

import (
	"fmt"
	"os/exec"
	"strings"

	. "github.com/pivotalservices/cfops/backup/modules/persistence"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exec Command Adaptor", func() {
	Context("an adapted exec.Command function", func() {
		It("Should on success call through Command().Output() and return Output() methods response w/ nil error", func() {
			controlResponseString := "some random output"
			controlByteResponse := []byte(fmt.Sprintf("%s\n", controlResponseString))
			testcmd := fmt.Sprintf("echo %s", controlResponseString)
			syscall := ExecCommandOutputterAdaptor(exec.Command)
			commandArr := strings.Split(testcmd, " ")
			controlResponse, controlErr := exec.Command(commandArr[0], commandArr[1:]...).Output()
			res, err := syscall.Output(testcmd)
			Ω(err).Should(BeNil())
			Ω(controlErr).Should(BeNil())
			Ω(res).Should(Equal(controlByteResponse))
			Ω(res).Should(Equal(controlResponse))
		})

		It("Should, on success call w/ multiple line cmd string, call through Command().Output() and return Output() methods response w/ nil error", func() {
			controlResponseString := `echo "[mysqldump]
user=%s
password=%s"
`
			controlByteResponse := []byte(fmt.Sprintf("%s\n", controlResponseString))
			testcmd := fmt.Sprintf("echo %s", controlResponseString)
			syscall := ExecCommandOutputterAdaptor(exec.Command)
			commandArr := strings.Split(testcmd, " ")
			controlResponse, controlErr := exec.Command(commandArr[0], commandArr[1:]...).Output()
			res, err := syscall.Output(testcmd)
			Ω(err).Should(BeNil())
			Ω(controlErr).Should(BeNil())
			Ω(res).Should(Equal(controlByteResponse))
			Ω(res).Should(Equal(controlResponse))
		})

		It("Should on error call through Command().Output() and return Output() methods response w/ non nil error", func() {
			testcmd := "exit 1"
			syscall := ExecCommandOutputterAdaptor(exec.Command)
			commandArr := strings.Split(testcmd, " ")
			controlResponse, controlErr := exec.Command(commandArr[0], commandArr[1:]...).Output()
			res, err := syscall.Output(testcmd)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(Equal(controlErr))
			Ω(res).Should(Equal(controlResponse))
		})

	})
})
