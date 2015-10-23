package lo

import (
	"os"

	"github.com/op/go-logging"
)

const (
	LOG_MODULE = "lo.G_logger"
)

var (
	G *logging.Logger
)

func init() {

	if logLevel, err := logging.LogLevel(os.Getenv("LOG_LEVEL")); err == nil {
		logging.SetLevel(logLevel, LOG_MODULE)

	} else {
		logging.SetLevel(logging.INFO, LOG_MODULE)
	}
	logging.SetFormatter(logging.GlogFormatter)
	G = logging.MustGetLogger(LOG_MODULE)
}
