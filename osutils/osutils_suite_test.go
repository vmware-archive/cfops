package osutils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOsutils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Osutils Suite")
}
