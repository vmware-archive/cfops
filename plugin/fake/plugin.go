package fake

import (
	"errors"

	"github.com/pivotalservices/cfops/plugin/cfopsplugin"
	"github.com/xchapter7x/lo"
)

//Plugin -- Here is a real implementation of a plugin
type Plugin struct {
	Meta cfopsplugin.Meta
}

//GetMeta ---
func (b Plugin) GetMeta() cfopsplugin.Meta {
	return b.Meta
}

//Setup --
func (Plugin) Setup(pcf cfopsplugin.PivotalCF) error {
	lo.G.Debug("mypcf: ", pcf)
	return nil
}

//Backup --
func (Plugin) Backup() error {
	lo.G.Debug("Backup!")
	lo.G.Debug("again")
	return errors.New("somethign happened")
}

//Restore --
func (Plugin) Restore() error {
	lo.G.Debug("restore!")
	return nil
}
