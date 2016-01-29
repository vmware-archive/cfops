package fake

import "github.com/pivotalservices/cfops/plugin"

//BackupRestorePlugin --
type BackupRestorePlugin struct {
	Meta         plugin.Meta
	RunCallCount int
	SpyPivotalCF plugin.PivotalCF
}

//GetMeta --
func (s *BackupRestorePlugin) GetMeta() (meta plugin.Meta) {
	return s.Meta
}

//Run --
func (s *BackupRestorePlugin) Run(pcf plugin.PivotalCF, args *[]string) (err error) {
	s.RunCallCount++
	s.SpyPivotalCF = pcf
	return nil
}
