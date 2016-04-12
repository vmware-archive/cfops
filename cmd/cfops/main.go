package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfbackup/tileregistry"
	_ "github.com/pivotalservices/cfbackup/tiles"
	_ "github.com/pivotalservices/cfops/plugin/load"
)

var (
	VERSION string
)

type ErrorHandler struct {
	ExitCode int
	Error    error
}

func main() {
	eh := new(ErrorHandler)
	eh.ExitCode = 0
	app := NewApp(eh)
	app.Run(os.Args)
	os.Exit(eh.ExitCode)
}

// NewApp creates a new cli app
func NewApp(eh *ErrorHandler) *cli.App {
	cli.AppHelpTemplate = CfopsHelpTemplate
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
		CreateBURACliCommand("backup", "creates a backup archive of the target tile", eh),
		CreateBURACliCommand("restore", "restores from an archive to the target tile", eh),
	}
	return app
}
