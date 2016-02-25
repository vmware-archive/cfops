package cfopsplugin_test

import (
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	opsfakes "github.com/pivotalservices/cfbackup/tiles/opsmanager/fakes"
	. "github.com/pivotalservices/cfops/plugin/cfopsplugin"
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
			server = opsfakes.NewFakeOpsManagerServer(ghttp.NewTLSServer(), http.StatusInternalServerError, "{}", http.StatusOK, string(fileBytes[:]))
			pluginTileBuilder = new(PluginTileBuilder)
			pluginTileBuilder.FilePath = "../load/fixture_plugins/" + runtime.GOOS + "/sample"
			pluginTileBuilder.Meta = Meta{Name: "backuprestore"}
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
			Ω(pluginTile).Should(BeAssignableToTypeOf(new(BackupRestoreRPC)))
		})
	})
})
