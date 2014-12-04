package localpubsub

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/xchapter7x/toggle"
)

type receiverInterface interface {
	Receive() interface{}
}

func PubSubReceiver(s receiverInterface, toggleList map[string]*toggle.Feature) {
	switch n := s.Receive().(type) {
	case redis.Message:
		toggleList[n.Channel].UpdateStatus(string(n.Data[:]))

	case error:
		fmt.Printf("error: %v\n", n)
		return
	}
}
