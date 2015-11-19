package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfops/tileregistry"
	"github.com/xchapter7x/lo"
)

func CreateBURACliCommand(name string, usage string, eh *ErrorHandler) (command cli.Command) {
	desc := fmt.Sprintf("%s --opsmanagerhost <host> --adminuser <usr> --adminpass <pass> --opsmanageruser <opsuser> --opsmanagerpass <opspass> -d <dir> --tile elastic-runtime", name)
	command = cli.Command{
		Name:        name,
		Usage:       usage,
		Description: desc,
		Flags:       buraFlags,
		Action:      buraAction(name, eh),
	}
	return
}

func buraAction(commandName string, eh *ErrorHandler) (action func(*cli.Context)) {
	action = func(c *cli.Context) {
		var (
			fs = &flagSet{
				host:           c.String(flagList[opsManagerHost].Flag[0]),
				adminUser:      c.String(flagList[adminUser].Flag[0]),
				adminPass:      c.String(flagList[adminPass].Flag[0]),
				opsManagerUser: c.String(flagList[opsManagerUser].Flag[0]),
				opsManagerPass: c.String(flagList[opsManagerPass].Flag[0]),
				dest:           c.String(flagList[dest].Flag[0]),
				tile:           c.String(flagList[tile].Flag[0]),
			}
		)

		if err := getTileFromRegistry(fs, commandName); err == nil {
			lo.G.Debug("Tile action completed successfully")

		} else {
			lo.G.Error("there was an error getting tile from registry:", err.Error())
			cli.ShowCommandHelp(c, commandName)
			eh.ExitCode = helpExitCode
			eh.Error = err
		}
	}
	return
}

func runTileAction(commandName string, tile tileregistry.Tile) (err error) {
	switch commandName {
	case "backup":
		err = tile.Backup()
	case "restore":
		err = tile.Restore()
	}
	return
}

func getTileFromRegistry(fs *flagSet, commandName string) (err error) {
	var tile tileregistry.Tile
	lo.G.Debug("checking registry for available tile")
	lo.G.Error(fs.Tile())

	if tileBuilder, ok := tileregistry.GetRegistry()[fs.Tile()]; ok {
		lo.G.Debug("found tile in registry")

		if hasValidBackupRestoreFlags(fs) {
			lo.G.Debug("we have all required flags and a proper builder: ", tileBuilder)
			tile, err = tileBuilder.New(tileregistry.TileSpec{
				OpsManagerHost:   fs.Host(),
				AdminUser:        fs.AdminUser(),
				AdminPass:        fs.AdminPass(),
				OpsManagerUser:   fs.OpsManagerUser(),
				OpsManagerPass:   fs.OpsManagerPass(),
				ArchiveDirectory: fs.Dest(),
			})
			err = runTileAction(commandName, tile)

		} else {
			err = ErrInvalidFlagArgs
		}

	} else {
		err = ErrInvalidTileSelection
	}
	return
}

var buraFlags = func() (flags []cli.Flag) {
	for _, v := range flagList {
		flags = append(flags, cli.StringFlag{
			Name:   strings.Join(v.Flag, ", "),
			Value:  "",
			Usage:  v.Desc,
			EnvVar: v.EnvVar,
		})
	}
	return
}()

const (
	errExitCode           = 1
	helpExitCode          = 2
	cleanExitCode         = 0
	opsManagerHost string = "opsmanagerHost"
	adminUser      string = "adminUser"
	adminPass      string = "adminPass"
	opsManagerUser string = "opsManagerUser"
	opsManagerPass string = "opsManagerPass"
	dest           string = "destination"
	tile           string = "tile"
)

var (
	ErrInvalidFlagArgs      = errors.New("invalid cli flag args")
	ErrInvalidTileSelection = errors.New("invalid tile selected. try the 'list-tiles' option to see a list of available tiles.")
	flagList                = map[string]flagBucket{
		opsManagerHost: flagBucket{
			Flag:   []string{"opsmanagerhost", "omh"},
			Desc:   "hostname for Ops Manager",
			EnvVar: "CFOPS_HOST",
		},
		adminUser: flagBucket{
			Flag:   []string{"adminuser", "du"},
			Desc:   "username for Ops Mgr admin (Ops Manager WebConsole Credentials)",
			EnvVar: "CFOPS_ADMIN_USER",
		},
		adminPass: flagBucket{
			Flag:   []string{"adminpass", "dp"},
			Desc:   "password for Ops Mgr admin (Ops Manager WebConsole Credentials)",
			EnvVar: "CFOPS_ADMIN_PASS",
		},
		opsManagerUser: flagBucket{
			Flag:   []string{"opsmanageruser", "omu"},
			Desc:   "username for Ops Manager VM Access (used for ssh connections)",
			EnvVar: "CFOPS_OM_USER",
		},
		opsManagerPass: flagBucket{
			Flag:   []string{"opsmanagerpass", "omp"},
			Desc:   "password for Ops Manager VM Access (used for ssh connections)",
			EnvVar: "CFOPS_OM_PASS",
		},
		dest: flagBucket{
			Flag:   []string{"destination", "d"},
			Desc:   "path of the Cloud Foundry archive",
			EnvVar: "CFOPS_DEST_PATH",
		},
		tile: flagBucket{
			Flag:   []string{"tile", "t"},
			Desc:   "a tile you would like to run the operation on",
			EnvVar: "CFOPS_TILE",
		},
	}
)

type (
	flagSet struct {
		host           string
		adminUser      string
		adminPass      string
		opsManagerUser string
		opsManagerPass string
		dest           string
		tile           string
	}

	flagBucket struct {
		Flag   []string
		Desc   string
		EnvVar string
	}
)

func (s *flagSet) Host() string {
	return s.host
}

func (s *flagSet) AdminUser() string {
	return s.adminUser
}

func (s *flagSet) AdminPass() string {
	return s.adminPass
}

func (s *flagSet) OpsManagerUser() string {
	return s.opsManagerUser
}

func (s *flagSet) OpsManagerPass() string {
	return s.opsManagerPass
}

func (s *flagSet) Dest() string {
	return s.dest
}

func (s *flagSet) Tile() string {
	return s.tile
}

func hasValidBackupRestoreFlags(fs *flagSet) bool {
	res := (fs.Host() != "" &&
		fs.AdminUser() != "" &&
		fs.AdminPass() != "" &&
		fs.OpsManagerUser() != "" &&
		fs.OpsManagerPass() != "" &&
		fs.Dest() != "" &&
		fs.Tile() != "")

	if res == false {
		lo.G.Debug("OpsManagerHost: ", fs.Host())
		lo.G.Debug("AdminUser: ", fs.AdminUser())
		lo.G.Debug("AdminPass: ", fs.AdminPass())
		lo.G.Debug("OpsManagerUser: ", fs.OpsManagerUser())
		lo.G.Debug("OpsManagerPass: ", fs.OpsManagerPass())
		lo.G.Debug("Destination: ", fs.Dest())
	}
	return res
}
