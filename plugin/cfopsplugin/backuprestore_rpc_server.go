package cfopsplugin

import "github.com/xchapter7x/lo"

//Backup --
func (s *BackupRestoreRPCServer) Backup(args interface{}, resp *error) error {
	lo.G.Debug("rpc server backup execution")
	return s.Impl.Backup()
}

//Restore --
func (s *BackupRestoreRPCServer) Restore(args interface{}, resp *error) error {
	lo.G.Debug("rpc server restore execution")
	return s.Impl.Restore()
}

//Restore --
func (s *BackupRestoreRPCServer) Setup(pcf PivotalCF, resp *error) error {
	lo.G.Debug("rpc server setup execution")
	return s.Impl.Setup(pcf)
}
