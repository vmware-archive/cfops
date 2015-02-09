package cfbackup

import (
	"github.com/cloudfoundry-incubator/cf-lager"
	"path"
)

// Tile is a deployable component that can be backed up
type Tile interface {
	Backup() error
	Restore() error
}

type BackupContext struct {
	TargetDir string
}

type action func() error

type actionAdaptor func(t Tile) action

func RunBackupPipeline(hostname, username, password, tempestpassword, destination string) (err error) {
	backup := func(t Tile) action {
		return func() error {
			return t.Backup()
		}
	}
	return runPipelines(hostname, username, password, tempestpassword, destination, "Backup", backup)
}

func RunRestorePipeline(hostname, username, password, tempestpassword, destination string) (err error) {
	restore := func(t Tile) action {
		return func() error {
			return t.Restore()
		}
	}
	return runPipelines(hostname, username, password, tempestpassword, destination, "Backup", restore)
}

func runPipelines(hostname, username, password, tempestpassword, destination, loggerName string, actionBuilder actionAdaptor) (err error) {
	var (
		opsmanager     Tile
		elasticRuntime Tile
	)
	backupLogger := cf_lager.New(loggerName)
	installationFilePath := path.Join(destination, OPSMGR_BACKUP_DIR, OPSMGR_INSTALLATION_SETTINGS_FILENAME)

	if opsmanager, err = NewOpsManager(hostname, username, password, tempestpassword, destination, backupLogger); err == nil {
		elasticRuntime = NewElasticRuntime(installationFilePath, destination, backupLogger)
		tiles := []action{
			actionBuilder(opsmanager),
			actionBuilder(elasticRuntime),
		}
		err = runActions(tiles)
	}
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
