package main

import (
	"strings"

	"github.com/codegangsta/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cli flags", func() {
	Context("When defined in codegangsta", func() {
		controlFields := map[string]bool{
			strings.Join(opsManagerHostFlag, ", "): true,
			strings.Join(adminUserFlag, ", "):      true,
			strings.Join(adminPassFlag, ", "):      true,
			strings.Join(opsManagerUserFlag, ", "): true,
			strings.Join(opsManagerPassFlag, ", "): true,
			strings.Join(destFlag, ", "):           true,
			strings.Join(tilelistFlag, ", "):       true,
		}

		It("should have all registered fields", func() {

			for _, v := range backupRestoreFlags {
				flagName := v.(cli.StringFlag).Name
				val, ok := controlFields[flagName]
				Ω(val).Should(BeTrue())
				Ω(ok).Should(BeTrue())
			}
		})

		It("should contain the proper number of flags available", func() {
			Ω(len(backupRestoreFlags)).Should(Equal(len(controlFields)))
		})
	})
})
