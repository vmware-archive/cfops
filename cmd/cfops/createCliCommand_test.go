package main_test

import (
	"flag"
	"fmt"

	"github.com/codegangsta/cli"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cfbackup/tileregistry"
	"github.com/pivotalservices/cfbackup/tileregistry/fake"
	. "github.com/pivotalservices/cfops/cmd/cfops"
)

var _ = Describe("given a CreateBURACliCommand func", func() {
	testTileAction("backup")
	testTileAction("restore")
})

func testTileAction(actionName string) {
	var (
		controlErrorHandler = new(ErrorHandler)
	)

	Context(fmt.Sprintf("when called with the %s command name", actionName), func() {
		controlName := actionName
		controlUsage := "something about the usage here"

		It("Then it should return a cli.Command with required values", func() {
			cmd := CreateBURACliCommand(controlName, controlUsage, controlErrorHandler)
			Ω(cmd.Name).Should(Equal(controlName))
			Ω(cmd.Usage).Should(Equal(controlUsage))
			Ω(cmd.Action).ShouldNot(BeNil())
		})

		Describe(fmt.Sprintf("given a %s command with an Action value", actionName), func() {
			var (
				controlExit       = 0
				controlCliContext *cli.Context
				controlCmd        cli.Command
			)
			BeforeEach(func() {
				controlErrorHandler = new(ErrorHandler)
				controlErrorHandler.ExitCode = controlExit
			})

			Context("when the action is called with all proper flags on a registered tile", func() {
				var (
					controlTileName      = "fake-tile"
					controlTileGenerator *fake.TileGenerator
					controlTile          *fake.Tile
					controlCloser        *fake.Closer
				)
				BeforeEach(func() {
					controlTile = new(fake.Tile)
					controlCloser = &fake.Closer{Executions: 0}
					controlTileGenerator = new(fake.TileGenerator)
					controlTileGenerator.TileSpy = controlTile
					controlTileGenerator.Closer = controlCloser
					tileregistry.Register(controlTileName, controlTileGenerator)
					set := flag.NewFlagSet("", 0)
					set.String("tile", controlTileName, "")
					set.String("opsmanagerhost", "*****", "")
					set.String("adminuser", "*****", "")
					set.String("adminpass", "*****", "")
					set.String("opsmanageruser", "*****", "")
					set.String("opsmanagerpass", "*****", "")
					set.String("destination", "*****", "")
					set.String("pluginargs", "*****", "")
					controlCliContext = cli.NewContext(cli.NewApp(), set, nil)
					controlCmd = CreateBURACliCommand(controlName, controlUsage, controlErrorHandler)
					controlCliContext.Command = controlCmd
				})
				It("then it should execute a call on the tiles Action func()", func() {
					controlCmd.Action.(func(*cli.Context))(controlCliContext)
					Ω(controlErrorHandler.ExitCode).Should(Equal(controlExit))
					Ω(controlErrorHandler.Error).ShouldNot(HaveOccurred())
					Ω(controlCloser.Executions).Should(Equal(1))
					switch controlName {
					case "backup":
						Ω(controlTile.BackupCallCount).ShouldNot(Equal(0))
					case "restore":
						Ω(controlTile.RestoreCallCount).ShouldNot(Equal(0))
					default:
						panic("this should never happen")
					}
				})
			})

			Context("when the action is called w/o a matching registered tile", func() {
				BeforeEach(func() {
					set := flag.NewFlagSet(controlName, 0)
					set.Bool("myflag", false, "doc")
					set.String("otherflag", "hello world", "doc")
					controlCliContext = cli.NewContext(cli.NewApp(), set, nil)
					controlCmd = CreateBURACliCommand(controlName, controlUsage, controlErrorHandler)
					controlCliContext.Command = controlCmd
				})

				It("then it should set an error and failure exit code", func() {
					controlCmd.Action.(func(*cli.Context))(controlCliContext)
					Ω(controlErrorHandler.ExitCode).ShouldNot(Equal(controlExit))
					Ω(controlErrorHandler.Error).Should(HaveOccurred())
					Ω(controlErrorHandler.Error).Should(Equal(ErrInvalidTileSelection))
				})
			})

			Context("when running a tile action returns an error", func() {
				var (
					controlTileName      = "fake-tile"
					controlTileGenerator *fake.TileGenerator
					controlTile          *fake.Tile
					controlCloser        *fake.Closer
				)
				BeforeEach(func() {
					controlTile = new(fake.Tile)
					controlCloser = &fake.Closer{Executions: 0}
					controlTile.ErrFake = fmt.Errorf("operation failed")
					controlTileGenerator = new(fake.TileGenerator)
					controlTileGenerator.TileSpy = controlTile
					controlTileGenerator.Closer = controlCloser
					tileregistry.Register(controlTileName, controlTileGenerator)
					set := flag.NewFlagSet("", 0)
					set.String("tile", controlTileName, "")
					set.String("opsmanagerhost", "*****", "")
					set.String("adminuser", "*****", "")
					set.String("adminpass", "*****", "")
					set.String("opsmanageruser", "*****", "")
					set.String("opsmanagerpass", "*****", "")
					set.String("destination", "*****", "")
					controlCliContext = cli.NewContext(cli.NewApp(), set, nil)
					controlCmd = CreateBURACliCommand(controlName, controlUsage, controlErrorHandler)
					controlCliContext.Command = controlCmd
				})
				It("then it should set an error and failure exit code", func() {
					controlCmd.Action.(func(*cli.Context))(controlCliContext)
					Ω(controlErrorHandler.ExitCode).ShouldNot(Equal(controlExit))
					Ω(controlErrorHandler.Error).Should(HaveOccurred())
				})
			})

			Context("when the action is called w/o proper flags but with a registered tile", func() {
				var controlTileName = "fake-tile"
				BeforeEach(func() {
					tileregistry.Register(controlTileName, new(fake.TileGenerator))
					set := flag.NewFlagSet("", 0)
					set.String("tile", controlTileName, "doc")
					set.Bool("myflag", false, "doc")
					set.String("otherflag", "hello world", "doc")
					controlCliContext = cli.NewContext(cli.NewApp(), set, nil)
					controlCmd = CreateBURACliCommand(controlName, controlUsage, controlErrorHandler)
					controlCliContext.Command = controlCmd
				})
				It("then it should set an error and failure exit code", func() {
					controlCmd.Action.(func(*cli.Context))(controlCliContext)
					Ω(controlErrorHandler.ExitCode).ShouldNot(Equal(controlExit))
					Ω(controlErrorHandler.Error).Should(HaveOccurred())
					Ω(controlErrorHandler.Error).Should(Equal(ErrInvalidFlagArgs))
				})
			})

			Context("when a tile builder returns an error", func() {
				var (
					controlTileName      = "fake-tile"
					controlTileGenerator *fake.TileGenerator
					controlTile          *fake.Tile
				)
				BeforeEach(func() {
					controlTile = new(fake.Tile)
					controlTileGenerator = new(fake.TileGenerator)
					controlTileGenerator.TileSpy = controlTile
					controlTileGenerator.ErrFake = fmt.Errorf("operation timed out")
					tileregistry.Register(controlTileName, controlTileGenerator)
					set := flag.NewFlagSet("", 0)
					set.String("tile", controlTileName, "")
					set.String("opsmanagerhost", "*****", "")
					set.String("adminuser", "*****", "")
					set.String("adminpass", "*****", "")
					set.String("opsmanageruser", "*****", "")
					set.String("opsmanagerpass", "*****", "")
					set.String("destination", "*****", "")
					controlCliContext = cli.NewContext(cli.NewApp(), set, nil)
					controlCmd = CreateBURACliCommand(controlName, controlUsage, controlErrorHandler)
					controlCliContext.Command = controlCmd
				})
				It("then it should set an error and failure exit code", func() {
					controlCmd.Action.(func(*cli.Context))(controlCliContext)
					Ω(controlErrorHandler.ExitCode).ShouldNot(Equal(controlExit))
					Ω(controlErrorHandler.Error).Should(HaveOccurred())
				})
			})
		})
	})
}
