package main

import (
	"os"

	"github.com/codegangsta/cli"
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
	app.Commands = append(app.Commands, []cli.Command{
		cli.Command{
			Name: "version",
			Action: func(c *cli.Context) {
				cli.ShowVersion(c)
			},
		},
	}...)
	return app
}
