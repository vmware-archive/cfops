package tileregistry

//Register -- add a Sku interface object to the Repo
func Register(name string, tile TileGenerator) {
	Repo[name] = tile
}

//GetRegistry -- gets the map of all registered Sku interface objects
func GetRegistry() map[string]TileGenerator {
	return Repo
}
