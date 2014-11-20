package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestConfiguration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

type TestConfig struct {
	TestString      string
	TestInt         int
	TestStringSlice []string
	TestBigInt      uint32
}
