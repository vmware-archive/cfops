package plugin_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops/plugin"
	"github.com/pivotalservices/cfops/tileregistry"
)

var _ = Describe("given a plugin tile builder", func() {
	var pluginTileBuilder *PluginTileBuilder
	var pluginTile tileregistry.Tile
	var err error
	Context("when New is called", func() {
		BeforeEach(func() {
			pluginTileBuilder = new(PluginTileBuilder)
			pluginTileBuilder.FilePath = "test"
			pluginTileBuilder.Meta = Meta{Name: "fakemeta"}
			pluginTile, err = pluginTileBuilder.New(tileregistry.TileSpec{})
		})
		It("should return a plugin tile object", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(pluginTile).Should(BeAssignableToTypeOf(new(PluginTile)))
		})
	})
})
