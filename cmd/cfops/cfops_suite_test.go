package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCfops(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cfops")
}
