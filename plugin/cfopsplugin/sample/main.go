package main

import (
	"errors"
	"os"

	cfopsplugin "github.com/pivotalservices/cfops/plugin/cfopsplugin"
	"github.com/pivotalservices/cfops/plugin/fake"
	"github.com/xchapter7x/lo"
)

func main() {

	if len(os.Args) == 2 {
		lo.G.Debug("executing run plugin", os.Args)
		runPlugin()
	} else {
		lo.G.Debug("executing call plugin", os.Args)
		callPlugin()
	}
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
	return nil
}

func runPlugin() {
	brt := new(BackupRestoreTest)
	brt.Meta = cfopsplugin.Meta{Name: "backuprestore", Role: "backup-restore"}
	cfopsplugin.Start(brt)
}

func callPlugin() {
	p, c := cfopsplugin.Call("backuprestore", os.Args[0])
	defer c.Kill()
	lo.G.Debug("", p.Setup(new(fake.PivotalCF)))
	lo.G.Debug("", p.Backup())
	lo.G.Debug("", p.Restore())
}
