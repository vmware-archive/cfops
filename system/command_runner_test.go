package system_test

import (
	"github.com/cloudfoundry/gosteno"
	"github.com/pivotalservices/cfops/system"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("A command runner", func() {
	Context("is called with a valid command", func() {
		It("should run successfully", func() {
			commandRunner := &system.OSCommandRunner{
				gosteno.NewLogger("TestLogger"),
			}
			err := commandRunner.Run("echo", "Hi", "there!")
			Ω(err).ToNot(HaveOccurred())
		})
	})
	Context("is called with an invalid command", func() {
		It("should fail with an error", func() {
			commandRunner := &system.OSCommandRunner{
				gosteno.NewLogger("TestLogger"),
			}
			err := commandRunner.Run("bad", "command")
			Ω(err).To(HaveOccurred())
		})
	})
})
