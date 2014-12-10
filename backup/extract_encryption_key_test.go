package backup_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	. "github.com/pivotalservices/cfops/backup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	eeksuccessControlOuput string = "successful execute"
	eekfailureControlOuput string = "failed to execute"
)

type eekMockSuccessCall struct{}

func (s eekMockSuccessCall) Execute(destination io.Writer, command string) (err error) {
	destination.Write([]byte(eeksuccessControlOuput))
	return
}

type eekMockFailCall struct{}

func (s eekMockFailCall) Execute(destination io.Writer, command string) (err error) {
	destination.Write([]byte(eekfailureControlOuput))
	err = fmt.Errorf("random mock error")
	return
}

var _ = Describe("ToggleCCJobRunner", func() {
	var (
		backupDir                string = "backupDirTest"
		deploymentDir            string = "deploymentDirTest"
		controlBackupFilename    string = path.Join(backupDir, "cc_db_encryption_key.txt")
		controlSuccessByteOutput []byte = []byte(eeksuccessControlOuput)
		controlFailureByteOutput []byte = []byte(eekfailureControlOuput)
	)

	BeforeEach(func() {
		os.Remove(controlBackupFilename)
	})

	AfterEach(func() {
		os.Remove(controlBackupFilename)
	})

	Context("successful call", func() {
		It("Should return nil error and pass through the cmd output", func() {
			err := ExtractEncryptionKey(backupDir, deploymentDir, &eekMockSuccessCall{})
			fileBytes, _ := ioutil.ReadFile(controlBackupFilename)
			Ω(err).Should(BeNil())
			Ω(fileBytes).Should(Equal(controlSuccessByteOutput))
		})
	})

	Context("failure call", func() {
		It("Should return non nil error and pass through the cmd output", func() {
			err := ExtractEncryptionKey(backupDir, deploymentDir, &eekMockFailCall{})
			fileBytes, _ := ioutil.ReadFile(controlBackupFilename)
			Ω(err).ShouldNot(BeNil())
			Ω(fileBytes).Should(Equal(controlFailureByteOutput))
		})
	})
})
