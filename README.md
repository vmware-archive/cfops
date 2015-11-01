cfops [![wercker status](https://app.wercker.com/status/d0a50d426b77a9f73da0fe4f383ad624/s/master "wercker status")](https://app.wercker.com/project/bykey/d0a50d426b77a9f73da0fe4f383ad624) [![GoDoc](http://godoc.org/github.com/pivotalservices/cfops?status.png)](http://godoc.org/github.com/pivotalservices/cfops)
======

### Version Compatibility

(as of release v1.1.0+)

**IAAS**
This utility is only designed to work for PivotalCF on VSphere

**BACKUP**
tested and known to work for **Ops Manager v1.5**
tested and known to work for **ER v1.5**

**RESTORE**
tested and known to work for **Ops Manager v1.5**
!! Partially working for **ER v1.5** 
   * Note: Mysql automatic restore failing. (ccdb, uaadb, consoledb & nfs all working)


### Overview

This is simply an automation that is based on the supported way to back up Pivotal Cloud Foundry (http://docs.pivotal.io/pivotalcf/customizing/backup-settings.html).

It may be extended in the future to support greater breadth of functionality.


### Install

Download the latest version here:
https://github.com/pivotalservices/cfops/releases

### Contributing

PRs welcome.

### Usage

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





