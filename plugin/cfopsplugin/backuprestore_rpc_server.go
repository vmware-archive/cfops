package cfopsplugin

//Backup --
func (s *BackupRestoreRPCServer) Backup(args interface{}, resp *error) error {
	*resp = s.Impl.Backup()
	return *resp
}

//Restore --
func (s *BackupRestoreRPCServer) Restore(args interface{}, resp *error) error {
	*resp = s.Impl.Restore()
	return *resp
}

//Restore --
func (s *BackupRestoreRPCServer) Setup(pcf PivotalCF, resp *error) error {
	*resp = s.Impl.Setup(pcf)
	return *resp
}
