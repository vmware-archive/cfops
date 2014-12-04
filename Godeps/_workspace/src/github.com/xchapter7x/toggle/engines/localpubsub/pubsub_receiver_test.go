package localpubsub_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/garyburd/redigo/redis"
	"github.com/xchapter7x/toggle"
	"github.com/xchapter7x/toggle/engines/localpubsub"
)

type PubSubErrorMock struct {
	Channel string
	Data    []byte
}

func (s *PubSubErrorMock) Receive() (rMsg interface{}) {
	rMsg = "failure not a redis message"
	return
}

type PubSubMock struct {
	Channel string
	Data    []byte
}

func (s *PubSubMock) Receive() (rMsg interface{}) {
	rMsg = redis.Message{
		Channel: s.Channel,
		Data:    s.Data,
	}
	return
}

var _ = Describe("localpubsub package", func() {
	Describe("PubSubReceiver function", func() {
		channel := "test"
		data := []byte("string")
		var psMock *PubSubMock
		var psErrorMock *PubSubErrorMock

		BeforeEach(func() {
			psMock = &PubSubMock{
				Channel: channel,
				Data:    data,
			}
			psErrorMock = &PubSubErrorMock{
				Channel: channel,
				Data:    data,
			}

		})

		AfterEach(func() {
			psMock = nil
			psErrorMock = nil
		})

		It("Should update the status of the proper feature object", func() {
			controlStatus := "controlString"
			togglelist := make(map[string]*toggle.Feature)
			togglelist[channel] = &toggle.Feature{
				Status: controlStatus,
			}
			localpubsub.PubSubReceiver(psMock, togglelist)
			Expect(togglelist[channel].Status).NotTo(Equal(controlStatus))
			Expect(togglelist[channel].Status).To(Equal(string(data[:])))
		})

		It("Should on error leave existing values in place", func() {
			controlStatus := "controlString"
			togglelist := make(map[string]*toggle.Feature)
			togglelist[channel] = &toggle.Feature{
				Status: controlStatus,
			}
			localpubsub.PubSubReceiver(psErrorMock, togglelist)
			Expect(togglelist[channel].Status).To(Equal(controlStatus))
		})
	})
})
