package plugin

import "github.com/pivotalservices/cfbackup/tiles/opsmanager"

//Meta - plugin meta data storage object
type Meta struct {
	Name                string
	Role                string
	Description         string
	SupportedActivities map[string]bool
}

//Plugin - a plugin interface definition
type Plugin interface {
	GetMeta() Meta
	Run(PivotalCF, *[]string) error
}

//Product - implementation
type Product struct{}

//Credential - credential implementation
type Credential struct{}

//PivotalCF - interface representing a pivotalcf
type PivotalCF interface {
	GetActivity() string
	GetProducts() []Product
	GetCredentials() []Credential
}

type wrappedPlugin struct {
	plugin Plugin
}

//PluginTileBuilder - factory for a tile wrapped plugin
type PluginTileBuilder struct {
	FilePath string
	Meta     Meta
}

//PluginTile - tile implementation of a plugin
type PluginTile struct {
	OpsManager *opsmanager.OpsManager
	FilePath   string
	Meta       Meta
}
