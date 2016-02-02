package cfopsplugin

//Backup --
func (g *BackupRestoreRPC) Backup() error {
	var resp error
	err := g.client.Call("Plugin.Backup", new(interface{}), &resp)
	return err
}

//Restore --
func (g *BackupRestoreRPC) Restore() error {
	var resp error
	err := g.client.Call("Plugin.Restore", new(interface{}), &resp)
	return err
}

//Restore --
func (g *BackupRestoreRPC) Setup(pcf PivotalCF) error {
	var resp error
	err := g.client.Call("Plugin.Setup", &pcf, &resp)
	return err
}
