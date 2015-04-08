package cfops

import (
	"fmt"

	"github.com/pivotalservices/cfbackup"
)

const (
	Restore = "restore"
	Backup  = "backup"
)

var (
	BuiltinPipelineExecution = map[string]func(string, string, string, string, string) error{
		Restore: cfbackup.RunRestorePipeline,
		Backup:  cfbackup.RunBackupPipeline,
	}
)

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

func RunPipeline(fs flagSet, action string) (err error) {

	if hasTilelistFlag(fs) {
		err = fmt.Errorf("tile list is not yet implemented")

	} else {
		err = BuiltinPipelineExecution[action](fs.Host(), fs.User(), fs.Pass(), fs.Tpass(), fs.Dest())
	}
	return
}
