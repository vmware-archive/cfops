[![wercker status](https://app.wercker.com/status/d0a50d426b77a9f73da0fe4f383ad624/m/master "wercker status")](https://app.wercker.com/project/bykey/d0a50d426b77a9f73da0fe4f383ad624)

[![GoDoc](http://godoc.org/github.com/pivotalservices/cfops?status.png)](http://godoc.org/github.com/pivotalservices/cfops)


CF Ops
======

A Cloud Foundry Operations tool for IaaS installation, deployment, and management automation


### Background

This is simply an automation that is based on the supported way to back up Pivotal Cloud Foundry (http://docs.pivotal.io/pivotalcf/customizing/backup-settings.html).

It may be extended in the future to support greater breadth of functionality.

PRs welcome.

### Install

download latest version for your system. details here:
https://github.com/pivotalservices/cfops/wiki


### Current

This initial version *only* provides backup and restore.

The project is written in "Go".

For example you can try the various commands, args and flags (and --help documentation) that are currently proposed, such as:

    $ ./cfops -help

    $ ./cfops backup

    $ ./cfops restore

etc.


Sample help output:
```
    $ ./cfops
NAME:
   cfops - Cloud Foundry Operations tool for IaaS installation, deployment, and management automation

USAGE:
   cfops [global options] command [command options] [arguments...]

VERSION:
   0.0.0

AUTHOR(S):

COMMANDS:
   backup, b	backup -host <host> -u <usr> -p <pass> --tp <tpass> -d <dir> --tl 'opsmanager, er'
   restore, r	restore --host <host> -u <usr> -p <pass> --tp <tpass> -d <dir>  --tl 'opsmanager, er'
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
    
    
    
    $ ./cfops help backup
NAME:
   backup - backup -host <host> -u <usr> -p <pass> --tp <tpass> -d <dir> --tl 'opsmanager, er'

USAGE:
   command backup [command options] [arguments...]

DESCRIPTION:
   backup a Cloud Foundry deployment, including Ops Manager configuration, databases, and blob store

OPTIONS:
   --hostname, --host 		hostname for Ops Manager
   --username, -u 		username for Ops Manager
   --password, -p 		password for Ops Manager
   --tempestpassword, --tp 	password for the Ops Manager tempest user
   --destination, -d 		directory of the Cloud Foundry backup archive
   --tilelist, --tl 		a csv list of the tiles you would like to run the operation on
   
   $ ./cfops help restore
NAME:
   restore - restore -host <host> -u <usr> -p <pass> --tp <tpass> -d <dir> --tl 'opsmanager, er'

USAGE:
   command restore [command options] [arguments...]

DESCRIPTION:
   restore a Cloud Foundry deployment, including Ops Manager configuration, databases, and blob store

OPTIONS:
   --hostname, --host 		hostname for Ops Manager
   --username, -u 		username for Ops Manager
   --password, -p 		password for Ops Manager
   --tempestpassword, --tp 	password for the Ops Manager tempest user
   --destination, -d 		directory of the Cloud Foundry backup archive
   --tilelist, --tl 		a csv list of the tiles you would like to run the operation on

```





