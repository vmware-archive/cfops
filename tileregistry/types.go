package tileregistry

type (
	//TileGenerator - interface for a tile creating object
	TileGenerator interface {
		New(tileSpec TileSpec) (Tile, error)
	}
	//Tile - definition for what a tile looks like
	Tile interface {
		Backup() error
		Restore() error
	}
	//TileSpec -- defines what a tile would need to be initialized
	TileSpec struct {
		OpsManagerHost    string
		AdminUser         string
		AdminPass         string
		OpsManagerUser    string
		OpsManagerPass    string
		ArchiveDirectory  string
		CryptKey          string
		ClearBoshManifest bool
		PluginArgs        string
	}
)
