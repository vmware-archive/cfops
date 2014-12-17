package itertools_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestItertools(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "itertools Suite")
}
