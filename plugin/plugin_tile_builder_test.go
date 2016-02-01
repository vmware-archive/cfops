package plugin_test

import (
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/pivotalservices/cfops/plugin"
	"github.com/pivotalservices/cfops/tileregistry"
)

var _ = Describe("given a plugin tile builder", func() {
	var pluginTileBuilder *PluginTileBuilder
	var pluginTile tileregistry.Tile
	var err error
	Context("when New is called", func() {

		var (
			server   *ghttp.Server
			fakeUser = "fakeuser"
			fakePass = "fakepass"
		)

		BeforeEach(func() {
			fileBytes, _ := ioutil.ReadFile("./fixtures/installation-settings-1-6-default.json")
			server = ghttp.NewTLSServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyBasicAuth(fakeUser, fakePass),
					ghttp.RespondWith(http.StatusOK, string(fileBytes[:])),
				),
			)

			pluginTileBuilder = new(PluginTileBuilder)
			pluginTileBuilder.FilePath = "test"
			pluginTileBuilder.Meta = Meta{Name: "fakemeta"}
			pluginTile, err = pluginTileBuilder.New(tileregistry.TileSpec{
				OpsManagerHost:   strings.Replace(server.URL(), "https://", "", 1),
				AdminUser:        fakeUser,
				AdminPass:        fakePass,
				OpsManagerUser:   "ubuntu",
				OpsManagerPass:   "xxx",
				ArchiveDirectory: "/tmp",
			})
		})

		AfterEach(func() {
			server.Close()
		})
		It("should return a plugin tile object", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(pluginTile).Should(BeAssignableToTypeOf(new(PluginTile)))
		})
	})
})
