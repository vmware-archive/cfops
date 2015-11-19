package fake

import (
	"github.com/pivotalservices/cfops/tileregistry"
	"github.com/xchapter7x/lo"
)

type TileGenerator struct {
	tileregistry.TileGenerator
	TileSpy tileregistry.Tile
	ErrFake error
}

func (s *TileGenerator) New(tileSpec tileregistry.TileSpec) (tileregistry.Tile, error) {
	return s.TileSpy, s.ErrFake
}

type Tile struct {
	ErrFake          error
	BackupCallCount  int
	RestoreCallCount int
}

func (s *Tile) Backup() error {
	lo.G.Debug("we fake backed up")
	s.BackupCallCount++
	return s.ErrFake
}

func (s *Tile) Restore() error {
	lo.G.Debug("we fake restored")
	s.RestoreCallCount++
	return s.ErrFake
}
