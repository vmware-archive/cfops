package fake

import (
	"github.com/pivotalservices/cfops/tileregistry"
	"github.com/xchapter7x/lo"
)

//TileGenerator --
type TileGenerator struct {
	tileregistry.TileGenerator
	TileSpy tileregistry.Tile
	ErrFake error
}

//New --
func (s *TileGenerator) New(tileSpec tileregistry.TileSpec) (tileregistry.Tile, error) {
	return s.TileSpy, s.ErrFake
}

//Tile --
type Tile struct {
	ErrFake          error
	BackupCallCount  int
	RestoreCallCount int
}

//Backup --
func (s *Tile) Backup() error {
	lo.G.Debug("we fake backed up")
	s.BackupCallCount++
	return s.ErrFake
}

//Restore --
func (s *Tile) Restore() error {
	lo.G.Debug("we fake restored")
	s.RestoreCallCount++
	return s.ErrFake
}
