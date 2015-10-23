package cfops

import (
	"fmt"
	"path"
	"strings"

	"github.com/pivotalservices/cfbackup"
	"github.com/xchapter7x/lo"
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
	BuiltinPipelineExecution = map[string]func(string, string, string, string, string, string) error{
		Restore: cfbackup.RunRestorePipeline,
		Backup:  cfbackup.RunBackupPipeline,
	}
	SupportedTiles map[string]func() (Tile, error)
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
	AdminUser() string
	AdminPass() string
	OpsManagerUser() string
	OpsManagerPass() string
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
			opsmgr, err = cfbackup.NewOpsManager(fs.Host(), fs.AdminUser(), fs.AdminPass(), fs.OpsManagerUser(), fs.OpsManagerPass(), fs.Dest())
			lo.G.Debug("Creating a new OpsManager object")
			return
		},
		ER: func() (er Tile, err error) {
			installationFilePath := path.Join(fs.Dest(), cfbackup.OPSMGR_BACKUP_DIR, cfbackup.OPSMGR_INSTALLATION_SETTINGS_FILENAME)
			er = cfbackup.NewElasticRuntime(installationFilePath, fs.Dest())
			lo.G.Debug("Creating a new ElasticRuntime object")
			return
		},
	}
}

func runTileUsingAction(t Tile, action string) (err error) {
	lo.G.Debug("Running on tile")
	switch action {
	case Restore:
		err = t.Restore()

	case Backup:
		err = t.Backup()
	}
	lo.G.Debug("Action complete")
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

func RunPipeline(fs flagSet, action string) (err error) {

	if hasTilelistFlag(fs) {
		lo.G.Debug("Running a tile list action")
		err = runTileListUsingAction(fs, action)

	} else {
		err = BuiltinPipelineExecution[action](fs.Host(), fs.AdminUser(), fs.AdminPass(), fs.OpsManagerUser(), fs.OpsManagerPass(), fs.Dest())
	}
	return
}
