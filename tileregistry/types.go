package tileregistry

type (
	//TileGenerator - interface for a tile creating object
	TileGenerator interface {
		New() Tile
	}
	//Tile - definition for what a tile looks like
	Tile interface {
		Backup() error
		Restore() error
	}
)
