package cfopsplugin

import (
	"io"
	"net/rpc"
	"os/exec"

	"github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfops/tileregistry"
)

//Meta - plugin meta data storage object
type Meta struct {
	Name                string
	Role                string
	Description         string
	SupportedActivities map[string]bool
}

//BackupRestorer - is the interface that we're exposing as a plugin.
type BackupRestorer interface {
	Backup() error
	Restore() error
	Setup(PivotalCF) error
}

//Plugin - is a interface plugin providers should implement
type Plugin interface {
	BackupRestorer
	GetMeta() Meta
}

//PivotalCF - interface representing a pivotalcf
type PivotalCF interface {
	GetHostDetails() tileregistry.TileSpec
	NewArchiveReader(name string) (io.ReadCloser, error)
	NewArchiveWriter(name string) (io.WriteCloser, error)
	GetInstallationSettings() cfbackup.InstallationSettings
}

//BackupRestorePlugin - this is an implementation of the rpc client and server wrapper
type BackupRestorePlugin struct {
	P BackupRestorer
}

//BackupRestoreRPCServer - this is an implementation of the rpc server
//for a backuprestorer
type BackupRestoreRPCServer struct {
	Impl BackupRestorer
}

//BackupRestoreRPC - is an implementation of a client that talks over RPC
type BackupRestoreRPC struct {
	client *rpc.Client
}

//DefaultPivotalCF - default implementation of PivotalCF interface
type DefaultPivotalCF struct {
	TileSpec             tileregistry.TileSpec
	InstallationSettings cfbackup.InstallationSettings
}

//PluginTileBuilder - factory for a tile wrapped plugin
type PluginTileBuilder struct {
	FilePath   string
	Meta       Meta
	CmdBuilder BuildCmd
}

//BuildCmd Command func
type BuildCmd func(filePath string, args string) *exec.Cmd
