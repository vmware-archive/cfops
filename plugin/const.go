package plugin

import (
	"fmt"
	"os"
)

const (
	//PluginPort - default plugin port
	PluginPort = int(1984)
	//PluginMeta - default plugin arg to show meta data
	PluginMeta = "plugin-meta"
)

var (
	//UIOutput - a function to control UIOutput
	UIOutput = fmt.Print
	//PluginArgs - cli args given to the plugin
	PluginArgs = os.Args
)
