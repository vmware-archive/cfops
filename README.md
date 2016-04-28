<a href="http://cfops.io"> <img src="logos/original-logos-2016-Feb-8505-11461245.png" width="70" ></a>
[![wercker status](https://app.wercker.com/status/d0a50d426b77a9f73da0fe4f383ad624/s/master "wercker status")](https://app.wercker.com/project/bykey/d0a50d426b77a9f73da0fe4f383ad624) [![GoDoc](http://godoc.org/github.com/pivotalservices/cfops?status.png)](http://godoc.org/github.com/pivotalservices/cfops) [![Gitter](https://badges.gitter.im/pivotalservices/cfops.svg)](https://gitter.im/pivotalservices/cfops?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
======

### Version Compatibility

**IAAS**
Supports PivotalCF on:
   - **VSphere**
      - ops-manager (backup/restore `verified`)
      - elastic-runtime
         - postgres datastore (backup/restore `verified`)
         - mysql datastore (backup/restore `verified`)
         - external datastore (backup/restore `no yet supported`)
   - **AWS**
      - ops-manager (backup `verified`)
      - elastic-runtime (`not yet supported`)


### Overview

CFOPS - is a self contained binary, which has no dependencies

This is simply an automation that is based on the supported way to back up Pivotal Cloud Foundry (http://docs.pivotal.io/pivotalcf/customizing/backup-restore/backup-pcf.html).

It may be extended in the future to support greater breadth of functionality.

**Backing up ER will take Cloud Controller offline for the duration of the backup, causing your foundation to become readonly for the duration of the backup**. App pushes etc will not work during this time.

### Install

Download the latest version here:
https://github.com/pivotalservices/cfops/releases/latest

### Contributing

PRs welcome. To get started follow these steps:


* Install [Go 1.6.x](https://golang.org)
* Create a directory where you would like to store the source for Go projects and their binaries (e.g. `$HOME/go`)
* Set an environment variable, `GOPATH`, pointing at the directory you created
* Get the `cf` source: `go get github.com/pivotalservices/cfops` (Ignore any warnings about "no buildable Go source files")
* [Fork this repository](https://help.github.com/articles/fork-a-repo/), adding your fork as a remote
* Install all the required tools - [glide](https://github.com/Masterminds/glide), [wercker cli](http://wercker.com/cli/)
```
$ brew install glide
$ brew tap wercker/wercker
$ brew install wercker-cli
```
* Wrecker requires a local installation of docker, ensure you have docker in place with the relevant environments set.
```
$ docker ps -a # should return something meaningful
```
* Pull in glide managed dependencies:
```
  $ cd $GOPATH/src/github.com/pivotalservices/cfops
  $ glide install
```
* Build the project:
```
  $ cd cmd/cfops/
  $ go build
```
* At this point you should see the cfops binary in the cmd/cfops folder
* Run wercker integration tests
```
$ cd $GOPATH/src/github.com/pivotalservices/cfops
$ ./testrunner
```
* At this point you have everything needed for local development, hack away and submit a [pull request](https://help.github.com/articles/using-pull-requests/) to the `develop` branch


### Differences between v1 and v2

While the core package cfops uses to do the backup & restore has not changed
the cli of cfops itself has been changed.

Major differences are:
- one must specify a single tile to backup/restore at a time.
- one can list the tiles that are supported via the cli.
- more to come...


### Usage (cfops v2.x.x+)

**general commands**

```
NAME:
   cfops - Cloud Foundry Operations Tool

USAGE:
   ./cfops command [command options] [arguments...]

COMMANDS:
   version	shows the application version currently in use
   list-tiles	shows a list of available backup/restore target tiles
   backup	creates a backup archive of the target tile
   restore	restores from an archive to the target tile
   help, h	Shows a list of commands or help for one command
```

**setting log levels**

```
LOG_LEVEL=(debug|info|error) ./cfops backup ...
```

**list available tiles**

```
$ cfops list-tiles
Available Tiles:
ops-manager
elastic-runtime
```

**list version**

```
$ cfops version
cfops version v2.0.0
```

** setting S3 domain **

```
export S3_DOMAIN=<some_compatible_s3_store_url>
```

**run a backup on a tile**

```
NAME:
   ./cfops backup - creates a backup archive of the target tile

USAGE:
   ./cfops backup [command options] [arguments...]

DESCRIPTION:
   backup --opsmanagerhost <host> --adminuser <usr> --adminpass <pass> --opsmanageruser <opsuser> --opsmanagerpass <opspass> -d <dir> --tile elastic-runtime

OPTIONS:
   --destination, -d 		path of the Cloud Foundry archive [$CFOPS_DEST_PATH]
   --tile, -t 			a tile you would like to run the operation on [$CFOPS_TILE]
   --opsmanagerhost, --omh 	hostname for Ops Manager [$CFOPS_HOST]
   --adminuser, --du 		username for Ops Mgr admin (Ops Manager WebConsole Credentials) [$CFOPS_ADMIN_USER]
   --adminpass, --dp 		password for Ops Mgr admin (Ops Manager WebConsole Credentials) [$CFOPS_ADMIN_PASS]
   --opsmanageruser, --omu 	username for Ops Manager VM Access (used for ssh connections) [$CFOPS_OM_USER]
   --opsmanagerpass, --omp 	password for Ops Manager VM Access (used for ssh connections) [$CFOPS_OM_PASS]
```


**run a restore on a tile**

```
NAME:
   ./cfops restore - restores from an archive to the target tile

USAGE:
   ./cfops restore [command options] [arguments...]   

DESCRIPTION:
   restore --opsmanagerhost <host> --adminuser <usr> --adminpass <pass> --opsmanageruser <opsuser> --opsmanagerpass <opspass> -d <dir> --tile elastic-runtime

OPTIONS:
   --adminuser, --du 		username for Ops Mgr admin (Ops Manager WebConsole Credentials) [$CFOPS_ADMIN_USER]
   --adminpass, --dp 		password for Ops Mgr admin (Ops Manager WebConsole Credentials) [$CFOPS_ADMIN_PASS]
   --opsmanageruser, --omu 	username for Ops Manager VM Access (used for ssh connections) [$CFOPS_OM_USER]
   --opsmanagerpass, --omp 	password for Ops Manager VM Access (used for ssh connections) [$CFOPS_OM_PASS]
   --destination, -d 		path of the Cloud Foundry archive [$CFOPS_DEST_PATH]
   --tile, -t 			a tile you would like to run the operation on [$CFOPS_TILE]
   --opsmanagerhost, --omh 	hostname for Ops Manager [$CFOPS_HOST]
```

---



### Usage (cfops v1.x.x)

For example you can try the various commands, args and flags (and --help documentation) that are currently proposed, such as:

    $ ./cfops -help

    $ ./cfops backup

    $ ./cfops restore

etc.


Sample help output:
```
$ ./cfops help backup
NAME:
   backup - backup --opsmanagerhost <host> --adminuser <usr> --adminpass <pass> --opsmanageruser <opsuser> --opsmanagerpass <opspass> -d <dir> --tl 'opsmanager, er'

USAGE:
   command backup [command options] [arguments...]

DESCRIPTION:
   backup a Cloud Foundry deployment, including Ops Manager configuration, databases, and blob store

OPTIONS:
   --tilelist, --tl 		a csv list of the tiles you would like to run the operation on [$CFOPS_TILE_LIST]
   --opsmanagerhost, --omh 	hostname for Ops Manager [$CFOPS_HOST]
   --adminuser, --du 		username for Ops Mgr admin (Ops Manager WebConsole Credentials) [$CFOPS_ADMIN_USER]
   --adminpass, --dp 		password for Ops Mgr admin (Ops Manager WebConsole Credentials) [$CFOPS_ADMIN_PASS]
   --opsmanageruser, --omu 	username for Ops Manager VM Access (used for ssh connections) [$CFOPS_OM_USER]
   --opsmanagerpass, --omp 	password for Ops Manager VM Access (used for ssh connections) [$CFOPS_OM_PASS]
   --destination, -d 		path of the Cloud Foundry backup archive [$CFOPS_BACKUP_PATH]







$ ./cfops help restore
NAME:
   restore - restore --opsmanagerhost <host> --adminuser <usr> --adminpass <pass> --opsmanageruser <opsuser> --opsmanagerpass <opspass> -d <dir> --tl 'opsmanager, er'

USAGE:
   command restore [command options] [arguments...]

DESCRIPTION:
   Restore a Cloud Foundry deployment, including Ops Manager configuration, databases, and blob store

OPTIONS:
   --tilelist, --tl 		a csv list of the tiles you would like to run the operation on [$CFOPS_TILE_LIST]
   --opsmanagerhost, --omh 	hostname for Ops Manager [$CFOPS_HOST]
   --adminuser, --du 		username for Ops Mgr admin (Ops Manager WebConsole Credentials) [$CFOPS_ADMIN_USER]
   --adminpass, --dp 		password for Ops Mgr admin (Ops Manager WebConsole Credentials) [$CFOPS_ADMIN_PASS]
   --opsmanageruser, --omu 	username for Ops Manager VM Access (used for ssh connections) [$CFOPS_OM_USER]
   --opsmanagerpass, --omp 	password for Ops Manager VM Access (used for ssh connections) [$CFOPS_OM_PASS]
   --destination, -d 		path of the Cloud Foundry backup archive [$CFOPS_BACKUP_PATH]

```
