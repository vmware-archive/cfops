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
	"github.com/pivotalservices/cfbackup/tileregistry"
	"github.com/pivotalservices/cfbackup/tiles/opsmanager"
	"github.com/xchapter7x/lo"
)

//DefaultCmdBuilder This is to build the default Cmd to execute the plugin
var DefaultCmdBuilder BuildCmd = func(filePath string, args string) *exec.Cmd {
	arguments := strings.Split(args, " ")
	arguments = append([]string{"plugin"}, arguments...)
	return exec.Command(filePath, arguments...)
}

//Close Let the client kill the method
func (clientCloser *ClientCloser) Close() {
	clientCloser.Client.Kill()
}

//New - method to create a plugin tile
func (s *PluginTileBuilder) New(tileSpec tileregistry.TileSpec) (tileCloser tileregistry.TileCloser, err error) {
	var opsManager *opsmanager.OpsManager
	var settingsReader io.Reader
	opsManager, err = opsmanager.NewOpsManager(tileSpec.OpsManagerHost, tileSpec.AdminUser, tileSpec.AdminPass, tileSpec.OpsManagerUser, tileSpec.OpsManagerPass, tileSpec.OpsManagerPassphrase, tileSpec.ArchiveDirectory, tileSpec.CryptKey)

	if settingsReader, err = opsManager.GetInstallationSettings(); err == nil {
		var brPlugin BackupRestorer
		installationSettings := cfbackup.NewConfigurationParserFromReader(settingsReader)
		pcf := NewPivotalCF(installationSettings.InstallationSettings, tileSpec)
		lo.G.Debug("", s.Meta.Name, s.FilePath, pcf)
		brPlugin, client := s.call(tileSpec)
		brPlugin.Setup(pcf)
		tileCloser = struct {
			tileregistry.Tile
			tileregistry.Closer
		}{
			brPlugin,
			&ClientCloser{Client: client},
		}
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
