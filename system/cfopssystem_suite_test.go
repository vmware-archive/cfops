package system

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var cfConfig Config

func TestBrokerintegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cfops System Test Suite")
}
