package load

import "errors"

//PluginDir - default plugin directory
var PluginDir = "./plugins"

//ErrInvalidPluginMeta - plugin error
var ErrInvalidPluginMeta = errors.New("invalid plugin meta")

const (
	//PluginDirEnv - env var to set a custom plugin directory
	PluginDirEnv = "CFOPS_PLUGIN_DIR"
)
