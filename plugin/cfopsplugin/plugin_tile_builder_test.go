package cfopsplugin_test

import (
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotalservices/cfbackup/tileregistry"
	opsfakes "github.com/pivotalservices/cfbackup/tiles/opsmanager/fakes"
	. "github.com/pivotalservices/cfops/plugin/cfopsplugin"
)

var _ = Describe("given a plugin tile builder", func() {
	testNewPluginTileBuilder("when new is called on a default PCF installation", "./fixtures/installation-settings-1-6-default.json")
	testNewPluginTileBuilder("when new is called on a AWS PCF installation", "./fixtures/installation-settings-1-6-aws.json")
	testNewPluginTileBuilder("when new is called on a PCF installation w/ RabbitMQ", "./fixtures/installation-settings-with-rabbit.json")
})

func testNewPluginTileBuilder(behavior string, fixtureInstallationSettingsPath string) {
	var pluginTileBuilder *PluginTileBuilder
	var pluginTile tileregistry.TileCloser
	var err error
	Context(behavior, func() {

		var (
			server   *ghttp.Server
			fakeUser = "fakeuser"
			fakePass = "fakepass"
		)

		BeforeEach(func() {
			fileBytes, _ := ioutil.ReadFile(fixtureInstallationSettingsPath)
			server = opsfakes.NewFakeOpsManagerServer(ghttp.NewTLSServer(), http.StatusInternalServerError, "{}", http.StatusOK, string(fileBytes[:]))
			pluginTileBuilder = new(PluginTileBuilder)
			pluginTileBuilder.FilePath = "../load/fixture_plugins/" + runtime.GOOS + "/sample"
			pluginTileBuilder.Meta = Meta{Name: "backuprestore"}
			pluginTileBuilder.CmdBuilder = DefaultCmdBuilder
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
			pluginTile.Close()
		})
		It("should return a plugin tile object", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should yield results from its methods", func() {
			Ω(pluginTile.Backup()).Should(HaveOccurred())
		})
		It("should yield results from its methods", func() {
			Ω(pluginTile.Restore()).ShouldNot(HaveOccurred())
		})
	})
}

var _ = Describe("given the default command builder", func() {
	Context("when a string with argument passed in", func() {
		var arguments = "--arg1 value1 --arg2 value2"
		var filePath = "path"
		It("should return a command with space delimited args", func() {
			cmd := DefaultCmdBuilder(filePath, arguments)
			Ω(cmd.Path).Should(Equal(filePath))
			Ω(cmd.Args[1]).Should(Equal("plugin"))
			Ω(cmd.Args[2]).Should(Equal("--arg1"))
			Ω(cmd.Args[3]).Should(Equal("value1"))
			Ω(cmd.Args[4]).Should(Equal("--arg2"))
			Ω(cmd.Args[5]).Should(Equal("value2"))
		})
	})
})
