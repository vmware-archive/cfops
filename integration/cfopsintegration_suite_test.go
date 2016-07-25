package cfopsintegration_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBrokerintegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cfops Integration Suite")
}
