package createfs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCreatefs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Createfs Suite")
}
