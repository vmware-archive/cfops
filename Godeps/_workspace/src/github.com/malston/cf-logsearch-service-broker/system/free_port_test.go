package system_test

import (
	"github.com/pivotal-golang/lager"
	"net"
	"os"
	"regexp"
	"strconv"

	"github.com/pivotal-cf/cf-redis-broker/system"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Next available port", func() {
	var logger lager.Logger

	BeforeEach(func() {
		logger = lager.NewLogger("free-port-test")
		logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	})

	It("finds a free tcp port", func() {
		port, _ := system.FindFreePort()
		portStr := strconv.Itoa(port)
		logger.Debug("next available tcp port", lager.Data{"port": portStr})

		matched, err := regexp.MatchString("^[0-9]+$", portStr)
		Ω(matched).To(Equal(true))

		l, err := net.Listen("tcp", ":"+portStr)
		Ω(err).ToNot(HaveOccurred())
		l.Close()
	})

})
