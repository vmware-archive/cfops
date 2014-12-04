package localpubsub

import (
	"os"

	"github.com/xchapter7x/toggle"
	"github.com/xchapter7x/toggle/engines/localengine"
	"github.com/xchapter7x/toggle/engines/storageinterface"
)

func NewLocalPubSubEngine(pubsub pubsubInterface, toggleList map[string]*toggle.Feature) storageinterface.StorageEngine {
	le := &localengine.LocalEngine{
		Getenv: os.Getenv,
	}
	engine := &LocalPubSubEngine{
		LocalEngine: le,
		PubSub:      pubsub,
	}
	engine.StartSubscriptionListener(toggleList)
	return engine
}

type LocalPubSubEngine struct {
	LocalEngine *localengine.LocalEngine
	PubSub      pubsubInterface
	quit        chan bool
}

func (s *LocalPubSubEngine) Close() (err error) {
	s.PubSub.Unsubscribe()
	s.quit <- true
	return
}

func (s *LocalPubSubEngine) StartSubscriptionListener(toggleList map[string]*toggle.Feature) {
	if s.quit == nil {
		s.quit = make(chan bool)

		go func() {
			for {
				select {
				case <-s.quit:
					return

				default:
					PubSubReceiver(s.PubSub, toggleList)
				}
			}
		}()
	}
}

func (s *LocalPubSubEngine) GetFeatureStatusValue(featureSignature string) (status string, err error) {
	s.PubSub.Subscribe(featureSignature)
	status, err = s.LocalEngine.GetFeatureStatusValue(featureSignature)
	return
}
