package backup

import "path"

// Tile is a deployable component that can be backed up
type Tile interface {
	Backup() error
	Restore() error
}

type BackupContext struct {
	TargetDir string
}

func RunBackupPipeline(hostname, username, password, tempestpassword, destination string) (err error) {
	var (
		opsmanager     Tile
		elasticRuntime Tile
	)
	installationFilePath := path.Join(destination, OPSMGR_BACKUP_DIR, OPSMGR_INSTALLATION_SETTINGS_FILENAME)

	if opsmanager, err = NewOpsManager(hostname, username, password, tempestpassword, destination); err == nil {
		elasticRuntime = NewElasticRuntime(installationFilePath, destination)
		tiles := []Tile{
			opsmanager,
			elasticRuntime,
		}
		err = runBackups(tiles)
	}
	return
}

func runBackups(tiles []Tile) (err error) {
	for _, tile := range tiles {

		if err = tile.Backup(); err != nil {
			break
		}
	}
	return
}
