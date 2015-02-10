package command_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	. "github.com/pivotalservices/gtils/command"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exec Command Adaptor", func() {
	Context("an adapted exec.Command function", func() {
		var filename string = "output.txt"

		BeforeEach(func() {
			os.Remove(filename)
		})

		AfterEach(func() {
			os.Remove(filename)
		})

		It("Should write to file when the io.Writer is a file reference", func() {
			controlString := "exit 1"
			testControlByte := []byte(fmt.Sprintf("%s\n", controlString))
			testcmd := fmt.Sprintf("echo %s", controlString)
			syscall := NewLocalExecuter()
			commandArr := strings.Split(testcmd, " ")
			exec.Command(commandArr[0], commandArr[1:]...).Output()
			b, err := os.Create(filename)
			defer b.Close()
			err = syscall.Execute(b, testcmd)
			fileBytes, _ := ioutil.ReadFile(filename)
			Ω(err).Should(BeNil())
			Ω(fileBytes).Should(Equal(testControlByte))
		})

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
