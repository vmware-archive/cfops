package cfops_test

import (
	. "github.com/pivotalservices/cfops"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tile Runner", func() {
	Describe("SetupSupportedTiles", func() {
		Context("when called with a valid flagSet", func() {
			It("should not panic", func() {
				Ω(func() {
					SetupSupportedTiles(&mockFlagSet{})
				}).ShouldNot(Panic())
			})
		})
	})

	Describe("RunPipeline", func() {

		Context("restore action", func() {
			testWithoutTilelist(Restore)
			testWithValidTilelist(Restore)
			testWithInvalidTileList(Restore)
			testPipelineExecutionError(Restore)
		})

		Context("backup action", func() {
			testWithoutTilelist(Backup)
			testWithValidTilelist(Backup)
			testWithInvalidTileList(Backup)
			testPipelineExecutionError(Backup)
		})

	})
})
