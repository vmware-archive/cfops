package persistence_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPersistance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TestPersistance Suite")
}
