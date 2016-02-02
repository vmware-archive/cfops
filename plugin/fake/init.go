package fake

import "encoding/gob"

func init() {
	gob.Register(new(PivotalCF))
}
