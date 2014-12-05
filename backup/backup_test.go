package backup_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	. "github.com/pivotalservices/cfops/backup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backup", func() {
	var (
		dir     string
		context *BackupContext
	)

	BeforeEach(func() {
		dir, _ = ioutil.TempDir("", "cfops-backup")
		context = &BackupContext{
			Hostname:  "localhost",
			Username:  "admin",
			Password:  "admin",
			TPassword: "admin",
			Target:    path.Join(dir, "backup"),
		}
	})

	AfterEach(func() {
		os.RemoveAll(dir)
	})

	Describe("Something", func() {
		Context("With an empty directory", func() {
			It("should create the parent directory", func() {
				Ω(context.Target).NotTo(BeEquivalentTo(""))

			})
		})
	})

	Describe("ExecutePipeline method", func() {
		var (
			callCountCutoff int = 1
			twiceCounter    int = 0
			errorCounter    int = 0
			successCounter  int = 0
		)

		failFunction := func() (err error) {
			errorCounter++
			err = fmt.Errorf("random mock error")
			return
		}

		successFunction := func() (err error) {
			successCounter++
			return
		}

		partialFailFunction := func() (err error) {
			if twiceCounter >= callCountCutoff {
				errorCounter++
				err = fmt.Errorf("mock error")
			} else {
				successCounter++
			}
			twiceCounter++
			return
		}

		var MockSuccessPipeline []func() error = []func() error{
			successFunction,
			successFunction,
			successFunction,
		}

		var MockErrorPipeline []func() error = []func() error{
			failFunction,
			failFunction,
			failFunction,
		}

		var MockPartialErrorPipeline []func() error = []func() error{
			partialFailFunction,
			partialFailFunction,
			partialFailFunction,
		}

		BeforeEach(func() {
			twiceCounter = 0
			errorCounter = 0
			successCounter = 0
		})

		AfterEach(func() {
			twiceCounter = 0
			errorCounter = 0
			successCounter = 0
		})

		Context("With a successful pipeline to run", func() {
			It("should run all elements and return nil error", func() {
				controlCount := len(MockSuccessPipeline)
				err := context.ExecutePipeline(MockSuccessPipeline)
				Ω(err).To(BeNil())
				Expect(successCounter).To(Equal(controlCount))
			})
		})

		Context("With a failed pipeline to run", func() {
			It("should not run all elements and return not nil", func() {
				controlCount := len(MockErrorPipeline)
				err := context.ExecutePipeline(MockErrorPipeline)
				Ω(err).NotTo(BeNil())
				Expect(errorCounter).NotTo(Equal(controlCount))
			})
		})

		Context("With a partially failed pipeline to run", func() {
			It("should not run all elements and return not nil", func() {
				controlCallCount := callCountCutoff
				err := context.ExecutePipeline(MockPartialErrorPipeline)
				Ω(err).NotTo(BeNil())
				Expect(successCounter).To(Equal(controlCallCount))
			})
		})
	})
})
