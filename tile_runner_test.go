package cfops_test

import (
	. "github.com/pivotalservices/cfops"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

var _ = Describe("Tile Runner", func() {
	Describe("RunPipeline", func() {

		Context("restore action", func() {
			testWithoutTilelist(Restore)
			testWithTilelist(Restore)
			testPipelineExecutionError(Restore)
		})

		Context("backup action", func() {
			testWithoutTilelist(Backup)
			testWithTilelist(Backup)
			testPipelineExecutionError(Backup)
		})

	})
})
