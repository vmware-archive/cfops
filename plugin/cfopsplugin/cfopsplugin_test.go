package cfopsplugin_test

import (
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops/plugin/cfopsplugin"
)

var _ = Describe("cfopsplugin", func() {
	Describe("given a Call function", func() {
		Context("when calling on a named plugin with path", func() {
			It("should return a plugin which yield results from its methods", func() {
				pl, cl := Call("backuprestore", "../load/fixture_plugins/"+runtime.GOOS+"/sample")
				defer cl.Kill()
				Ω(pl.Backup()).Should(HaveOccurred())
			})

			It("should return a plugin which yield results from its methods", func() {
				pl, cl := Call("backuprestore", "../load/fixture_plugins/"+runtime.GOOS+"/sample")
				defer cl.Kill()
				Ω(pl.Restore()).ShouldNot(HaveOccurred())
			})
		})
	})
})
