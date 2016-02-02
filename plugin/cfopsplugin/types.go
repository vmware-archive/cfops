package cfopsplugin

import (
	"net/rpc"

	"github.com/pivotalservices/cfbackup"
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

type Plugin interface {
	BackupRestorer
	GetMeta() Meta
}

//PivotalCF - interface representing a pivotalcf
type PivotalCF interface {
	GetProducts() map[string]cfbackup.Products
	GetCredentials() map[string]map[string][]cfbackup.Properties
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
	activity             string
	installationSettings *cfbackup.ConfigurationParser
}

//PluginTileBuilder - factory for a tile wrapped plugin
type PluginTileBuilder struct {
	FilePath string
	Meta     Meta
}

//PluginTile - tile implementation of a plugin
type PluginTile struct {
	PivotalCF PivotalCF
	FilePath  string
	Meta      Meta
}
