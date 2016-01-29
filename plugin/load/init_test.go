package load_test

import (
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops/plugin/load"
	"github.com/pivotalservices/cfops/tileregistry"
)

var _ = Describe("given LoadPlugins", func() {
	Context("when called with a valid plugin directory containing plugins", func() {
		It("then it should register all plugins in the plugins directory", func() {
			controlTileLength := len(tileregistry.GetRegistry())
			err := LoadPlugins("fixture_plugins/" + runtime.GOOS)
			tileCount := len(tileregistry.GetRegistry())
			Ω(err).ShouldNot(HaveOccurred())
			Ω(tileCount).Should(BeNumerically(">", controlTileLength))
		})
	})
	Context("when called on a invalid or empty directory", func() {
		var err error
		BeforeEach(func() {
			err = LoadPlugins("dir-does-not-exist")
		})
		It("then it should yield an error", func() {
			Ω(err).Should(HaveOccurred())
		})
	})
	Context("when a plugin is not able to be registered", func() {
		var err error
		BeforeEach(func() {
			err = LoadPlugins("busted_plugins")
		})
		It("then it should yield an error", func() {
			Ω(err).Should(HaveOccurred())
		})
	})
})
