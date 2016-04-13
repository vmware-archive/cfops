package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cfbackup/tileregistry"
	"github.com/xchapter7x/lo"
)

//CreateBURACliCommand - this will create a cli command object for backup / restore
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
				host:                 c.String(flagList[opsManagerHost].Flag[0]),
				adminUser:            c.String(flagList[adminUser].Flag[0]),
				adminPass:            c.String(flagList[adminPass].Flag[0]),
				opsManagerUser:       c.String(flagList[opsManagerUser].Flag[0]),
				opsManagerPass:       c.String(flagList[opsManagerPass].Flag[0]),
				opsManagerPassphrase: c.String(flagList[opsManagerPassphrase].Flag[0]),
				dest:                 c.String(flagList[dest].Flag[0]),
				tile:                 c.String(flagList[tile].Flag[0]),
				encryptionKey:        c.String(flagList[encryptionKey].Flag[0]),
				clearBoshManifest:    c.Bool(flagList[clearBoshManifest].Flag[0]),
				pluginArgs:           c.String(flagList[pluginArgs].Flag[0]),
			}
		)

		if tileCloser, err := getTileFromRegistry(fs, commandName); err == nil {
			defer tileCloser.Close()
			if err = runTileAction(commandName, tileCloser); err != nil {
				lo.G.Errorf("there was an error: %s running %s on %s tile:%v", err.Error(), commandName, fs.Tile(), tile)
				exitOnError(eh, c, commandName, err)
				return
			}
		} else {
			lo.G.Errorf("there was an error getting tile from registry: %s", err.Error())
			exitOnError(eh, c, commandName, err)
			return
		}
		lo.G.Debug("Tile action completed successfully")
	}
	return
}

func exitOnError(eh *ErrorHandler, c *cli.Context, commandName string, err error) {
	cli.ShowCommandHelp(c, commandName)
	eh.ExitCode = helpExitCode
	eh.Error = err
}

func runTileAction(commandName string, tile tileregistry.Tile) (err error) {
	lo.G.Debug("Running %s for tile: %+v", commandName, tile)
	switch commandName {
	case "backup":
		err = tile.Backup()
	case "restore":
		err = tile.Restore()
	}
	return
}

func getTileFromRegistry(fs *flagSet, commandName string) (tileCloser tileregistry.TileCloser, err error) {
	lo.G.Debug("checking registry for '%s' tile", fs.Tile())

	if tileBuilder, ok := tileregistry.GetRegistry()[fs.Tile()]; ok {
		lo.G.Debug("found tile in registry")

		if hasValidBackupRestoreFlags(fs) {
			lo.G.Debug("we have all required flags and a proper builder")
			tileCloser, err = tileBuilder.New(tileregistry.TileSpec{
				OpsManagerHost:       fs.Host(),
				AdminUser:            fs.AdminUser(),
				AdminPass:            fs.AdminPass(),
				OpsManagerUser:       fs.OpsManagerUser(),
				OpsManagerPass:       fs.OpsManagerPass(),
				OpsManagerPassphrase: fs.OpsManagerPassphrase(),
				ArchiveDirectory:     fs.Dest(),
				CryptKey:             fs.EncryptionKey(),
				ClearBoshManifest:    fs.ClearBoshManifest(),
				PluginArgs:           fs.PluginArgs(),
			})
			if err != nil {
				return nil, fmt.Errorf("failure to connect to ops manager host: %s", err.Error())
			}

		} else {
			err = ErrInvalidFlagArgs
		}

	} else {
		lo.G.Errorf("tile '%s' not found", fs.Tile())
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
	errExitCode                 = 1
	helpExitCode                = 2
	cleanExitCode               = 0
	opsManagerHost       string = "opsmanagerHost"
	adminUser            string = "adminUser"
	adminPass            string = "adminPass"
	opsManagerUser       string = "opsManagerUser"
	opsManagerPass       string = "opsManagerPass"
	opsManagerPassphrase string = "opsManagerPassphrase"
	dest                 string = "destination"
	tile                 string = "tile"
	encryptionKey        string = "encryptionKey"
	clearBoshManifest    string = "clearboshmanifest"
	pluginArgs           string = "pluginArgs"
)

var (
	//ErrInvalidFlagArgs - error for invalid flags
	ErrInvalidFlagArgs = errors.New("invalid cli flag args")
	//ErrInvalidTileSelection - error for invalid tile
	ErrInvalidTileSelection = errors.New("invalid tile selected. try the 'list-tiles' option to see a list of available tiles")
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
		opsManagerPassphrase: flagBucket{
			Flag:   []string{"opsmanagerpassphrase", "omr"},
			Desc:   "passphrase is used by Ops Manager 1.7 to decrypt the installation files during restore",
			EnvVar: "CFOPS_OM_PASSPHRASE",
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
		encryptionKey: flagBucket{
			Flag:   []string{"encryptionkey", "k"},
			Desc:   "encryption key to encrypt/decrypt your archive (key lengths supported are 16, 24, 32 for AES-128, AES-192, or AES-256)",
			EnvVar: "CFOPS_ENCRYPTION_KEY",
		},
		clearBoshManifest: flagBucket{
			Flag:   []string{"clear-bosh-manifest"},
			Desc:   "set this flag if you would like to clear the bosh-deployments.yml (this should only affect a restore of Ops-Manager)",
			EnvVar: "CFOPS_CLEAR_BOSH_MANIFEST",
		},
		pluginArgs: flagBucket{
			Flag:   []string{"pluginargs", "p"},
			Desc:   "Arguments for plugin to execute",
			EnvVar: "CFOPS_PLUGIN_ARGS",
		},
	}
)

type (
	flagSet struct {
		host                 string
		adminUser            string
		adminPass            string
		opsManagerUser       string
		opsManagerPass       string
		opsManagerPassphrase string
		dest                 string
		tile                 string
		encryptionKey        string
		clearBoshManifest    bool
		pluginArgs           string
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

func (s *flagSet) OpsManagerPassphrase() string {
	return s.opsManagerPassphrase
}

func (s *flagSet) Dest() string {
	return s.dest
}

func (s *flagSet) Tile() string {
	return s.tile
}

func (s *flagSet) EncryptionKey() string {
	return s.encryptionKey
}

func (s *flagSet) ClearBoshManifest() bool {
	return s.clearBoshManifest
}

func (s *flagSet) PluginArgs() string {
	return s.pluginArgs
}

func hasValidBackupRestoreFlags(fs *flagSet) bool {
	res := (fs.Host() != "" &&
		fs.AdminUser() != "" &&
		fs.AdminPass() != "" &&
		fs.OpsManagerUser() != "" &&
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
