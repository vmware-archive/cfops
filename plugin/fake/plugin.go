package fake

import "github.com/pivotalservices/cfops/plugin/cfopsplugin"

//BackupRestorePlugin --
type BackupRestorePlugin struct {
	Meta         cfopsplugin.Meta
	RunCallCount int
	SpyPivotalCF cfopsplugin.PivotalCF
}

//GetMeta --
func (s *BackupRestorePlugin) GetMeta() (meta cfopsplugin.Meta) {
	return s.Meta
}

//Run --
func (s *BackupRestorePlugin) Run(pcf cfopsplugin.PivotalCF, args *[]string) (err error) {
	s.RunCallCount++
	s.SpyPivotalCF = pcf
	return nil
}
