package gogobosh_test

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoGoBosh(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoGoBOSH suite")
}
