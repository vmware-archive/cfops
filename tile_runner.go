package cfops

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/gtils/log"
)

const (
	ErrUnsupportedTileFormat = "you have a unsupported tile in your list: %s"
	loggerName               = "cfops"
	Restore                  = "restore"
	Backup                   = "backup"
	OpsMgr                   = "OPSMANAGER"
	ER                       = "ER"
)

var (
	BuiltinPipelineExecution = map[string]func(string, string, string, string, string) error{
		Restore: cfbackup.RunRestorePipeline,
		Backup:  cfbackup.RunBackupPipeline,
	}
	SupportedTiles map[string]func() (Tile, error)
	backupLogger   log.Logger = log.LogFactory("cfops default logger", log.Lager, os.Stdout)
)

func ErrUnsupportedTile(errString string) error {
	return fmt.Errorf(ErrUnsupportedTileFormat, errString)
}

type Tile interface {
	Backup() error
	Restore() error
}

func hasTilelistFlag(fs flagSet) bool {
	return (fs.Tilelist() != "")
}

type flagSet interface {
	Host() string
	User() string
	Pass() string
	Tpass() string
	Dest() string
	Tilelist() string
}

func formatArray(a []string) []string {
	for i, v := range a {
		a[i] = strings.ToUpper(strings.TrimSpace(v))
	}
	return a
}

func SetupSupportedTiles(fs flagSet) {
	SupportedTiles = map[string]func() (Tile, error){
		OpsMgr: func() (opsmgr Tile, err error) {
			opsmgr, err = cfbackup.NewOpsManager(fs.Host(), fs.User(), fs.Pass(), fs.Tpass(), fs.Dest(), backupLogger)
			backupLogger.Debug("Creating a new OpsManager object", log.Data{"err": err, "opsmgr": opsmgr})
			return
		},
		ER: func() (er Tile, err error) {
			installationFilePath := path.Join(fs.Dest(), cfbackup.OPSMGR_BACKUP_DIR, cfbackup.OPSMGR_INSTALLATION_SETTINGS_FILENAME)
			er = cfbackup.NewElasticRuntime(installationFilePath, fs.Dest(), backupLogger)
			backupLogger.Debug("Creating a new ElasticRuntime object", log.Data{"er": er})
			return
		},
	}
}

func runTileUsingAction(t Tile, action string) (err error) {
	backupLogger.Debug("Running on tile", log.Data{"tile": t, "action": action})
	switch action {
	case Restore:
		err = t.Restore()

	case Backup:
		err = t.Backup()
	}
	backupLogger.Debug("Action complete", log.Data{"tile": t, "action": action, "err": err})
	return
}

func getSupportedTile(tilename string) (tile Tile, err error) {
	var (
		ok          bool
		tileFactory func() (Tile, error)
	)

	if tileFactory, ok = SupportedTiles[tilename]; !ok {
		err = ErrUnsupportedTile(tilename)

	} else {
		tile, err = tileFactory()
	}
	return
}

func runTileListUsingAction(fs flagSet, action string) (err error) {
	tiles := formatArray(strings.Split(fs.Tilelist(), ","))

	for _, tileName := range tiles {
		var tile Tile

		if tile, err = getSupportedTile(tileName); err == nil {
			err = runTileUsingAction(tile, action)
		}

		if err != nil {
			break
		}
	}
	return
}

func SetLogger(logger log.Logger) {
	backupLogger = logger
}

func RunPipeline(fs flagSet, action string) (err error) {

	if hasTilelistFlag(fs) {
		backupLogger.Debug("Running a tile list action", nil)
		err = runTileListUsingAction(fs, action)

	} else {
		cfbackup.SetLogger(backupLogger)
		err = BuiltinPipelineExecution[action](fs.Host(), fs.User(), fs.Pass(), fs.Tpass(), fs.Dest())
	}
	return
}
