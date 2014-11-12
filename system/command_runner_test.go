package system_test

import (
	"github.com/malston/cf-logsearch-broker/system"
	"github.com/pivotal-golang/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("A command runner", func() {
	Context("is called with a valid command", func() {
		It("should run successfully", func() {
			commandRunner := &system.OSCommandRunner{
				Logger: lagertest.NewTestLogger("command-runner-test"),
			}
			err := commandRunner.Run("echo", "Hi", "there!")
			Ω(err).ToNot(HaveOccurred())
		})
	})
	Context("is called with an invalid command", func() {
		It("should fail with an error", func() {
			commandRunner := &system.OSCommandRunner{
				Logger: lagertest.NewTestLogger("command-runner-test"),
			}
			err := commandRunner.Run("bad", "command")
			Ω(err).To(HaveOccurred())
		})
	})
})
