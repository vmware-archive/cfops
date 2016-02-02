package cfopsplugin

import "encoding/gob"

func init() {
	gob.Register(make(map[string]interface{}))
	gob.Register(new(DefaultPivotalCF))
	gob.Register(new(BackupRestoreRPC))
}
