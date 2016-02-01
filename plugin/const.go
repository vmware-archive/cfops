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

	//PluginActivityRestore - activity keyname
	PluginActivityRestore = "restore"
	//PluginActivityBackup - activity keyname
	PluginActivityBackup = "backup"
	//PluginActivityExecute - activity keyname
	PluginActivityExecute = "execute"
)

var (
	//UIOutput - a function to control UIOutput
	UIOutput = fmt.Print
	//PluginArgs - cli args given to the plugin
	PluginArgs = os.Args
)
