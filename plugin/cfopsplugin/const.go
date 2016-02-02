package cfopsplugin

import (
	"fmt"

	"github.com/hashicorp/go-plugin"
)

const (
	//PluginMeta - default plugin arg to show meta data
	PluginMeta = "plugin-meta"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

var pluginMap = make(map[string]plugin.Plugin)

var (
	//UIOutput - a function to control UIOutput
	UIOutput = fmt.Print
)
