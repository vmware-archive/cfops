package cfopsplugin

import (
	"encoding/gob"
	"encoding/json"
	"os"

	"github.com/hashicorp/go-plugin"
	"github.com/xchapter7x/lo"
)

//Start - takes a given plugin and starts it
func Start(plgn Plugin) {
	gob.Register(plgn)
	RegisterPlugin(plgn.GetMeta().Name, plgn)

	if len(os.Args) == 2 && os.Args[1] == PluginMeta {
		b, _ := json.Marshal(plgn.GetMeta())
		UIOutput(string(b))

	} else {

		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: handshakeConfig,
			Plugins:         GetPlugins(),
		})
	}
}

//RegisterPlugin - register a plugin as available
func RegisterPlugin(name string, plugin BackupRestorer) {
	lo.G.Debug("registering plugin: ", name, plugin)
	pluginMap[name] = &BackupRestorePlugin{
		P: plugin,
	}
}

//GetPlugins - returns the list of registered plugins
func GetPlugins() map[string]plugin.Plugin {
	return pluginMap
}

//GetHandshake - gets the handshake object
func GetHandshake() plugin.HandshakeConfig {
	return handshakeConfig
}
