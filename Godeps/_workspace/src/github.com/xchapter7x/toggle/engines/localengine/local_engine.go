package localengine

import (
	"fmt"
	"os"

	"github.com/xchapter7x/toggle/engines/storageinterface"
)

func NewLocalEngine() (engine storageinterface.StorageEngine) {
	engine = &LocalEngine{
		Getenv: os.Getenv,
	}
	return
}

type LocalEngine struct {
	Getenv func(string) string
}

func (s *LocalEngine) GetFeatureStatusValue(featureSignature string) (status string, err error) {
	status = s.Getenv(featureSignature)

	if status == "" {
		err = fmt.Errorf("toggle value not set")
	}
	return
}

func (s *LocalEngine) Close() (err error) {
	return
}
