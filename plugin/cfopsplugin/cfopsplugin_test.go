package cfopsplugin_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops/plugin/cfopsplugin"
)

var _ = Describe("cfopsplugin", func() {
	XDescribe("given a Call function", func() {
		Context("when calling on a named plugin with path", func() {
			It("", func() {
				Call("", "")
				Ω(true).Should(BeFalse())
			})
		})
	})

	XDescribe("given a Start function", func() {
		Context("when calling for meta data against a plugin", func() {

		})
	})
})
