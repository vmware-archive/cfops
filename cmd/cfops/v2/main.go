package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	_ "github.com/pivotalservices/cfbackup"
	"github.com/pivotalservices/cfops/tileregistry"
)

var (
	VERSION  string
	ExitCode int
	logLevel string
)

func main() {
	app := NewApp()
	app.Run(os.Args)
	os.Exit(ExitCode)
}

// NewApp creates a new cli app
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Version = VERSION
	app.Name = "cfops"
	app.Usage = "Cloud Foundry Operations Tool"
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "version",
			Usage: "shows the application version currently in use",
			Action: func(c *cli.Context) {
				cli.ShowVersion(c)
			},
		},
		cli.Command{
			Name:  "list-tiles",
			Usage: "shows a list of available backup/restore target tiles",
			Action: func(c *cli.Context) {
				fmt.Println("Available Tiles:")
				for n, _ := range tileregistry.GetRegistry() {
					fmt.Println(n)
				}
			},
		},
		cli.Command{
			Name:  "backup",
			Usage: "creates a backup archive of the target tile",
			Action: func(c *cli.Context) {
				cli.ShowVersion(c)
			},
		},
		cli.Command{
			Name:  "restore",
			Usage: "restores from an archive to the target tile",
			Action: func(c *cli.Context) {
				cli.ShowVersion(c)
			},
		},
	}
	return app
}
