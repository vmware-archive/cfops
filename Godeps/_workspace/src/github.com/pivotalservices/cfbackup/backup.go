package cfbackup

import "path"

const (
	//BackupLoggerName --
	BackupLoggerName = "Backup"
	//RestoreLoggerName --
	RestoreLoggerName = "Restore"
)

var (
	//TileRestoreAction -- executes a restore action on the given tile
	TileRestoreAction = func(t Tile) func() error {
		return t.Restore
	}
	//TileBackupAction - executes a backup action on a given tile
	TileBackupAction = func(t Tile) func() error {
		return t.Backup
	}
)

//RunBackupPipeline - Backup the list of all default tiles
func RunBackupPipeline(hostname, adminUsername, adminPassword, opsManagerUsername, opsManagerPassword, destination string) (err error) {
	var tiles []Tile
	conn := connectionBucket{
		hostname:           hostname,
		adminUsername:      adminUsername,
		adminPassword:      adminPassword,
		opsManagerUsername: opsManagerUsername,
		opsManagerPassword: opsManagerPassword,
		destination:        destination,
	}

	if tiles, err = fullTileList(conn, BackupLoggerName); err == nil {
		err = RunPipeline(TileBackupAction, tiles)
	}
	return
}

//RunRestorePipeline - Restore the list of all default tiles
func RunRestorePipeline(hostname, adminUsername, adminPassword, opsManagerUser, opsManagerPassword, destination string) (err error) {
	var tiles []Tile
	conn := connectionBucket{
		hostname:           hostname,
		adminUsername:      adminUsername,
		adminPassword:      adminPassword,
		opsManagerUsername: opsManagerUser,
		opsManagerPassword: opsManagerPassword,
		destination:        destination,
	}

	if tiles, err = fullTileList(conn, RestoreLoggerName); err == nil {
		err = RunPipeline(TileRestoreAction, tiles)
	}
	return
}

//Runs a pipeline action (restore/backup) on a list of tiles
var RunPipeline = func(actionBuilder func(Tile) func() error, tiles []Tile) (err error) {
	var pipeline []action

	for _, tile := range tiles {
		pipeline = append(pipeline, actionBuilder(tile))
	}
	err = runActions(pipeline)
	return
}

func runActions(actions []action) (err error) {
	for _, action := range actions {

		if err = action(); err != nil {
			break
		}
	}
	return
}

func fullTileList(conn connBucketInterface, loggerName string) (tiles []Tile, err error) {
	var (
		opsmanager     Tile
		elasticRuntime Tile
	)
	installationFilePath := path.Join(conn.Destination(), OpsMgrBackupDir, OpsMgrInstallationSettingsFilename)

	if opsmanager, err = NewOpsManager(conn.Host(), conn.AdminUser(), conn.AdminPass(), conn.OpsManagerUser(), conn.OpsManagerPass(), conn.Destination()); err == nil {
		elasticRuntime = NewElasticRuntime(installationFilePath, conn.Destination())
		tiles = []Tile{
			opsmanager,
			elasticRuntime,
		}
	}
	return
}
