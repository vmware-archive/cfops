cfops [![wercker status](https://app.wercker.com/status/d0a50d426b77a9f73da0fe4f383ad624/s/master "wercker status")](https://app.wercker.com/project/bykey/d0a50d426b77a9f73da0fe4f383ad624) [![GoDoc](http://godoc.org/github.com/pivotalservices/cfops?status.png)](http://godoc.org/github.com/pivotalservices/cfops)
======

### Version Compatibility

(as of release v1.1.0+)

**IAAS**
This utility is only designed to work for PivotalCF on VSphere

**BACKUP**
tested and known to work for **Ops Manager v1.6**
tested and known to work for **ER v1.6**

**RESTORE**
tested and known to work for **Ops Manager v1.6**
tested and known to work for **ER v1.6** 


### Overview

This is simply an automation that is based on the supported way to back up Pivotal Cloud Foundry (http://docs.pivotal.io/pivotalcf/customizing/backup-restore/backup-pcf.html).

It may be extended in the future to support greater breadth of functionality.

**Backing up ER will take Cloud Controller offline for the duration of the backup, causing your foundation to become readonly for the duration of the backup**. App pushes etc will not work during this time. 

### Install

Download the latest version here:
https://github.com/pivotalservices/cfops/releases/latest

### Contributing

PRs welcome.


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





