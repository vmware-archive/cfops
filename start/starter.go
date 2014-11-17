package start

import (
	"github.com/pivotalservices/cfops/cli"
)

type Starter struct {
}

func New(factory cli.CommandFactory) Starter {

	factory.Register("start", StartCommand{})

	return Starter{}
}
