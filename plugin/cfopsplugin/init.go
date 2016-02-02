package cfopsplugin

import "encoding/gob"

func init() {
	gob.Register(new(DefaultPivotalCF))
}
