package cfbackup_test

import (
	"io/ioutil"
	"os"

	. "github.com/pivotalservices/cfbackup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backup / Restore", func() {
	var (
		dir string
	)

	BeforeEach(func() {
		dir, _ = ioutil.TempDir("", "cfops-backup")
	})

	AfterEach(func() {
		os.RemoveAll(dir)
	})

	Describe("RunBackupPipeline", func() {
		var (
			origRunPipeline func(actionBuilder func(Tile) func() error, tiles []Tile) (err error)
			err             error
		)

		BeforeEach(func() {
			origRunPipeline = RunPipeline
			RunPipeline = func(actionBuilder func(Tile) func() error, tiles []Tile) (err error) {
				return
			}
			err = RunBackupPipeline("", "", "", "", "")
		})

		AfterEach(func() {
			RunPipeline = origRunPipeline
		})

		Context("empty argument call", func() {
			It("should return error", func() {
				Ω(err).ShouldNot(BeNil())
			})
		})
	})

	Describe("RunRestorePipeline", func() {
		var (
			origRunPipeline func(actionBuilder func(Tile) func() error, tiles []Tile) (err error)
			err             error
		)

		BeforeEach(func() {
			origRunPipeline = RunPipeline
			RunPipeline = func(actionBuilder func(Tile) func() error, tiles []Tile) (err error) {
				return
			}
			err = RunRestorePipeline("", "", "", "", "")
		})

		AfterEach(func() {
			RunPipeline = origRunPipeline
		})

		Context("empty argument call", func() {
			It("should return error", func() {
				Ω(err).ShouldNot(BeNil())
			})
		})
	})

	Describe("RunPipeline", func() {
		var err error
		var tile *mockTile

		Context("backup", func() {
			var (
				controlBackupCount int
			)

			BeforeEach(func() {
				tile = &mockTile{}
				controlBackupCount = tile.BackupCalled
				err = RunPipeline(TILE_BACKUP_ACTION, []Tile{tile})
			})

			It("should return nil error and successfully call Backup function", func() {
				Ω(err).Should(BeNil())
				Ω(tile.BackupCalled).Should(BeNumerically(">", controlBackupCount))
			})

			Context("failed backup call", func() {
				BeforeEach(func() {
					tile = &mockTile{ErrBackup: mockTileBackupError}
					controlBackupCount = tile.BackupCalled
					err = RunPipeline(TILE_BACKUP_ACTION, []Tile{tile})
				})

				It("should return backup error and successfully call Backup function", func() {
					Ω(err).Should(Equal(mockTileBackupError))
					Ω(tile.BackupCalled).Should(BeNumerically(">", controlBackupCount))
				})
			})
		})

		Context("restore", func() {
			var (
				controlRestoreCount int
			)

			BeforeEach(func() {
				tile = &mockTile{}
				controlRestoreCount = tile.RestoreCalled
				err = RunPipeline(TILE_RESTORE_ACTION, []Tile{tile})
			})

			It("should return nil error and successfully call Restore function", func() {
				Ω(err).Should(BeNil())
				Ω(tile.RestoreCalled).Should(BeNumerically(">", controlRestoreCount))
			})

			Context("failed restore call", func() {
				BeforeEach(func() {
					tile = &mockTile{ErrRestore: mockTileRestoreError}
					controlRestoreCount = tile.RestoreCalled
					err = RunPipeline(TILE_RESTORE_ACTION, []Tile{tile})
				})

				It("should return restore error and successfully call Restore function", func() {
					Ω(err).Should(Equal(mockTileRestoreError))
					Ω(tile.RestoreCalled).Should(BeNumerically(">", controlRestoreCount))
				})
			})
		})
	})
})
