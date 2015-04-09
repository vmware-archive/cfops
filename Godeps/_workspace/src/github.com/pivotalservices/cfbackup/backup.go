package cfbackup

import (
	"os"
	"path"

	"github.com/pivotalservices/gtils/log"
)

const (
	BACKUP_LOGGER_NAME  = "Backup"
	RESTORE_LOGGER_NAME = "Restore"
)

var (
	TILE_RESTORE_ACTION = func(t Tile) func() error {
		return t.Restore
	}
	TILE_BACKUP_ACTION = func(t Tile) func() error {
		return t.Backup
	}
	backupLogger log.Logger
)

// Tile is a deployable component that can be backed up
type Tile interface {
	Backup() error
	Restore() error
}

type connBucketInterface interface {
	Host() string
	User() string
	Pass() string
	TempestPass() string
	Destination() string
}

type BackupContext struct {
	TargetDir string
}

type action func() error

type actionAdaptor func(t Tile) action

//SetLogger - lets us set the logger object
func SetLogger(logger log.Logger) {
	backupLogger = logger
}

func init() {
	backupLogger = log.LogFactory("cfbackup default logger", log.Lager, os.Stdout)
}

//Backup the list of all default tiles
func RunBackupPipeline(hostname, username, password, tempestpassword, destination string) (err error) {
	var tiles []Tile
	conn := connectionBucket{
		hostname:        hostname,
		username:        username,
		password:        password,
		tempestPassword: tempestpassword,
		destination:     destination,
	}

	if tiles, err = fullTileList(conn, BACKUP_LOGGER_NAME); err == nil {
		err = RunPipeline(TILE_BACKUP_ACTION, tiles)
	}
	return
}

//Restore the list of all default tiles
func RunRestorePipeline(hostname, username, password, tempestpassword, destination string) (err error) {
	var tiles []Tile
	conn := connectionBucket{
		hostname:        hostname,
		username:        username,
		password:        password,
		tempestPassword: tempestpassword,
		destination:     destination,
	}

	if tiles, err = fullTileList(conn, RESTORE_LOGGER_NAME); err == nil {
		err = RunPipeline(TILE_RESTORE_ACTION, tiles)
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
	installationFilePath := path.Join(conn.Destination(), OPSMGR_BACKUP_DIR, OPSMGR_INSTALLATION_SETTINGS_FILENAME)

	if opsmanager, err = NewOpsManager(conn.Host(), conn.User(), conn.Pass(), conn.TempestPass(), conn.Destination(), backupLogger); err == nil {
		elasticRuntime = NewElasticRuntime(installationFilePath, conn.Destination(), backupLogger)
		tiles = []Tile{
			opsmanager,
			elasticRuntime,
		}
	}
	return
}
