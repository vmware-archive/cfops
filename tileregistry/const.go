package tileregistry

var (
	//Repo -- repo holds the registered sku interfaces
	Repo = make(map[string]TileGenerator)
)
