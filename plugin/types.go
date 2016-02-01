package plugin

import "github.com/pivotalservices/cfbackup"

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
type Product cfbackup.Products

//Credential - credential implementation
type Credential cfbackup.Properties

//PivotalCF - interface representing a pivotalcf
type PivotalCF interface {
	SetActivity(string)
	GetActivity() string
	GetProducts() map[string]cfbackup.Products
	GetCredentials() map[string]map[string][]cfbackup.Properties
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
	PivotalCF PivotalCF
	FilePath  string
	Meta      Meta
}

type defaultPivotalCF struct {
	activity             string
	installationSettings *cfbackup.ConfigurationParser
}
