package main

import (
	"fmt"

	"github.com/pivotalservices/cfops/plugin"
	"github.com/xchapter7x/lo"
)

type brPlugin struct {
	Meta plugin.Meta
}

func (s brPlugin) GetMeta() plugin.Meta {
	return s.Meta
}

func (s brPlugin) Run(pcf plugin.PivotalCF, args *[]string) (err error) {
	lo.G.Debug("we are here")
	lo.G.Debug("pcf: ", pcf)
	lo.G.Debug("args", args)

	switch pcf.GetActivity() {
	case "backup":
		fmt.Println("this is a backup")
	case "restore":
		fmt.Println("this is a restore")
	default:
		fmt.Println("not sure what this is")
	}
	return
}

func main() {
	myPlugin := new(brPlugin)
	myPlugin.Meta = plugin.Meta{
		Name:        "something",
		Role:        "backup-restore",
		Description: "this plugin doesnt do anything",
	}
	plugin.Start(myPlugin)
}
