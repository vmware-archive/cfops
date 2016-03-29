package main

import (
	"errors"
	"os"

	cfopsplugin "github.com/pivotalservices/cfops/plugin/cfopsplugin"
	"github.com/xchapter7x/lo"
)

func main() {
	runPlugin()
}

// Here is a real implementation of a plugin
type BackupRestoreTest struct {
	Meta cfopsplugin.Meta
}

//GetMeta ---
func (b BackupRestoreTest) GetMeta() cfopsplugin.Meta {
	return b.Meta
}

//Setup --
func (BackupRestoreTest) Setup(pcf cfopsplugin.PivotalCF) error {
	lo.G.Debug("mypcf: ", pcf)
	return nil
}

//Backup --
func (BackupRestoreTest) Backup() error {
	lo.G.Debug("Backup!")
	lo.G.Debug("again")
	return errors.New("somethign happened")
}

//Restore --
func (BackupRestoreTest) Restore() error {
	lo.G.Debug("restore!")
	lo.G.Debug("Arguments %s", os.Args[2:])
	return nil
}

func runPlugin() {
	brt := new(BackupRestoreTest)
	brt.Meta = cfopsplugin.Meta{Name: "backuprestore", Role: "backup-restore"}
	cfopsplugin.Start(brt)
}
