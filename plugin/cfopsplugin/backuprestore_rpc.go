package cfopsplugin

import (
	"errors"

	"github.com/xchapter7x/lo"
)

//Backup --
func (g *BackupRestoreRPC) Backup() (err error) {
	lo.G.Debug("running backuprestorerpc backup: ")
	resp := err
	err = g.client.Call("Plugin.Backup", new(interface{}), &resp)
	lo.G.Debug("done calling plugin.backup: ", err)
	return
}

//Restore --
func (g *BackupRestoreRPC) Restore() (err error) {
	lo.G.Debug("running backuprestorerpc restore: ")
	resp := errors.New("")
	err = g.client.Call("Plugin.Restore", new(interface{}), &resp)
	lo.G.Debug("done calling plugin.restore: ", err)
	return
}

//Restore --
func (g *BackupRestoreRPC) Setup(pcf PivotalCF) error {
	lo.G.Debug("running backuprestorerpc setup ", pcf)
	resp := errors.New("")
	err := g.client.Call("Plugin.Setup", &pcf, &resp)
	lo.G.Debug("done calling plugin.setup: ", err)
	return err
}
