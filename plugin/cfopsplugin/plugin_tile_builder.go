package cfopsplugin

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfbackup/tiles/opsmanager"
	"github.com/pivotalservices/cfops/tileregistry"
	"github.com/xchapter7x/lo"
)

//DefaultCmdBuilder This is to build the default Cmd to execute the plugin
var DefaultCmdBuilder BuildCmd = func(filePath string, args string) *exec.Cmd {
	arguments := strings.Split(args, " ")
	arguments = append([]string{"plugin"}, arguments...)
	return exec.Command(filePath, arguments...)
}

//New - method to create a plugin tile
func (s *PluginTileBuilder) New(tileSpec tileregistry.TileSpec) (tile tileregistry.Tile, err error) {
	var opsManager *opsmanager.OpsManager
	var settingsReader io.Reader
	opsManager, err = opsmanager.NewOpsManager(tileSpec.OpsManagerHost, tileSpec.AdminUser, tileSpec.AdminPass, tileSpec.OpsManagerUser, tileSpec.OpsManagerPass, tileSpec.ArchiveDirectory, tileSpec.CryptKey)

	if settingsReader, err = opsManager.GetInstallationSettings(); err == nil {
		var brPlugin BackupRestorer
		installationSettings := cfbackup.NewConfigurationParserFromReader(settingsReader)
		pcf := NewPivotalCF(installationSettings.InstallationSettings, tileSpec)
		lo.G.Debug("", s.Meta.Name, s.FilePath, pcf)
		brPlugin, _ = s.call(tileSpec)
		brPlugin.Setup(pcf)
		tile = brPlugin
	}
	lo.G.Debug("error from getinstallationsettings: ", err)
	return
}

func (s *PluginTileBuilder) call(tileSpec tileregistry.TileSpec) (BackupRestorer, *plugin.Client) {
	RegisterPlugin(s.Meta.Name, new(BackupRestoreRPC))
	log.SetOutput(ioutil.Discard)

	client := plugin.NewClient(&plugin.ClientConfig{
		Stderr:          os.Stderr,
		SyncStdout:      os.Stdout,
		SyncStderr:      os.Stderr,
		HandshakeConfig: GetHandshake(),
		Plugins:         GetPlugins(),
		Cmd:             s.CmdBuilder(s.FilePath, tileSpec.PluginArgs),
	})
	rpcClient, err := client.Client()

	if err != nil {
		lo.G.Debug("rpcclient error: ", err)
		log.Fatal(err)
	}
	raw, err := rpcClient.Dispense(s.Meta.Name)

	if err != nil {
		lo.G.Debug("error: ", err)
		log.Fatal(err)
	}
	return raw.(BackupRestorer), client
}
