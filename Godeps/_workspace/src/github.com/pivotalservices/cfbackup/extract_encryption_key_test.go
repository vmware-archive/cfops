package cfbackup_test

import (
	"bytes"

	. "github.com/pivotalservices/cfbackup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExtractEncryptionKey", func() {

	Context("successful call", func() {
		var (
			deploymentDir string = "fixtures/encryptionkey"
		)
		It("Should return nil error and write the correct key", func() {
			var keystring bytes.Buffer
			err := ExtractEncryptionKey(&keystring, deploymentDir)
			Ω(err).Should(BeNil())
			Ω(keystring.String()).Should(Equal("a5f5bc93ea6221499492"))
		})
	})

	Context("yaml dir doesnt exist", func() {
		var deploymentDir string = "dirdoesntexist"

		It("Should return non nil error and an empty writer", func() {
			var keystring bytes.Buffer
			err := ExtractEncryptionKey(&keystring, deploymentDir)
			Ω(err).ShouldNot(BeNil())
			Ω(keystring.String()).Should(BeEmpty())
		})
	})

	Context("invalid yaml file", func() {
		var deploymentDir string = "fixtures/encryptionkey/invalidfileerror"

		It("Should return non nil error and an empty writer", func() {
			var keystring bytes.Buffer
			err := ExtractEncryptionKey(&keystring, deploymentDir)
			Ω(err).ShouldNot(BeNil())
			Ω(keystring.String()).Should(BeEmpty())
		})
	})

	Context("yaml dir doesnt exist", func() {
		var deploymentDir string = "fixtures/encryptionkey/emptyerror"

		It("Should return non nil error and an empty writer", func() {
			var keystring bytes.Buffer
			err := ExtractEncryptionKey(&keystring, deploymentDir)
			Ω(err).ShouldNot(BeNil())
			Ω(keystring.String()).Should(BeEmpty())
		})
	})
})
